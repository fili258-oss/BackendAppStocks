-- Migration: 001_create_stocks_table.down.sql
-- Descripción: Revierte la creación de la tabla stocks
-- Fecha: 2025-01-30

-- Eliminar índices
DROP INDEX IF EXISTS idx_stocks_updated_at;
DROP INDEX IF EXISTS idx_stocks_volume;
DROP INDEX IF EXISTS idx_stocks_market_cap;
DROP INDEX IF EXISTS idx_stocks_exchange;
DROP INDEX IF EXISTS idx_stocks_name;
DROP INDEX IF EXISTS idx_stocks_symbol;

-- Eliminar tabla
DROP TABLE IF EXISTS stocks;
