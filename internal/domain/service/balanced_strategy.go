package service

import (
	"context"
	"fmt"

	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// BalancedStrategy implementa una estrategia de análisis equilibrado
// Considera múltiples factores: precio, momentum, fundamentales
type BalancedStrategy struct {
	scorer *BaseScorer
}

// NewBalancedStrategy crea una nueva instancia
func NewBalancedStrategy() *BalancedStrategy {
	return &BalancedStrategy{
		scorer: NewBaseScorer(),
	}
}

// GetName retorna el nombre de la estrategia
func (s *BalancedStrategy) GetName() string {
	return string(entity.StrategyBalanced)
}

// GetDescription retorna la descripción de la estrategia
func (s *BalancedStrategy) GetDescription() string {
	return "Análisis equilibrado que considera precio, momentum y fundamentales"
}

// Analyze analiza un stock y genera una recomendación
func (s *BalancedStrategy) Analyze(ctx context.Context, stock *entity.Stock) (*entity.Recommendation, error) {
	// Extraer métricas del stock
	metrics := s.extractMetrics(stock)

	// Definir pesos para cada métrica
	weights := map[string]float64{
		"price_momentum":  0.20, // 20% - Cambio de precio
		"volume":          0.15, // 15% - Volumen
		"market_cap":      0.15, // 15% - Capitalización
		"pe_ratio":        0.20, // 20% - P/E ratio
		"dividend_yield":  0.10, // 10% - Dividendos
		"volatility":      0.20, // 20% - Volatilidad (inversa)
	}

	// Calcular score
	score := s.scorer.CalculateScore(metrics, weights)

	// Calcular confianza
	requiredMetrics := []string{"price_momentum", "volume", "market_cap", "pe_ratio"}
	confidence := s.scorer.CalculateConfidence(metrics, requiredMetrics)

	// Determinar tipo de recomendación
	recType := s.scorer.DetermineRecommendationType(score)

	// Generar highlights
	highlights := s.generateHighlights(stock, metrics)

	// Generar razón
	reason := s.scorer.GenerateReason(recType, s.GetName(), highlights)

	// Crear recomendación
	recommendation, err := entity.NewRecommendation(
		stock.ID,
		stock.Symbol,
		s.GetName(),
		score,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating recommendation: %w", err)
	}

	recommendation.SetConfidence(confidence)
	recommendation.SetReason(reason)

	// Agregar métricas usadas
	for key, value := range metrics {
		recommendation.AddMetric(key, value)
	}

	return recommendation, nil
}

// AnalyzeBatch analiza múltiples stocks
func (s *BalancedStrategy) AnalyzeBatch(ctx context.Context, stocks []*entity.Stock) ([]*entity.Recommendation, error) {
	recommendations := make([]*entity.Recommendation, 0, len(stocks))

	for _, stock := range stocks {
		rec, err := s.Analyze(ctx, stock)
		if err != nil {
			// Continuar con los demás si uno falla
			continue
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

// extractMetrics extrae las métricas relevantes del stock
func (s *BalancedStrategy) extractMetrics(stock *entity.Stock) map[string]float64 {
	metrics := make(map[string]float64)

	// Momentum del precio (normalizado)
	if stock.ChangePercent != 0 {
		// Normalizar: +5% = 1.0, -5% = 0.0, 0% = 0.5
		momentum := (stock.ChangePercent + 5) / 10
		if momentum < 0 {
			momentum = 0
		}
		if momentum > 1 {
			momentum = 1
		}
		metrics["price_momentum"] = momentum
	}

	// Volumen (normalizado)
	if stock.Volume > 0 {
		// Volumen alto = bueno
		// 10M = 0.5, 50M = 1.0
		volumeScore := float64(stock.Volume) / 50000000
		if volumeScore > 1 {
			volumeScore = 1
		}
		metrics["volume"] = volumeScore
	}

	// Market Cap (normalizado)
	if stock.MarketCap > 0 {
		// $100B = 0.5, $2T = 1.0
		capScore := stock.MarketCap / 2000000000000
		if capScore > 1 {
			capScore = 1
		}
		metrics["market_cap"] = capScore
	}

	// P/E Ratio (inverso - menor es mejor)
	if stock.PERatio > 0 {
		// P/E de 15 = 1.0, P/E de 30 = 0.5, P/E de 50+ = 0
		peScore := 1.0 - ((stock.PERatio - 15) / 35)
		if peScore < 0 {
			peScore = 0
		}
		if peScore > 1 {
			peScore = 1
		}
		metrics["pe_ratio"] = peScore
	}

	// Dividend Yield
	if stock.DividendYield > 0 {
		// 2% = 0.5, 4% = 1.0
		divScore := stock.DividendYield / 4
		if divScore > 1 {
			divScore = 1
		}
		metrics["dividend_yield"] = divScore
	}

	// Volatilidad (basada en 52-week range)
	if stock.Week52High > 0 && stock.Week52Low > 0 {
		volatilityRange := (stock.Week52High - stock.Week52Low) / stock.Week52Low
		// Volatilidad baja = bueno
		// 20% range = 1.0, 100% range = 0.0
		volScore := 1.0 - (volatilityRange / 1.0)
		if volScore < 0 {
			volScore = 0
		}
		if volScore > 1 {
			volScore = 1
		}
		metrics["volatility"] = volScore
	}

	return metrics
}

// generateHighlights genera puntos destacados del análisis
func (s *BalancedStrategy) generateHighlights(stock *entity.Stock, metrics map[string]float64) []string {
	highlights := make([]string, 0)

	// Momentum positivo
	if stock.ChangePercent > 2 {
		highlights = append(highlights, fmt.Sprintf("momentum positivo (+%.2f%%)", stock.ChangePercent))
	} else if stock.ChangePercent < -2 {
		highlights = append(highlights, fmt.Sprintf("momentum negativo (%.2f%%)", stock.ChangePercent))
	}

	// Volumen alto
	if stock.Volume > 30000000 {
		highlights = append(highlights, "alto volumen de trading")
	}

	// P/E atractivo
	if stock.PERatio > 0 && stock.PERatio < 20 {
		highlights = append(highlights, fmt.Sprintf("P/E atractivo (%.2f)", stock.PERatio))
	}

	// Dividendos
	if stock.DividendYield > 2 {
		highlights = append(highlights, fmt.Sprintf("buenos dividendos (%.2f%%)", stock.DividendYield))
	}

	// Market cap
	if stock.MarketCap > 1000000000000 {
		highlights = append(highlights, "gran capitalización")
	}

	return highlights
}
