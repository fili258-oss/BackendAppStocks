package usecase

import (
	"context"
	"fmt"

	"github.com/marino/stock-analyzer/internal/application/dto"
	"github.com/marino/stock-analyzer/internal/domain/repository"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
)

// SearchStocksUseCase caso de uso para buscar stocks
type SearchStocksUseCase struct {
	stockRepo repository.StockRepository
	adapter   *external.FinnhubAdapter
}

// NewSearchStocksUseCase crea una nueva instancia del use case
func NewSearchStocksUseCase(
	stockRepo repository.StockRepository,
	adapter *external.FinnhubAdapter,
) *SearchStocksUseCase {
	return &SearchStocksUseCase{
		stockRepo: stockRepo,
		adapter:   adapter,
	}
}

// Execute ejecuta el caso de uso
func (uc *SearchStocksUseCase) Execute(ctx context.Context, request dto.StockSearchRequest) (*dto.StockListResponse, error) {
	if request.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	// Establecer valores por defecto
	if request.Limit == 0 {
		request.Limit = 20
	}
	if request.SortBy == "" {
		request.SortBy = "symbol"
	}
	if request.SortDirection == "" {
		request.SortDirection = "asc"
	}

	// Buscar en la base de datos
	filters := repository.StockFilters{
		Exchange:      request.Exchange,
		MinPrice:      request.MinPrice,
		MaxPrice:      request.MaxPrice,
		MinVolume:     request.MinVolume,
		MinMarketCap:  request.MinMarketCap,
		SortBy:        request.SortBy,
		SortDirection: request.SortDirection,
	}

	stocks, err := uc.stockRepo.Search(ctx, request.Query, filters)
	if err != nil {
		return nil, fmt.Errorf("error searching stocks: %w", err)
	}

	// Convertir a DTOs
	stockDTOs := dto.ToStockDTOList(stocks)

	// Aplicar paginación
	total := len(stockDTOs)
	start := request.Offset
	end := request.Offset + request.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedStocks := stockDTOs[start:end]

	return &dto.StockListResponse{
		Stocks: paginatedStocks,
		Total:  total,
		Limit:  request.Limit,
		Offset: request.Offset,
	}, nil
}

// ExecuteWithAPI busca en la API externa si no encuentra en BD
func (uc *SearchStocksUseCase) ExecuteWithAPI(ctx context.Context, request dto.StockSearchRequest) (*dto.StockListResponse, error) {
	// Primero buscar en BD
	response, err := uc.Execute(ctx, request)
	if err != nil {
		return nil, err
	}

	// Si encontramos resultados, retornar
	if len(response.Stocks) > 0 {
		return response, nil
	}

	// Si no encontramos en BD, buscar en API
	limit := request.Limit
	if limit == 0 {
		limit = 10
	}

	apiStocks, err := uc.adapter.SearchStocks(ctx, request.Query, limit)
	if err != nil {
		return nil, fmt.Errorf("error searching in API: %w", err)
	}

	// Convertir a DTOs
	stockDTOs := dto.ToStockDTOList(apiStocks)

	return &dto.StockListResponse{
		Stocks: stockDTOs,
		Total:  len(stockDTOs),
		Limit:  limit,
		Offset: 0,
	}, nil
}
