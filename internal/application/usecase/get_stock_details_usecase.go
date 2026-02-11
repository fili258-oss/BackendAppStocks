package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/marino/stock-analyzer/internal/application/dto"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/domain/repository"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
)

// GetStockDetailsUseCase caso de uso para obtener detalles de un stock
type GetStockDetailsUseCase struct {
	stockRepo repository.StockRepository
	adapter   *external.FinnhubAdapter
}

// NewGetStockDetailsUseCase crea una nueva instancia del use case
func NewGetStockDetailsUseCase(
	stockRepo repository.StockRepository,
	adapter *external.FinnhubAdapter,
) *GetStockDetailsUseCase {
	return &GetStockDetailsUseCase{
		stockRepo: stockRepo,
		adapter:   adapter,
	}
}

// Execute ejecuta el caso de uso
// Busca primero en BD, si no encuentra o está desactualizado, busca en API
func (uc *GetStockDetailsUseCase) Execute(ctx context.Context, symbol string) (*dto.StockDTO, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	// Intentar obtener de BD
	stock, err := uc.stockRepo.FindBySymbol(ctx, symbol)
	
	// Si no existe en BD, obtener de API
	if err == entity.ErrStockNotFound {
		return uc.fetchFromAPI(ctx, symbol)
	}
	
	if err != nil {
		return nil, fmt.Errorf("error fetching from DB: %w", err)
	}

	// Si existe pero está desactualizado (más de 5 minutos), actualizar
	if time.Since(stock.UpdatedAt) > 5*time.Minute {
		return uc.refreshStock(ctx, stock)
	}

	// Retornar stock de BD
	stockDTO := dto.ToStockDTO(stock)
	return &stockDTO, nil
}

// ExecuteFromDB fuerza obtener solo de BD (no busca en API)
func (uc *GetStockDetailsUseCase) ExecuteFromDB(ctx context.Context, symbol string) (*dto.StockDTO, error) {
	stock, err := uc.stockRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	stockDTO := dto.ToStockDTO(stock)
	return &stockDTO, nil
}

// ExecuteFromAPI fuerza obtener de API (ignora BD)
func (uc *GetStockDetailsUseCase) ExecuteFromAPI(ctx context.Context, symbol string) (*dto.StockDTO, error) {
	return uc.fetchFromAPI(ctx, symbol)
}

// fetchFromAPI obtiene stock de API y lo guarda en BD
func (uc *GetStockDetailsUseCase) fetchFromAPI(ctx context.Context, symbol string) (*dto.StockDTO, error) {
	// Obtener de API
	stock, err := uc.adapter.GetStock(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("error fetching from API: %w", err)
	}

	// Intentar guardar en BD
	if err := uc.stockRepo.Create(ctx, stock); err != nil {
		// Si falla al guardar, igual retornar los datos
		// (puede fallar si ya existe por race condition)
	}

	stockDTO := dto.ToStockDTO(stock)
	return &stockDTO, nil
}

// refreshStock actualiza stock existente desde API
func (uc *GetStockDetailsUseCase) refreshStock(ctx context.Context, existing *entity.Stock) (*dto.StockDTO, error) {
	// Obtener datos actualizados
	updated, err := uc.adapter.GetStock(ctx, existing.Symbol)
	if err != nil {
		// Si falla la actualización, retornar datos existentes
		stockDTO := dto.ToStockDTO(existing)
		return &stockDTO, nil
	}

	// Mantener ID y CreatedAt
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	// Actualizar en BD
	if err := uc.stockRepo.Update(ctx, updated); err != nil {
		// Si falla actualización, retornar datos nuevos de API de todas formas
		stockDTO := dto.ToStockDTO(updated)
		return &stockDTO, nil
	}

	stockDTO := dto.ToStockDTO(updated)
	return &stockDTO, nil
}
