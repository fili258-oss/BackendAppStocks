package repository

import (
	"context"
	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// StockRepository define el contrato para el acceso a datos de stocks
// Interface segregation principle: solo métodos necesarios
// Dependency inversion: el dominio define la interfaz, infrastructure la implementa
type StockRepository interface {
	// Create guarda un nuevo stock en la base de datos
	Create(ctx context.Context, stock *entity.Stock) error

	// Update actualiza un stock existente
	Update(ctx context.Context, stock *entity.Stock) error

	// Delete elimina un stock por su ID
	Delete(ctx context.Context, id string) error

	// FindByID busca un stock por su ID
	FindByID(ctx context.Context, id string) (*entity.Stock, error)

	// FindBySymbol busca un stock por su símbolo (ticker)
	FindBySymbol(ctx context.Context, symbol string) (*entity.Stock, error)

	// FindAll retorna todos los stocks con paginación
	FindAll(ctx context.Context, limit, offset int) ([]*entity.Stock, error)

	// Search busca stocks por nombre o símbolo con filtros
	Search(ctx context.Context, query string, filters StockFilters) ([]*entity.Stock, error)

	// BulkCreate inserta múltiples stocks de manera eficiente
	// Útil para la carga inicial desde la API externa
	BulkCreate(ctx context.Context, stocks []*entity.Stock) error

	// Count retorna el número total de stocks en la base de datos
	Count(ctx context.Context) (int64, error)

	// GetTopByVolume retorna los N stocks con mayor volumen de transacciones
	GetTopByVolume(ctx context.Context, limit int) ([]*entity.Stock, error)

	// GetTopByMarketCap retorna los N stocks con mayor capitalización
	GetTopByMarketCap(ctx context.Context, limit int) ([]*entity.Stock, error)
}

// StockFilters define los filtros disponibles para búsqueda de stocks
type StockFilters struct {
	Exchange      string  // Filtrar por bolsa
	MinPrice      float64 // Precio mínimo
	MaxPrice      float64 // Precio máximo
	MinVolume     int64   // Volumen mínimo
	MinMarketCap  float64 // Capitalización mínima
	SortBy        string  // Campo para ordenar (price, volume, market_cap)
	SortDirection string  // Dirección de ordenamiento (asc, desc)
}
