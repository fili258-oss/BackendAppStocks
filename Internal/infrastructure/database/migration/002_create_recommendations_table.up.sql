-- Migration: 002_create_recommendations_table.up.sql
-- Descripción: Crea la tabla recommendations para almacenar recomendaciones de inversión
-- Fecha: 2025-01-30

-- Tipo ENUM para tipo de recomendación
CREATE TYPE recommendation_type AS ENUM ('BUY', 'SELL', 'HOLD', 'STRONG_BUY');

-- Crear tabla recommendations
CREATE TABLE IF NOT EXISTS recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_id UUID NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    
    -- Tipo y scoring
    type recommendation_type NOT NULL,
    score DECIMAL(5, 2) NOT NULL CHECK (score >= 0 AND score <= 100),
    confidence DECIMAL(3, 2) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    
    -- Información adicional
    reason TEXT,
    strategy VARCHAR(50) NOT NULL,
    metrics JSONB DEFAULT '{}',
    
    -- Validez
    valid_until TIMESTAMP NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT current_timestamp(),
    updated_at TIMESTAMP DEFAULT current_timestamp(),
    
    -- Foreign key a stocks
    CONSTRAINT fk_recommendations_stock
        FOREIGN KEY (stock_id) 
        REFERENCES stocks(id)
        ON DELETE CASCADE
);

-- Índices para mejorar performance
CREATE INDEX IF NOT EXISTS idx_recommendations_stock_id ON recommendations(stock_id);
CREATE INDEX IF NOT EXISTS idx_recommendations_stock_symbol ON recommendations(stock_symbol);
CREATE INDEX IF NOT EXISTS idx_recommendations_type ON recommendations(type);
CREATE INDEX IF NOT EXISTS idx_recommendations_score ON recommendations(score DESC);
CREATE INDEX IF NOT EXISTS idx_recommendations_valid_until ON recommendations(valid_until);
CREATE INDEX IF NOT EXISTS idx_recommendations_created_at ON recommendations(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_recommendations_strategy ON recommendations(strategy);

-- Índice compuesto para búsquedas comunes
CREATE INDEX IF NOT EXISTS idx_recommendations_stock_created 
    ON recommendations(stock_id, created_at DESC);

-- Índice para JSONB metrics (permite búsquedas eficientes en JSON)
CREATE INDEX IF NOT EXISTS idx_recommendations_metrics 
    ON recommendations USING GIN (metrics);

-- Comentarios en la tabla
COMMENT ON TABLE recommendations IS 'Almacena recomendaciones de inversión generadas por algoritmos';
COMMENT ON COLUMN recommendations.score IS 'Score de 0 a 100 para la recomendación';
COMMENT ON COLUMN recommendations.confidence IS 'Nivel de confianza de 0 a 1';
COMMENT ON COLUMN recommendations.metrics IS 'Métricas usadas en el cálculo (formato JSON)';
COMMENT ON COLUMN recommendations.valid_until IS 'Fecha hasta la cual la recomendación es válida';
