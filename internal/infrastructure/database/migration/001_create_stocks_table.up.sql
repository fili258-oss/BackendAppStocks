-- Migration: 001_create_stocks_table.up.sql
-- Descripción: Crea la tabla stocks para almacenar información de acciones del mercado
-- Fecha: 2025-01-30

-- Crear tabla stocks
CREATE TABLE IF NOT EXISTS stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    exchange VARCHAR(50),
    currency VARCHAR(10) DEFAULT 'USD',
    
    -- Precios
    price DECIMAL(18, 4) DEFAULT 0,
    open_price DECIMAL(18, 4) DEFAULT 0,
    high_price DECIMAL(18, 4) DEFAULT 0,
    low_price DECIMAL(18, 4) DEFAULT 0,
    close_price DECIMAL(18, 4) DEFAULT 0,
    
    -- Cambios
    change DECIMAL(18, 4) DEFAULT 0,
    change_percent DECIMAL(10, 4) DEFAULT 0,
    
    -- Métricas
    volume BIGINT DEFAULT 0,
    market_cap DECIMAL(20, 2) DEFAULT 0,
    pe_ratio DECIMAL(10, 2) DEFAULT 0,
    dividend_yield DECIMAL(10, 4) DEFAULT 0,
    week_52_high DECIMAL(18, 4) DEFAULT 0,
    week_52_low DECIMAL(18, 4) DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT current_timestamp(),
    updated_at TIMESTAMP DEFAULT current_timestamp()
);

-- Índices para mejorar performance de búsquedas
CREATE INDEX IF NOT EXISTS idx_stocks_symbol ON stocks(symbol);
CREATE INDEX IF NOT EXISTS idx_stocks_name ON stocks(name);
CREATE INDEX IF NOT EXISTS idx_stocks_exchange ON stocks(exchange);
CREATE INDEX IF NOT EXISTS idx_stocks_market_cap ON stocks(market_cap DESC);
CREATE INDEX IF NOT EXISTS idx_stocks_volume ON stocks(volume DESC);
CREATE INDEX IF NOT EXISTS idx_stocks_updated_at ON stocks(updated_at DESC);

-- Comentarios en la tabla y columnas (documentación)
COMMENT ON TABLE stocks IS 'Almacena información de acciones del mercado de valores';
COMMENT ON COLUMN stocks.symbol IS 'Símbolo ticker de la acción (ej: AAPL, GOOGL)';
COMMENT ON COLUMN stocks.market_cap IS 'Capitalización de mercado en USD';
COMMENT ON COLUMN stocks.pe_ratio IS 'Price-to-Earnings ratio';
COMMENT ON COLUMN stocks.dividend_yield IS 'Rendimiento por dividendos en porcentaje';
