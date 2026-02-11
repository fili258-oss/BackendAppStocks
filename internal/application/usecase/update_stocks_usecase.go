package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marino/stock-analyzer/internal/application/dto"
	"github.com/marino/stock-analyzer/internal/domain/repository"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
)

// UpdateStocksUseCase caso de uso para actualizar stocks existentes
type UpdateStocksUseCase struct {
	stockRepo repository.StockRepository
	adapter   *external.FinnhubAdapter
}

// NewUpdateStocksUseCase crea una nueva instancia del use case
func NewUpdateStocksUseCase(
	stockRepo repository.StockRepository,
	adapter *external.FinnhubAdapter,
) *UpdateStocksUseCase {
	return &UpdateStocksUseCase{
		stockRepo: stockRepo,
		adapter:   adapter,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateStocksUseCase) Execute(ctx context.Context, request dto.UpdateStocksRequest) (*dto.UpdateStocksResponse, error) {
	response := &dto.UpdateStocksResponse{
		Errors: make([]string, 0),
	}

	// Determinar qué stocks actualizar
	var symbolsToUpdate []string
	
	if len(request.Symbols) > 0 {
		// Actualizar símbolos específicos
		symbolsToUpdate = request.Symbols
	} else {
		// Actualizar todos los stocks en la BD
		stocks, err := uc.stockRepo.FindAll(ctx, 1000, 0) // Límite alto para obtener todos
		if err != nil {
			return nil, fmt.Errorf("error fetching stocks from DB: %w", err)
		}
		
		symbolsToUpdate = make([]string, len(stocks))
		for i, stock := range stocks {
			symbolsToUpdate[i] = stock.Symbol
		}
	}

	if len(symbolsToUpdate) == 0 {
		return response, nil
	}

	// Actualizar cada stock
	for _, symbol := range symbolsToUpdate {
		if err := uc.updateStock(ctx, symbol, request.Force); err != nil {
			response.Failed++
			response.Errors = append(response.Errors, 
				fmt.Sprintf("Error updating %s: %v", symbol, err))
			continue
		}
		response.Updated++
	}

	return response, nil
}

// updateStock actualiza un stock individual
func (uc *UpdateStocksUseCase) updateStock(ctx context.Context, symbol string, force bool) error {
	// Verificar si el stock existe en BD
	existing, err := uc.stockRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return fmt.Errorf("stock not found in DB: %w", err)
	}

	// Si no es forzado, verificar si necesita actualización
	if !force {
		// Si fue actualizado hace menos de 5 minutos, skip
		if time.Since(existing.UpdatedAt) < 5*time.Minute {
			log.Printf("Stock %s recently updated, skipping", symbol)
			return nil
		}
	}

	// Obtener datos actualizados de la API
	updated, err := uc.adapter.GetStock(ctx, symbol)
	if err != nil {
		return fmt.Errorf("error fetching from API: %w", err)
	}

	// Mantener ID y CreatedAt originales
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	// Actualizar en BD
	if err := uc.stockRepo.Update(ctx, updated); err != nil {
		return fmt.Errorf("error updating in DB: %w", err)
	}

	return nil
}
