package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/domain/repository"
)

// StockRepositoryImpl implementa repository.StockRepository
// Implementa el patrón Repository para acceso a datos de stocks
type StockRepositoryImpl struct {
	db *sql.DB
}

// NewStockRepository crea una nueva instancia de StockRepositoryImpl
func NewStockRepository(db *sql.DB) repository.StockRepository {
	return &StockRepositoryImpl{db: db}
}

// Create inserta un nuevo stock en la base de datos
func (r *StockRepositoryImpl) Create(ctx context.Context, stock *entity.Stock) error {
	// Generar UUID si no existe
	if stock.ID == "" {
		stock.ID = uuid.New().String()
	}

	query := `
		INSERT INTO stocks (
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18,
			$19, $20
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		stock.ID, stock.Symbol, stock.Name, stock.Exchange, stock.Currency,
		stock.Price, stock.OpenPrice, stock.HighPrice, stock.LowPrice, stock.ClosePrice,
		stock.Change, stock.ChangePercent, stock.Volume, stock.MarketCap, stock.PERatio,
		stock.DividendYield, stock.Week52High, stock.Week52Low,
		stock.CreatedAt, stock.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al crear stock: %w", err)
	}

	return nil
}

// Update actualiza un stock existente
func (r *StockRepositoryImpl) Update(ctx context.Context, stock *entity.Stock) error {
	stock.UpdatedAt = time.Now()

	query := `
		UPDATE stocks SET
			name = $2, exchange = $3, currency = $4,
			price = $5, open_price = $6, high_price = $7, low_price = $8, close_price = $9,
			change = $10, change_percent = $11, volume = $12, market_cap = $13, pe_ratio = $14,
			dividend_yield = $15, week_52_high = $16, week_52_low = $17,
			updated_at = $18
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		stock.ID, stock.Name, stock.Exchange, stock.Currency,
		stock.Price, stock.OpenPrice, stock.HighPrice, stock.LowPrice, stock.ClosePrice,
		stock.Change, stock.ChangePercent, stock.Volume, stock.MarketCap, stock.PERatio,
		stock.DividendYield, stock.Week52High, stock.Week52Low,
		stock.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrStockNotFound
	}

	return nil
}

// Delete elimina un stock por su ID
func (r *StockRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM stocks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrStockNotFound
	}

	return nil
}

// FindByID busca un stock por su ID
func (r *StockRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Stock, error) {
	query := `
		SELECT 
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		FROM stocks
		WHERE id = $1
	`

	stock := &entity.Stock{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stock.ID, &stock.Symbol, &stock.Name, &stock.Exchange, &stock.Currency,
		&stock.Price, &stock.OpenPrice, &stock.HighPrice, &stock.LowPrice, &stock.ClosePrice,
		&stock.Change, &stock.ChangePercent, &stock.Volume, &stock.MarketCap, &stock.PERatio,
		&stock.DividendYield, &stock.Week52High, &stock.Week52Low,
		&stock.CreatedAt, &stock.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrStockNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar stock por ID: %w", err)
	}

	return stock, nil
}

// FindBySymbol busca un stock por su símbolo
func (r *StockRepositoryImpl) FindBySymbol(ctx context.Context, symbol string) (*entity.Stock, error) {
	query := `
		SELECT 
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		FROM stocks
		WHERE symbol = $1
	`

	stock := &entity.Stock{}
	err := r.db.QueryRowContext(ctx, query, symbol).Scan(
		&stock.ID, &stock.Symbol, &stock.Name, &stock.Exchange, &stock.Currency,
		&stock.Price, &stock.OpenPrice, &stock.HighPrice, &stock.LowPrice, &stock.ClosePrice,
		&stock.Change, &stock.ChangePercent, &stock.Volume, &stock.MarketCap, &stock.PERatio,
		&stock.DividendYield, &stock.Week52High, &stock.Week52Low,
		&stock.CreatedAt, &stock.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrStockNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar stock por símbolo: %w", err)
	}

	return stock, nil
}

// FindAll retorna todos los stocks con paginación
func (r *StockRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]*entity.Stock, error) {
	query := `
		SELECT 
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		FROM stocks
		ORDER BY symbol
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al buscar todos los stocks: %w", err)
	}
	defer rows.Close()

	return r.scanStocks(rows)
}

// Search busca stocks por nombre o símbolo con filtros
func (r *StockRepositoryImpl) Search(ctx context.Context, query string, filters repository.StockFilters) ([]*entity.Stock, error) {
	// Construir query SQL dinámicamente basado en filtros
	sqlQuery := `
		SELECT 
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		FROM stocks
		WHERE (LOWER(symbol) LIKE LOWER($1) OR LOWER(name) LIKE LOWER($1))
	`

	args := []interface{}{"%" + query + "%"}
	argIndex := 2

	// Aplicar filtros adicionales
	if filters.Exchange != "" {
		sqlQuery += fmt.Sprintf(" AND exchange = $%d", argIndex)
		args = append(args, filters.Exchange)
		argIndex++
	}

	if filters.MinPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, filters.MinPrice)
		argIndex++
	}

	if filters.MaxPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, filters.MaxPrice)
		argIndex++
	}

	if filters.MinVolume > 0 {
		sqlQuery += fmt.Sprintf(" AND volume >= $%d", argIndex)
		args = append(args, filters.MinVolume)
		argIndex++
	}

	if filters.MinMarketCap > 0 {
		sqlQuery += fmt.Sprintf(" AND market_cap >= $%d", argIndex)
		args = append(args, filters.MinMarketCap)
		argIndex++
	}

	// Ordenamiento
	orderBy := "symbol"
	if filters.SortBy != "" {
		orderBy = filters.SortBy
	}

	sortDirection := "ASC"
	if filters.SortDirection == "desc" {
		sortDirection = "DESC"
	}

	sqlQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT 100", orderBy, sortDirection)

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error al buscar stocks: %w", err)
	}
	defer rows.Close()

	return r.scanStocks(rows)
}

// BulkCreate inserta múltiples stocks de manera eficiente
func (r *StockRepositoryImpl) BulkCreate(ctx context.Context, stocks []*entity.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	// Usar transacción para bulk insert
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO stocks (
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18,
			$19, $20
		)
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement: %w", err)
	}
	defer stmt.Close()

	for _, stock := range stocks {
		if stock.ID == "" {
			stock.ID = uuid.New().String()
		}

		_, err := stmt.ExecContext(ctx,
			stock.ID, stock.Symbol, stock.Name, stock.Exchange, stock.Currency,
			stock.Price, stock.OpenPrice, stock.HighPrice, stock.LowPrice, stock.ClosePrice,
			stock.Change, stock.ChangePercent, stock.Volume, stock.MarketCap, stock.PERatio,
			stock.DividendYield, stock.Week52High, stock.Week52Low,
			stock.CreatedAt, stock.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("error al insertar stock %s: %w", stock.Symbol, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al hacer commit: %w", err)
	}

	return nil
}

// Count retorna el número total de stocks
func (r *StockRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM stocks`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar stocks: %w", err)
	}
	return count, nil
}

// GetTopByVolume retorna los N stocks con mayor volumen
func (r *StockRepositoryImpl) GetTopByVolume(ctx context.Context, limit int) ([]*entity.Stock, error) {
	query := `
		SELECT 
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		FROM stocks
		ORDER BY volume DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error al buscar top por volumen: %w", err)
	}
	defer rows.Close()

	return r.scanStocks(rows)
}

// GetTopByMarketCap retorna los N stocks con mayor capitalización
func (r *StockRepositoryImpl) GetTopByMarketCap(ctx context.Context, limit int) ([]*entity.Stock, error) {
	query := `
		SELECT 
			id, symbol, name, exchange, currency,
			price, open_price, high_price, low_price, close_price,
			change, change_percent, volume, market_cap, pe_ratio,
			dividend_yield, week_52_high, week_52_low,
			created_at, updated_at
		FROM stocks
		ORDER BY market_cap DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error al buscar top por market cap: %w", err)
	}
	defer rows.Close()

	return r.scanStocks(rows)
}

// scanStocks es un helper para escanear múltiples filas de stocks
func (r *StockRepositoryImpl) scanStocks(rows *sql.Rows) ([]*entity.Stock, error) {
	var stocks []*entity.Stock

	for rows.Next() {
		stock := &entity.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Symbol, &stock.Name, &stock.Exchange, &stock.Currency,
			&stock.Price, &stock.OpenPrice, &stock.HighPrice, &stock.LowPrice, &stock.ClosePrice,
			&stock.Change, &stock.ChangePercent, &stock.Volume, &stock.MarketCap, &stock.PERatio,
			&stock.DividendYield, &stock.Week52High, &stock.Week52Low,
			&stock.CreatedAt, &stock.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas: %w", err)
	}

	return stocks, nil
}
