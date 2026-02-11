package service

import (
	"context"
	"fmt"

	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// MomentumStrategy implementa análisis basado en momentum
// Se enfoca en tendencias de precio y volumen
type MomentumStrategy struct {
	scorer *BaseScorer
}

// NewMomentumStrategy crea una nueva instancia
func NewMomentumStrategy() *MomentumStrategy {
	return &MomentumStrategy{
		scorer: NewBaseScorer(),
	}
}

// GetName retorna el nombre de la estrategia
func (s *MomentumStrategy) GetName() string {
	return string(entity.StrategyMomentum)
}

// GetDescription retorna la descripción
func (s *MomentumStrategy) GetDescription() string {
	return "Análisis basado en tendencias de precio y momentum del mercado"
}

// Analyze analiza un stock
func (s *MomentumStrategy) Analyze(ctx context.Context, stock *entity.Stock) (*entity.Recommendation, error) {
	metrics := s.extractMetrics(stock)

	// Pesos enfocados en momentum
	weights := map[string]float64{
		"price_change":       0.35, // 35% - Cambio de precio
		"price_trend":        0.25, // 25% - Tendencia (vs 52-week)
		"volume_momentum":    0.25, // 25% - Volumen
		"volatility_score":   0.15, // 15% - Volatilidad favorable
	}

	score := s.scorer.CalculateScore(metrics, weights)
	
	requiredMetrics := []string{"price_change", "price_trend", "volume_momentum"}
	confidence := s.scorer.CalculateConfidence(metrics, requiredMetrics)

	recType := s.scorer.DetermineRecommendationType(score)
	highlights := s.generateHighlights(stock, metrics)
	reason := s.scorer.GenerateReason(recType, s.GetName(), highlights)

	recommendation, err := entity.NewRecommendation(
		stock.ID,
		stock.Symbol,
		s.GetName(),
		score,
	)
	if err != nil {
		return nil, err
	}

	recommendation.SetConfidence(confidence)
	recommendation.SetReason(reason)

	for key, value := range metrics {
		recommendation.AddMetric(key, value)
	}

	return recommendation, nil
}

// AnalyzeBatch analiza múltiples stocks
func (s *MomentumStrategy) AnalyzeBatch(ctx context.Context, stocks []*entity.Stock) ([]*entity.Recommendation, error) {
	recommendations := make([]*entity.Recommendation, 0, len(stocks))

	for _, stock := range stocks {
		rec, err := s.Analyze(ctx, stock)
		if err != nil {
			continue
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

// extractMetrics extrae métricas de momentum
func (s *MomentumStrategy) extractMetrics(stock *entity.Stock) map[string]float64 {
	metrics := make(map[string]float64)

	// Cambio de precio (más agresivo que balanced)
	if stock.ChangePercent != 0 {
		// +10% = 1.0, -10% = 0.0
		changeScore := (stock.ChangePercent + 10) / 20
		if changeScore < 0 {
			changeScore = 0
		}
		if changeScore > 1 {
			changeScore = 1
		}
		metrics["price_change"] = changeScore
	}

	// Tendencia (posición en rango 52-week)
	if stock.Week52High > 0 && stock.Week52Low > 0 && stock.Price > 0 {
		// Precio cerca del máximo = bueno
		position := (stock.Price - stock.Week52Low) / (stock.Week52High - stock.Week52Low)
		metrics["price_trend"] = position
	}

	// Volume momentum
	if stock.Volume > 0 {
		// Volumen muy alto = fuerte momentum
		volumeScore := float64(stock.Volume) / 100000000 // 100M = 1.0
		if volumeScore > 1 {
			volumeScore = 1
		}
		metrics["volume_momentum"] = volumeScore
	}

	// Volatilidad como oportunidad
	if stock.Week52High > 0 && stock.Week52Low > 0 {
		volatilityRange := (stock.Week52High - stock.Week52Low) / stock.Week52Low
		// Para momentum, cierta volatilidad es buena
		// 50% range = 1.0
		volScore := volatilityRange / 0.5
		if volScore > 1 {
			volScore = 1
		}
		metrics["volatility_score"] = volScore
	}

	return metrics
}

// generateHighlights genera puntos destacados
func (s *MomentumStrategy) generateHighlights(stock *entity.Stock, metrics map[string]float64) []string {
	highlights := make([]string, 0)

	// Momentum fuerte
	if stock.ChangePercent > 5 {
		highlights = append(highlights, fmt.Sprintf("fuerte momentum alcista (+%.2f%%)", stock.ChangePercent))
	} else if stock.ChangePercent < -5 {
		highlights = append(highlights, fmt.Sprintf("momentum bajista (%.2f%%)", stock.ChangePercent))
	}

	// Posición en rango
	if trendScore, exists := metrics["price_trend"]; exists {
		if trendScore > 0.8 {
			highlights = append(highlights, "cerca del máximo de 52 semanas")
		} else if trendScore < 0.2 {
			highlights = append(highlights, "cerca del mínimo de 52 semanas")
		}
	}

	// Volumen excepcional
	if stock.Volume > 50000000 {
		highlights = append(highlights, "volumen excepcional")
	}

	return highlights
}
