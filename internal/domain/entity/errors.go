package entity

import "errors"

// Domain Errors - Errores de negocio del dominio
// Estos errores representan violaciones de reglas de negocio

var (
	// Stock errors
	ErrInvalidSymbol = errors.New("invalid stock symbol: symbol cannot be empty")
	ErrInvalidName   = errors.New("invalid stock name: name cannot be empty")
	ErrInvalidPrice  = errors.New("invalid price: price cannot be negative")
	ErrInvalidVolume = errors.New("invalid volume: volume cannot be negative")
	ErrStockNotFound = errors.New("stock not found")

	// Recommendation errors
	ErrInvalidScore           = errors.New("invalid score: score must be between 0 and 100")
	ErrInvalidReason          = errors.New("invalid reason: reason cannot be empty")
	ErrNoStocksToRecommend    = errors.New("no stocks available to recommend")
	ErrRecommendationNotFound = errors.New("recommendation not found")
)
