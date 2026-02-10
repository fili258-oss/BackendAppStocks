package service

import (
	"context"
	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// RecommendationStrategy define la interfaz para diferentes algoritmos de recomendación
// Strategy Pattern: permite intercambiar algoritmos de recomendación sin modificar el código cliente
// Open/Closed Principle: abierto para extensión (nuevas estrategias), cerrado para modificación
type RecommendationStrategy interface {
	// Analyze analiza un stock y retorna una recomendación
	Analyze(ctx context.Context, stock *entity.Stock) (*entity.Recommendation, error)

	// AnalyzeBatch analiza múltiples stocks y retorna recomendaciones
	// Útil para análisis masivo y comparaciones
	AnalyzeBatch(ctx context.Context, stocks []*entity.Stock) ([]*entity.Recommendation, error)

	// GetName retorna el nombre de la estrategia
	GetName() string

	// GetDescription retorna una descripción de cómo funciona la estrategia
	GetDescription() string
}

// RecommendationScorer define la interfaz para calcular scores de stocks
// Permite diferentes sistemas de puntuación
type RecommendationScorer interface {
	// CalculateScore calcula un score de 0-100 para un stock
	CalculateScore(stock *entity.Stock) float64

	// GetWeights retorna los pesos usados en el cálculo
	GetWeights() map[string]float64
}

// StrategyType define los tipos de estrategias disponibles
type StrategyType string

const (
	// StrategyTypeValueInvesting estrategia basada en value investing
	StrategyTypeValueInvesting StrategyType = "VALUE_INVESTING"

	// StrategyTypeMomentum estrategia basada en momentum de precio
	StrategyTypeMomentum StrategyType = "MOMENTUM"

	// StrategyTypeDividend estrategia basada en dividendos
	StrategyTypeDividend StrategyType = "DIVIDEND"

	// StrategyTypeGrowth estrategia basada en crecimiento
	StrategyTypeGrowth StrategyType = "GROWTH"

	// StrategyTypeBalanced estrategia balanceada que combina múltiples factores
	StrategyTypeBalanced StrategyType = "BALANCED"
)
