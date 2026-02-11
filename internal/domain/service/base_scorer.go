package service

import (
	"fmt"

	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// BaseScorer implementa la lógica común de scoring
type BaseScorer struct{}

// NewBaseScorer crea una nueva instancia
func NewBaseScorer() *BaseScorer {
	return &BaseScorer{}
}

// CalculateScore calcula un score normalizado de 0 a 100
func (s *BaseScorer) CalculateScore(metrics map[string]float64, weights map[string]float64) float64 {
	totalWeight := 0.0
	weightedSum := 0.0

	for metric, value := range metrics {
		weight, exists := weights[metric]
		if !exists {
			continue
		}

		// Normalizar el valor entre 0 y 1
		normalizedValue := s.normalizeValue(value)
		
		weightedSum += normalizedValue * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 50.0 // Score neutral si no hay métricas
	}

	// Convertir a escala 0-100
	score := (weightedSum / totalWeight) * 100

	// Asegurar que esté en rango
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// normalizeValue normaliza un valor a rango 0-1
func (s *BaseScorer) normalizeValue(value float64) float64 {
	// Para valores ya normalizados (0-1)
	if value >= 0 && value <= 1 {
		return value
	}

	// Para valores positivos grandes, usar función sigmoide
	if value > 1 {
		return 1 / (1 + 1/value)
	}

	// Para valores negativos, invertir
	if value < 0 {
		return 0
	}

	return 0.5
}

// CalculateConfidence calcula el nivel de confianza basado en disponibilidad de datos
func (s *BaseScorer) CalculateConfidence(metrics map[string]float64, requiredMetrics []string) float64 {
	if len(requiredMetrics) == 0 {
		return 1.0
	}

	available := 0
	for _, metric := range requiredMetrics {
		if value, exists := metrics[metric]; exists && value != 0 {
			available++
		}
	}

	confidence := float64(available) / float64(len(requiredMetrics))
	return confidence
}

// DetermineRecommendationType determina el tipo de recomendación basado en score
func (s *BaseScorer) DetermineRecommendationType(score float64) entity.RecommendationType {
	switch {
	case score >= 80:
		return entity.RecommendationStrongBuy
	case score >= 60:
		return entity.RecommendationBuy
	case score >= 40:
		return entity.RecommendationHold
	default:
		return entity.RecommendationSell
	}
}

// GenerateReason genera una razón descriptiva para la recomendación
func (s *BaseScorer) GenerateReason(
	recType entity.RecommendationType,
	strategy string,
	highlights []string,
) string {
	baseReason := fmt.Sprintf("Análisis %s: ", strategy)

	switch recType {
	case entity.RecommendationStrongBuy:
		baseReason += "Excelente oportunidad de compra. "
	case entity.RecommendationBuy:
		baseReason += "Buena oportunidad de compra. "
	case entity.RecommendationHold:
		baseReason += "Se recomienda mantener. "
	case entity.RecommendationSell:
		baseReason += "Se recomienda vender. "
	}

	// Agregar highlights
	if len(highlights) > 0 {
		baseReason += "Factores clave: "
		for i, highlight := range highlights {
			if i > 0 {
				baseReason += ", "
			}
			baseReason += highlight
		}
		baseReason += "."
	}

	return baseReason
}
