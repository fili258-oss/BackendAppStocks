package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
	"github.com/marino/stock-analyzer/internal/infrastructure/config"
)

// Connection representa una conexión a la base de datos
type Connection struct {
	DB *sql.DB
}

// NewConnection crea una nueva conexión a CockroachDB
func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
	// Abrir conexión
	db, err := sql.Open("postgres", cfg.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al hacer ping a la base de datos: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close cierra la conexión a la base de datos
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// HealthCheck verifica que la conexión esté saludable
func (c *Connection) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.DB.PingContext(ctx)
}

// GetStats retorna estadísticas de la conexión
func (c *Connection) GetStats() sql.DBStats {
	return c.DB.Stats()
}
