package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// DatabaseConfig contiene la configuración para la conexión a la base de datos
type DatabaseConfig struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Database              string
	SSLMode               string
	MaxConnections        int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

// LoadDatabaseConfig carga la configuración de base de datos desde variables de entorno
func LoadDatabaseConfig() (*DatabaseConfig, error) {
	port, err := strconv.Atoi(getEnv("DB_PORT", "26257"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	maxConn, err := strconv.Atoi(getEnv("DB_MAX_CONNECTIONS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONNECTIONS: %w", err)
	}

	maxIdle, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNECTIONS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNECTIONS: %w", err)
	}

	lifetime, err := time.ParseDuration(getEnv("DB_CONNECTION_MAX_LIFETIME", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONNECTION_MAX_LIFETIME: %w", err)
	}

	return &DatabaseConfig{
		Host:                  getEnv("DB_HOST", "localhost"),
		Port:                  port,
		User:                  getEnv("DB_USER", "root"),
		Password:              getEnv("DB_PASSWORD", ""),
		Database:              getEnv("DB_NAME", "stock_analyzer"),
		SSLMode:               getEnv("DB_SSL_MODE", "disable"),
		MaxConnections:        maxConn,
		MaxIdleConnections:    maxIdle,
		ConnectionMaxLifetime: lifetime,
	}, nil
}

// GetConnectionString retorna el string de conexión PostgreSQL
func (c *DatabaseConfig) GetConnectionString() string {
	if c.Password != "" {
		return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
	}
	return fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=%s",
		c.User, c.Host, c.Port, c.Database, c.SSLMode)
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
