package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/fili258-oss/BackendAppStocks/internal/domain/entity"
	"github.com/fili258-oss/BackendAppStocks/internal/domain/repository"
)

// RecommendationRepositoryImpl implementa repository.RecommendationRepository
type RecommendationRepositoryImpl struct {
	db *sql.DB
}

// NewRecommendationRepository crea una nueva instancia de RecommendationRepositoryImpl
func NewRecommendationRepository(db *sql.DB) repository.RecommendationRepository {
	return &RecommendationRepositoryImpl{db: db}
}

// Create inserta una nueva recomendación
func (r *RecommendationRepositoryImpl) Create(ctx context.Context, rec *entity.Recommendation) error {
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}

	// Convertir metrics map a JSON
	metricsJSON, err := json.Marshal(rec.Metrics)
	if err != nil {
		return fmt.Errorf("error al serializar metrics: %w", err)
	}

	query := `
		INSERT INTO recommendations (
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12
		)
	`

	_, err = r.db.ExecContext(ctx, query,
		rec.ID, rec.StockID, rec.StockSymbol, rec.Type, rec.Score, rec.Confidence,
		rec.Reason, rec.Strategy, metricsJSON, rec.ValidUntil,
		rec.CreatedAt, rec.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al crear recommendation: %w", err)
	}

	return nil
}

// Update actualiza una recomendación existente
func (r *RecommendationRepositoryImpl) Update(ctx context.Context, rec *entity.Recommendation) error {
	rec.UpdatedAt = time.Now()

	metricsJSON, err := json.Marshal(rec.Metrics)
	if err != nil {
		return fmt.Errorf("error al serializar metrics: %w", err)
	}

	query := `
		UPDATE recommendations SET
			type = $2, score = $3, confidence = $4,
			reason = $5, strategy = $6, metrics = $7, valid_until = $8,
			updated_at = $9
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		rec.ID, rec.Type, rec.Score, rec.Confidence,
		rec.Reason, rec.Strategy, metricsJSON, rec.ValidUntil,
		rec.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar recommendation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrRecommendationNotFound
	}

	return nil
}

// Delete elimina una recomendación por su ID
func (r *RecommendationRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM recommendations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar recommendation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrRecommendationNotFound
	}

	return nil
}

// FindByID busca una recomendación por su ID
func (r *RecommendationRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Recommendation, error) {
	query := `
		SELECT 
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		FROM recommendations
		WHERE id = $1
	`

	rec := &entity.Recommendation{}
	var metricsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rec.ID, &rec.StockID, &rec.StockSymbol, &rec.Type, &rec.Score, &rec.Confidence,
		&rec.Reason, &rec.Strategy, &metricsJSON, &rec.ValidUntil,
		&rec.CreatedAt, &rec.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrRecommendationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar recommendation por ID: %w", err)
	}

	// Deserializar metrics JSON
	if err := json.Unmarshal(metricsJSON, &rec.Metrics); err != nil {
		return nil, fmt.Errorf("error al deserializar metrics: %w", err)
	}

	return rec, nil
}

// FindByStockID busca todas las recomendaciones de un stock
func (r *RecommendationRepositoryImpl) FindByStockID(ctx context.Context, stockID string) ([]*entity.Recommendation, error) {
	query := `
		SELECT 
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		FROM recommendations
		WHERE stock_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, stockID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar recommendations por stock_id: %w", err)
	}
	defer rows.Close()

	return r.scanRecommendations(rows)
}

// FindByType busca recomendaciones por tipo
func (r *RecommendationRepositoryImpl) FindByType(ctx context.Context, recType entity.RecommendationType, limit int) ([]*entity.Recommendation, error) {
	query := `
		SELECT 
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		FROM recommendations
		WHERE type = $1
		ORDER BY score DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, recType, limit)
	if err != nil {
		return nil, fmt.Errorf("error al buscar recommendations por tipo: %w", err)
	}
	defer rows.Close()

	return r.scanRecommendations(rows)
}

// FindTopRecommendations retorna las mejores N recomendaciones
func (r *RecommendationRepositoryImpl) FindTopRecommendations(ctx context.Context, limit int) ([]*entity.Recommendation, error) {
	query := `
		SELECT 
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		FROM recommendations
		WHERE valid_until > CURRENT_TIMESTAMP
		ORDER BY score DESC, confidence DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error al buscar top recommendations: %w", err)
	}
	defer rows.Close()

	return r.scanRecommendations(rows)
}

// FindValidRecommendations retorna solo recomendaciones válidas
func (r *RecommendationRepositoryImpl) FindValidRecommendations(ctx context.Context, limit int) ([]*entity.Recommendation, error) {
	query := `
		SELECT 
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		FROM recommendations
		WHERE valid_until > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error al buscar valid recommendations: %w", err)
	}
	defer rows.Close()

	return r.scanRecommendations(rows)
}

// FindByDateRange busca recomendaciones en un rango de fechas
func (r *RecommendationRepositoryImpl) FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Recommendation, error) {
	query := `
		SELECT 
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		FROM recommendations
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("error al buscar recommendations por rango de fechas: %w", err)
	}
	defer rows.Close()

	return r.scanRecommendations(rows)
}

// BulkCreate inserta múltiples recomendaciones
func (r *RecommendationRepositoryImpl) BulkCreate(ctx context.Context, recommendations []*entity.Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO recommendations (
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12
		)
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement: %w", err)
	}
	defer stmt.Close()

	for _, rec := range recommendations {
		if rec.ID == "" {
			rec.ID = uuid.New().String()
		}

		metricsJSON, err := json.Marshal(rec.Metrics)
		if err != nil {
			return fmt.Errorf("error al serializar metrics: %w", err)
		}

		_, err = stmt.ExecContext(ctx,
			rec.ID, rec.StockID, rec.StockSymbol, rec.Type, rec.Score, rec.Confidence,
			rec.Reason, rec.Strategy, metricsJSON, rec.ValidUntil,
			rec.CreatedAt, rec.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("error al insertar recommendation: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al hacer commit: %w", err)
	}

	return nil
}

// DeleteExpired elimina recomendaciones expiradas
func (r *RecommendationRepositoryImpl) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM recommendations WHERE valid_until < CURRENT_TIMESTAMP`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("error al eliminar recommendations expiradas: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	return rowsAffected, nil
}

// GetLatestByStock retorna la recomendación más reciente para un stock
func (r *RecommendationRepositoryImpl) GetLatestByStock(ctx context.Context, stockID string) (*entity.Recommendation, error) {
	query := `
		SELECT 
			id, stock_id, stock_symbol, type, score, confidence,
			reason, strategy, metrics, valid_until,
			created_at, updated_at
		FROM recommendations
		WHERE stock_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	rec := &entity.Recommendation{}
	var metricsJSON []byte

	err := r.db.QueryRowContext(ctx, query, stockID).Scan(
		&rec.ID, &rec.StockID, &rec.StockSymbol, &rec.Type, &rec.Score, &rec.Confidence,
		&rec.Reason, &rec.Strategy, &metricsJSON, &rec.ValidUntil,
		&rec.CreatedAt, &rec.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrRecommendationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar latest recommendation: %w", err)
	}

	if err := json.Unmarshal(metricsJSON, &rec.Metrics); err != nil {
		return nil, fmt.Errorf("error al deserializar metrics: %w", err)
	}

	return rec, nil
}

// Count retorna el número total de recomendaciones
func (r *RecommendationRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM recommendations`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar recommendations: %w", err)
	}
	return count, nil
}

// scanRecommendations es un helper para escanear múltiples filas
func (r *RecommendationRepositoryImpl) scanRecommendations(rows *sql.Rows) ([]*entity.Recommendation, error) {
	var recommendations []*entity.Recommendation

	for rows.Next() {
		rec := &entity.Recommendation{}
		var metricsJSON []byte

		err := rows.Scan(
			&rec.ID, &rec.StockID, &rec.StockSymbol, &rec.Type, &rec.Score, &rec.Confidence,
			&rec.Reason, &rec.Strategy, &metricsJSON, &rec.ValidUntil,
			&rec.CreatedAt, &rec.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear recommendation: %w", err)
		}

		// Deserializar metrics JSON
		if err := json.Unmarshal(metricsJSON, &rec.Metrics); err != nil {
			return nil, fmt.Errorf("error al deserializar metrics: %w", err)
		}

		recommendations = append(recommendations, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas: %w", err)
	}

	return recommendations, nil
}
