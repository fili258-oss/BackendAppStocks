package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marino/stock-analyzer/internal/application/dto"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/domain/repository"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
)

// SyncStocksUseCase caso de uso para sincronización batch de stocks
// Útil para mantener la BD actualizada con múltiples stocks
type SyncStocksUseCase struct {
	stockRepo repository.StockRepository
	adapter   *external.FinnhubAdapter
}

// NewSyncStocksUseCase crea una nueva instancia del use case
func NewSyncStocksUseCase(
	stockRepo repository.StockRepository,
	adapter *external.FinnhubAdapter,
) *SyncStocksUseCase {
	return &SyncStocksUseCase{
		stockRepo: stockRepo,
		adapter:   adapter,
	}
}

// Execute ejecuta el caso de uso
func (uc *SyncStocksUseCase) Execute(ctx context.Context, request dto.SyncStocksRequest) (*dto.SyncStocksResponse, error) {
	if len(request.Symbols) == 0 {
		return nil, fmt.Errorf("no symbols provided")
	}

	// Configurar tiempo de actualización por defecto
	updateOlder := request.UpdateOlder
	if updateOlder == 0 {
		updateOlder = 5 * time.Minute
	}

	response := &dto.SyncStocksResponse{
		Errors: make([]string, 0),
	}

	// Procesar cada símbolo
	for _, symbol := range request.Symbols {
		if err := uc.syncStock(ctx, symbol, updateOlder, response); err != nil {
			response.Failed++
			response.Errors = append(response.Errors, 
				fmt.Sprintf("Error syncing %s: %v", symbol, err))
		}
	}

	// Si se solicita eliminar faltantes, procesar
	if request.DeleteMissing {
		deleted, err := uc.deleteMissingStocks(ctx, request.Symbols)
		if err != nil {
			log.Printf("Error deleting missing stocks: %v", err)
		} else {
			response.Deleted = deleted
		}
	}

	return response, nil
}

// syncStock sincroniza un stock individual
func (uc *SyncStocksUseCase) syncStock(ctx context.Context, symbol string, updateOlder time.Duration, response *dto.SyncStocksResponse) error {
	// Verificar si existe en BD
	existing, err := uc.stockRepo.FindBySymbol(ctx, symbol)
	
	// No existe, crear nuevo
	if err != nil {
		return uc.addNewStock(ctx, symbol, response)
	}

	// Existe, verificar si necesita actualización
	if time.Since(existing.UpdatedAt) < updateOlder {
		// No necesita actualización
		return nil
	}

	// Necesita actualización
	return uc.updateExistingStock(ctx, existing, response)
}

// addNewStock agrega un nuevo stock desde la API
func (uc *SyncStocksUseCase) addNewStock(ctx context.Context, symbol string, response *dto.SyncStocksResponse) error {
	// Obtener de API
	stock, err := uc.adapter.GetStock(ctx, symbol)
	if err != nil {
		return fmt.Errorf("error fetching from API: %w", err)
	}

	// Guardar en BD
	if err := uc.stockRepo.Create(ctx, stock); err != nil {
		return fmt.Errorf("error saving to DB: %w", err)
	}

	response.Added++
	log.Printf("Added new stock: %s", symbol)
	return nil
}

// updateExistingStock actualiza un stock existente
func (uc *SyncStocksUseCase) updateExistingStock(ctx context.Context, existing *entity.Stock, response *dto.SyncStocksResponse) error {
	// Obtener datos actualizados
	updated, err := uc.adapter.GetStock(ctx, existing.Symbol)
	if err != nil {
		return fmt.Errorf("error fetching from API: %w", err)
	}

	// Mantener ID y CreatedAt
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	// Actualizar en BD
	if err := uc.stockRepo.Update(ctx, updated); err != nil {
		return fmt.Errorf("error updating in DB: %w", err)
	}

	response.Updated++
	log.Printf("Updated stock: %s", existing.Symbol)
	return nil
}

// deleteMissingStocks elimina stocks que ya no están en la lista
func (uc *SyncStocksUseCase) deleteMissingStocks(ctx context.Context, keepSymbols []string) (int, error) {
	// Obtener todos los stocks de BD
	allStocks, err := uc.stockRepo.FindAll(ctx, 10000, 0)
	if err != nil {
		return 0, fmt.Errorf("error fetching all stocks: %w", err)
	}

	// Crear mapa de símbolos a mantener
	keepMap := make(map[string]bool)
	for _, symbol := range keepSymbols {
		keepMap[symbol] = true
	}

	// Eliminar stocks que no están en la lista
	deleted := 0
	for _, stock := range allStocks {
		if !keepMap[stock.Symbol] {
			if err := uc.stockRepo.Delete(ctx, stock.ID); err != nil {
				log.Printf("Error deleting stock %s: %v", stock.Symbol, err)
				continue
			}
			deleted++
			log.Printf("Deleted stock: %s", stock.Symbol)
		}
	}

	return deleted, nil
}
