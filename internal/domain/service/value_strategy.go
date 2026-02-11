package service

import (
	"context"
	"fmt"

	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// ValueStrategy implementa análisis de value investing
// Se enfoca en fundamentales y valoración
type ValueStrategy struct {
	scorer *BaseScorer
}

// NewValueStrategy crea una nueva instancia
func NewValueStrategy() *ValueStrategy {
	return &ValueStrategy{
		scorer: NewBaseScorer(),
	}
}

// GetName retorna el nombre de la estrategia
func (s *ValueStrategy) GetName() string {
	return string(entity.StrategyValue)
}

// GetDescription retorna la descripción
func (s *ValueStrategy) GetDescription() string {
	return "Análisis de value investing enfocado en fundamentales y valoración"
}

// Analyze analiza un stock
func (s *ValueStrategy) Analyze(ctx context.Context, stock *entity.Stock) (*entity.Recommendation, error) {
	metrics := s.extractMetrics(stock)

	// Pesos enfocados en value
	weights := map[string]float64{
		"pe_ratio_score":     0.30, // 30% - P/E ratio
		"market_cap_score":   0.20, // 20% - Capitalización (estabilidad)
		"dividend_score":     0.20, // 20% - Dividendos
		"price_position":     0.20, // 20% - Posición en rango 52-week
		"volume_stability":   0.10, // 10% - Estabilidad de volumen
	}

	score := s.scorer.CalculateScore(metrics, weights)
	
	requiredMetrics := []string{"pe_ratio_score", "market_cap_score"}
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
func (s *ValueStrategy) AnalyzeBatch(ctx context.Context, stocks []*entity.Stock) ([]*entity.Recommendation, error) {
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

// extractMetrics extrae métricas de value
func (s *ValueStrategy) extractMetrics(stock *entity.Stock) map[string]float64 {
	metrics := make(map[string]float64)

	// P/E Ratio (más estricto para value)
	if stock.PERatio > 0 {
		// P/E < 15 = excelente (1.0)
		// P/E 15-25 = bueno (0.5-1.0)
		// P/E > 25 = caro (0-0.5)
		var peScore float64
		if stock.PERatio <= 15 {
			peScore = 1.0
		} else if stock.PERatio <= 25 {
			peScore = 1.0 - ((stock.PERatio - 15) / 10) * 0.5
		} else {
			peScore = 0.5 - ((stock.PERatio - 25) / 50) * 0.5
			if peScore < 0 {
				peScore = 0
			}
		}
		metrics["pe_ratio_score"] = peScore
	}

	// Market Cap (preferencia por large caps estables)
	if stock.MarketCap > 0 {
		// > $200B = 1.0 (mega cap)
		// $50B-$200B = 0.7-1.0 (large cap)
		// $10B-$50B = 0.4-0.7 (mid cap)
		// < $10B = 0-0.4 (small cap)
		var capScore float64
		if stock.MarketCap >= 200000000000 {
			capScore = 1.0
		} else if stock.MarketCap >= 50000000000 {
			capScore = 0.7 + (stock.MarketCap-50000000000)/(200000000000-50000000000)*0.3
		} else if stock.MarketCap >= 10000000000 {
			capScore = 0.4 + (stock.MarketCap-10000000000)/(50000000000-10000000000)*0.3
		} else {
			capScore = (stock.MarketCap / 10000000000) * 0.4
		}
		metrics["market_cap_score"] = capScore
	}

	// Dividend Yield (importante para value)
	if stock.DividendYield > 0 {
		// > 4% = 1.0
		// 2-4% = 0.5-1.0
		// < 2% = 0-0.5
		divScore := stock.DividendYield / 4
		if divScore > 1 {
			divScore = 1
		}
		metrics["dividend_score"] = divScore
	}

	// Posición en rango 52-week (para value, más bajo es mejor)
	if stock.Week52High > 0 && stock.Week52Low > 0 && stock.Price > 0 {
		position := (stock.Price - stock.Week52Low) / (stock.Week52High - stock.Week52Low)
		// Invertir: cerca del mínimo = mejor oportunidad
		priceScore := 1.0 - position
		metrics["price_position"] = priceScore
	}

	// Volumen estable
	if stock.Volume > 0 {
		// Volumen moderado y estable es bueno para value
		// 1M-20M = óptimo (1.0)
		var volScore float64
		if stock.Volume >= 1000000 && stock.Volume <= 20000000 {
			volScore = 1.0
		} else if stock.Volume < 1000000 {
			volScore = float64(stock.Volume) / 1000000
		} else {
			volScore = 1.0 - ((float64(stock.Volume) - 20000000) / 80000000)
			if volScore < 0.5 {
				volScore = 0.5
			}
		}
		metrics["volume_stability"] = volScore
	}

	return metrics
}

// generateHighlights genera puntos destacados
func (s *ValueStrategy) generateHighlights(stock *entity.Stock, metrics map[string]float64) []string {
	highlights := make([]string, 0)

	// P/E atractivo
	if stock.PERatio > 0 && stock.PERatio < 15 {
		highlights = append(highlights, fmt.Sprintf("P/E muy atractivo (%.2f)", stock.PERatio))
	} else if stock.PERatio > 0 && stock.PERatio < 20 {
		highlights = append(highlights, fmt.Sprintf("P/E razonable (%.2f)", stock.PERatio))
	}

	// Dividendos sólidos
	if stock.DividendYield >= 3 {
		highlights = append(highlights, fmt.Sprintf("dividendos sólidos (%.2f%%)", stock.DividendYield))
	}

	// Large cap estable
	if stock.MarketCap > 100000000000 {
		highlights = append(highlights, "compañía de gran capitalización")
	}

	// Precio cerca del mínimo
	if pricePos, exists := metrics["price_position"]; exists && pricePos > 0.7 {
		highlights = append(highlights, "precio cerca del mínimo de 52 semanas")
	}

	return highlights
}
