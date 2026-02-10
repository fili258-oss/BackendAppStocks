package repository

import (
	"context"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"time"
)

// RecommendationRepository define el contrato para el acceso a datos de recomendaciones
type RecommendationRepository interface {
	// Create guarda una nueva recomendación
	Create(ctx context.Context, recommendation *entity.Recommendation) error

	// Update actualiza una recomendación existente
	Update(ctx context.Context, recommendation *entity.Recommendation) error

	// Delete elimina una recomendación por su ID
	Delete(ctx context.Context, id string) error

	// FindByID busca una recomendación por su ID
	FindByID(ctx context.Context, id string) (*entity.Recommendation, error)

	// FindByStockID busca todas las recomendaciones de un stock específico
	FindByStockID(ctx context.Context, stockID string) ([]*entity.Recommendation, error)

	// FindByType busca recomendaciones por tipo (BUY, SELL, etc.)
	FindByType(ctx context.Context, recType entity.RecommendationType, limit int) ([]*entity.Recommendation, error)

	// FindTopRecommendations retorna las mejores N recomendaciones ordenadas por score
	FindTopRecommendations(ctx context.Context, limit int) ([]*entity.Recommendation, error)

	// FindValidRecommendations retorna solo las recomendaciones que aún son válidas
	FindValidRecommendations(ctx context.Context, limit int) ([]*entity.Recommendation, error)

	// FindByDateRange busca recomendaciones creadas en un rango de fechas
	FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Recommendation, error)

	// BulkCreate inserta múltiples recomendaciones
	BulkCreate(ctx context.Context, recommendations []*entity.Recommendation) error

	// DeleteExpired elimina recomendaciones que ya no son válidas
	DeleteExpired(ctx context.Context) (int64, error)

	// GetLatestByStock retorna la recomendación más reciente para un stock
	GetLatestByStock(ctx context.Context, stockID string) (*entity.Recommendation, error)

	// Count retorna el número total de recomendaciones
	Count(ctx context.Context) (int64, error)
}
