package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/marino/stock-analyzer/internal/application/dto"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/domain/repository"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
)

// FetchStocksUseCase caso de uso para traer stocks desde la API externa
type FetchStocksUseCase struct {
	stockRepo repository.StockRepository
	adapter   *external.FinnhubAdapter
}

// NewFetchStocksUseCase crea una nueva instancia del use case
func NewFetchStocksUseCase(
	stockRepo repository.StockRepository,
	adapter *external.FinnhubAdapter,
) *FetchStocksUseCase {
	return &FetchStocksUseCase{
		stockRepo: stockRepo,
		adapter:   adapter,
	}
}

// Execute ejecuta el caso de uso
func (uc *FetchStocksUseCase) Execute(ctx context.Context, request dto.FetchStocksRequest) (*dto.FetchStocksResponse, error) {
	if len(request.Symbols) == 0 {
		return nil, fmt.Errorf("no symbols provided")
	}

	response := &dto.FetchStocksResponse{
		Stocks: make([]dto.StockDTO, 0),
		Errors: make([]string, 0),
	}

	// Obtener stocks desde la API
	stocks, err := uc.adapter.GetMultipleStocks(ctx, request.Symbols)
	if err != nil {
		log.Printf("Warning: some stocks failed to fetch: %v", err)
	}

	// Procesar cada stock obtenido
	for _, stock := range stocks {
		// Si Save es true, guardar o actualizar en BD
		if request.Save {
			if err := uc.saveOrUpdateStock(ctx, stock); err != nil {
				response.Failed++
				response.Errors = append(response.Errors, 
					fmt.Sprintf("Error saving %s: %v", stock.Symbol, err))
				continue
			}
		}

		// Agregar a respuesta
		response.Stocks = append(response.Stocks, dto.ToStockDTO(stock))
		response.Success++
	}

	// Calcular stocks que fallaron
	response.Failed = len(request.Symbols) - response.Success

	return response, nil
}

// saveOrUpdateStock guarda o actualiza un stock en la BD
func (uc *FetchStocksUseCase) saveOrUpdateStock(ctx context.Context, stock *entity.Stock) error {
	// Intentar encontrar stock existente
	existing, err := uc.stockRepo.FindBySymbol(ctx, stock.Symbol)
	
	if err == entity.ErrStockNotFound {
		// No existe, crear nuevo
		return uc.stockRepo.Create(ctx, stock)
	}
	
	if err != nil {
		return fmt.Errorf("error checking existing stock: %w", err)
	}

	// Existe, actualizar
	stock.ID = existing.ID
	stock.CreatedAt = existing.CreatedAt
	return uc.stockRepo.Update(ctx, stock)
}
