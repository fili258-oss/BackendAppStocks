package service

import (
	"context"
	"fmt"

	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// GrowthStrategy implementa análisis de growth investing
// Se enfoca en potencial de crecimiento y momentum alcista
type GrowthStrategy struct {
	scorer *BaseScorer
}

// NewGrowthStrategy crea una nueva instancia
func NewGrowthStrategy() *GrowthStrategy {
	return &GrowthStrategy{
		scorer: NewBaseScorer(),
	}
}

// GetName retorna el nombre de la estrategia
func (s *GrowthStrategy) GetName() string {
	return string(entity.StrategyGrowth)
}

// GetDescription retorna la descripción
func (s *GrowthStrategy) GetDescription() string {
	return "Análisis de growth investing enfocado en potencial de crecimiento"
}

// Analyze analiza un stock
func (s *GrowthStrategy) Analyze(ctx context.Context, stock *entity.Stock) (*entity.Recommendation, error) {
	metrics := s.extractMetrics(stock)

	// Pesos enfocados en crecimiento
	weights := map[string]float64{
		"price_momentum":     0.30, // 30% - Momentum alcista
		"uptrend_strength":   0.25, // 25% - Fuerza de tendencia
		"volume_growth":      0.20, // 20% - Crecimiento de volumen
		"market_interest":    0.15, // 15% - Interés del mercado
		"growth_potential":   0.10, // 10% - Potencial (basado en P/E)
	}

	score := s.scorer.CalculateScore(metrics, weights)
	
	requiredMetrics := []string{"price_momentum", "uptrend_strength"}
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
func (s *GrowthStrategy) AnalyzeBatch(ctx context.Context, stocks []*entity.Stock) ([]*entity.Recommendation, error) {
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

// extractMetrics extrae métricas de growth
func (s *GrowthStrategy) extractMetrics(stock *entity.Stock) map[string]float64 {
	metrics := make(map[string]float64)

	// Momentum de precio (agresivo)
	if stock.ChangePercent != 0 {
		// +15% = 1.0, 0% = 0.5, -15% = 0.0
		momentum := (stock.ChangePercent + 15) / 30
		if momentum < 0 {
			momentum = 0
		}
		if momentum > 1 {
			momentum = 1
		}
		metrics["price_momentum"] = momentum
	}

	// Fuerza de tendencia alcista (posición en rango 52-week)
	if stock.Week52High > 0 && stock.Week52Low > 0 && stock.Price > 0 {
		// Para growth, cerca del máximo es excelente
		position := (stock.Price - stock.Week52Low) / (stock.Week52High - stock.Week52Low)
		// Bonus si está en el top 20%
		if position > 0.8 {
			position = 0.8 + (position - 0.8) * 2 // Amplificar
			if position > 1 {
				position = 1
			}
		}
		metrics["uptrend_strength"] = position
	}

	// Crecimiento de volumen (volumen alto = interés)
	if stock.Volume > 0 {
		// Volumen muy alto indica fuerte interés
		// 50M+ = 1.0
		volGrowth := float64(stock.Volume) / 50000000
		if volGrowth > 1 {
			volGrowth = 1
		}
		metrics["volume_growth"] = volGrowth
	}

	// Interés del mercado (combinación de factores)
	if stock.MarketCap > 0 && stock.Volume > 0 {
		// Balance entre tamaño y actividad
		// Mid-large caps con alto volumen = óptimo
		capFactor := 0.5
		if stock.MarketCap >= 10000000000 && stock.MarketCap <= 500000000000 {
			capFactor = 1.0
		} else if stock.MarketCap > 500000000000 {
			capFactor = 0.7
		} else {
			capFactor = stock.MarketCap / 10000000000 * 0.5
		}

		volFactor := float64(stock.Volume) / 30000000
		if volFactor > 1 {
			volFactor = 1
		}

		metrics["market_interest"] = (capFactor + volFactor) / 2
	}

	// Potencial de crecimiento (P/E alto puede ser bueno para growth)
	if stock.PERatio > 0 {
		// Para growth, P/E alto puede indicar expectativas altas
		// P/E 20-40 = óptimo (1.0)
		// P/E < 20 = puede no ser growth (0.5)
		// P/E > 40 = sobrevalorado (0.5)
		var peScore float64
		if stock.PERatio >= 20 && stock.PERatio <= 40 {
			peScore = 0.7 + (1.0 - abs(stock.PERatio-30)/10) * 0.3
		} else if stock.PERatio < 20 {
			peScore = 0.5 + (stock.PERatio / 20) * 0.2
		} else {
			peScore = 1.0 - ((stock.PERatio - 40) / 60) * 0.5
			if peScore < 0.3 {
				peScore = 0.3
			}
		}
		metrics["growth_potential"] = peScore
	}

	return metrics
}

// abs retorna el valor absoluto
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// generateHighlights genera puntos destacados
func (s *GrowthStrategy) generateHighlights(stock *entity.Stock, metrics map[string]float64) []string {
	highlights := make([]string, 0)

	// Momentum fuerte
	if stock.ChangePercent > 10 {
		highlights = append(highlights, fmt.Sprintf("momentum excepcional (+%.2f%%)", stock.ChangePercent))
	} else if stock.ChangePercent > 5 {
		highlights = append(highlights, fmt.Sprintf("fuerte momentum alcista (+%.2f%%)", stock.ChangePercent))
	}

	// En máximos
	if trendScore, exists := metrics["uptrend_strength"]; exists && trendScore > 0.85 {
		highlights = append(highlights, "en máximos de 52 semanas")
	} else if exists && trendScore > 0.7 {
		highlights = append(highlights, "cerca de máximos históricos")
	}

	// Alto volumen
	if stock.Volume > 50000000 {
		highlights = append(highlights, "volumen excepcional")
	} else if stock.Volume > 30000000 {
		highlights = append(highlights, "alto volumen de trading")
	}

	// P/E de growth
	if stock.PERatio > 25 && stock.PERatio < 45 {
		highlights = append(highlights, "valuación de growth stock")
	}

	return highlights
}
