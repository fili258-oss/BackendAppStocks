-- Migration: 002_create_recommendations_table.down.sql
-- Descripción: Revierte la creación de la tabla recommendations
-- Fecha: 2025-01-30

-- Eliminar índices
DROP INDEX IF EXISTS idx_recommendations_metrics;
DROP INDEX IF EXISTS idx_recommendations_stock_created;
DROP INDEX IF EXISTS idx_recommendations_strategy;
DROP INDEX IF EXISTS idx_recommendations_created_at;
DROP INDEX IF EXISTS idx_recommendations_valid_until;
DROP INDEX IF EXISTS idx_recommendations_score;
DROP INDEX IF EXISTS idx_recommendations_type;
DROP INDEX IF EXISTS idx_recommendations_stock_symbol;
DROP INDEX IF EXISTS idx_recommendations_stock_id;

-- Eliminar tabla
DROP TABLE IF EXISTS recommendations;

-- Eliminar tipo ENUM
DROP TYPE IF EXISTS recommendation_type;
