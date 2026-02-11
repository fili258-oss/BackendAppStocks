package service

import (
	"context"
	"fmt"

	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// DividendStrategy implementa análisis enfocado en dividendos
// Ideal para inversores que buscan ingresos pasivos
type DividendStrategy struct {
	scorer *BaseScorer
}

// NewDividendStrategy crea una nueva instancia
func NewDividendStrategy() *DividendStrategy {
	return &DividendStrategy{
		scorer: NewBaseScorer(),
	}
}

// GetName retorna el nombre de la estrategia
func (s *DividendStrategy) GetName() string {
	return string(entity.StrategyDividend)
}

// GetDescription retorna la descripción
func (s *DividendStrategy) GetDescription() string {
	return "Análisis enfocado en generación de ingresos pasivos mediante dividendos"
}

// Analyze analiza un stock
func (s *DividendStrategy) Analyze(ctx context.Context, stock *entity.Stock) (*entity.Recommendation, error) {
	metrics := s.extractMetrics(stock)

	// Pesos enfocados en dividendos e estabilidad
	weights := map[string]float64{
		"dividend_yield_score": 0.40, // 40% - Dividend yield
		"company_stability":    0.25, // 25% - Estabilidad (market cap)
		"price_stability":      0.20, // 20% - Estabilidad de precio
		"pe_sustainability":    0.15, // 15% - P/E ratio (sostenibilidad)
	}

	score := s.scorer.CalculateScore(metrics, weights)
	
	requiredMetrics := []string{"dividend_yield_score", "company_stability"}
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
func (s *DividendStrategy) AnalyzeBatch(ctx context.Context, stocks []*entity.Stock) ([]*entity.Recommendation, error) {
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

// extractMetrics extrae métricas de dividendos
func (s *DividendStrategy) extractMetrics(stock *entity.Stock) map[string]float64 {
	metrics := make(map[string]float64)

	// Dividend Yield (crítico)
	if stock.DividendYield > 0 {
		// > 5% = excelente (1.0)
		// 3-5% = muy bueno (0.7-1.0)
		// 2-3% = bueno (0.5-0.7)
		// < 2% = bajo (0-0.5)
		var divScore float64
		if stock.DividendYield >= 5 {
			divScore = 1.0
		} else if stock.DividendYield >= 3 {
			divScore = 0.7 + ((stock.DividendYield - 3) / 2) * 0.3
		} else if stock.DividendYield >= 2 {
			divScore = 0.5 + ((stock.DividendYield - 2) / 1) * 0.2
		} else {
			divScore = (stock.DividendYield / 2) * 0.5
		}
		metrics["dividend_yield_score"] = divScore
	} else {
		// Sin dividendos = score muy bajo
		metrics["dividend_yield_score"] = 0.1
	}

	// Estabilidad de la compañía (market cap)
	if stock.MarketCap > 0 {
		// Para dividendos, preferimos compañías grandes y estables
		// > $100B = 1.0
		// $50B-$100B = 0.8-1.0
		// $20B-$50B = 0.5-0.8
		// < $20B = 0-0.5
		var stabScore float64
		if stock.MarketCap >= 100000000000 {
			stabScore = 1.0
		} else if stock.MarketCap >= 50000000000 {
			stabScore = 0.8 + (stock.MarketCap-50000000000)/(100000000000-50000000000)*0.2
		} else if stock.MarketCap >= 20000000000 {
			stabScore = 0.5 + (stock.MarketCap-20000000000)/(50000000000-20000000000)*0.3
		} else {
			stabScore = (stock.MarketCap / 20000000000) * 0.5
		}
		metrics["company_stability"] = stabScore
	}

	// Estabilidad de precio (baja volatilidad)
	if stock.Week52High > 0 && stock.Week52Low > 0 {
		volatilityRange := (stock.Week52High - stock.Week52Low) / stock.Week52Low
		// Menor volatilidad = mejor para dividendos
		// < 20% = excelente (1.0)
		// 20-40% = bueno (0.7-1.0)
		// > 40% = alto riesgo (0-0.7)
		var priceStab float64
		if volatilityRange <= 0.2 {
			priceStab = 1.0
		} else if volatilityRange <= 0.4 {
			priceStab = 1.0 - ((volatilityRange - 0.2) / 0.2) * 0.3
		} else {
			priceStab = 0.7 - ((volatilityRange - 0.4) / 0.6) * 0.7
			if priceStab < 0 {
				priceStab = 0
			}
		}
		metrics["price_stability"] = priceStab
	}

	// P/E sostenible (ni muy alto ni muy bajo)
	if stock.PERatio > 0 {
		// P/E 10-20 = óptimo para dividendos (1.0)
		// P/E < 10 = posible distress (0.5)
		// P/E > 20 = sobrevalorado (0.5)
		var peScore float64
		if stock.PERatio >= 10 && stock.PERatio <= 20 {
			peScore = 1.0
		} else if stock.PERatio < 10 {
			peScore = 0.5 + (stock.PERatio / 10) * 0.5
		} else {
			peScore = 1.0 - ((stock.PERatio - 20) / 30) * 0.5
			if peScore < 0.3 {
				peScore = 0.3
			}
		}
		metrics["pe_sustainability"] = peScore
	}

	return metrics
}

// generateHighlights genera puntos destacados
func (s *DividendStrategy) generateHighlights(stock *entity.Stock, metrics map[string]float64) []string {
	highlights := make([]string, 0)

	// Dividend yield destacado
	if stock.DividendYield >= 4 {
		highlights = append(highlights, fmt.Sprintf("excelente dividend yield (%.2f%%)", stock.DividendYield))
	} else if stock.DividendYield >= 2.5 {
		highlights = append(highlights, fmt.Sprintf("buen dividend yield (%.2f%%)", stock.DividendYield))
	} else if stock.DividendYield > 0 {
		highlights = append(highlights, fmt.Sprintf("dividend yield moderado (%.2f%%)", stock.DividendYield))
	}

	// Estabilidad
	if stock.MarketCap > 100000000000 {
		highlights = append(highlights, "compañía muy estable (mega cap)")
	} else if stock.MarketCap > 50000000000 {
		highlights = append(highlights, "compañía estable (large cap)")
	}

	// Baja volatilidad
	if priceStab, exists := metrics["price_stability"]; exists && priceStab > 0.8 {
		highlights = append(highlights, "baja volatilidad de precio")
	}

	return highlights
}
