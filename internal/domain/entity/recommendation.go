package entity

import (
	"time"
)

// RecommendationType define los tipos de recomendación disponibles
type RecommendationType string

const (
	RecommendationBuy       RecommendationType = "BUY"
	RecommendationHold      RecommendationType = "HOLD"
	RecommendationSell      RecommendationType = "SELL"
	RecommendationStrongBuy RecommendationType = "STRONG_BUY"
	
	// Aliases para compatibilidad
	RecommendationTypeBuy    = RecommendationBuy
	RecommendationTypeHold   = RecommendationHold
	RecommendationTypeSell   = RecommendationSell
	RecommendationTypeStrong = RecommendationStrongBuy
)

// StrategyType define los tipos de estrategia de análisis
type StrategyType string

const (
	StrategyBalanced StrategyType = "BALANCED"
	StrategyMomentum StrategyType = "MOMENTUM"
	StrategyValue    StrategyType = "VALUE"
	StrategyDividend StrategyType = "DIVIDEND"
	StrategyGrowth   StrategyType = "GROWTH"
)

// Recommendation representa una recomendación de inversión para un stock
// Contiene el análisis y score calculado por el algoritmo
type Recommendation struct {
	ID          string             `json:"id"`
	StockID     string             `json:"stock_id"`
	StockSymbol string             `json:"stock_symbol"`
	Type        RecommendationType `json:"type"`
	Score       float64            `json:"score"`        // Score de 0-100
	Confidence  float64            `json:"confidence"`   // Nivel de confianza 0-1
	Reason      string             `json:"reason"`       // Razón de la recomendación
	Metrics     map[string]float64 `json:"metrics"`      // Métricas usadas en el cálculo
	Strategy    string             `json:"strategy"`     // Estrategia utilizada
	ValidUntil  time.Time          `json:"valid_until"`  // Validez de la recomendación
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// NewRecommendation crea una nueva recomendación con validaciones
func NewRecommendation(stockID, stockSymbol, strategy string, score float64) (*Recommendation, error) {
	if stockID == "" {
		return nil, ErrStockNotFound
	}
	if score < 0 || score > 100 {
		return nil, ErrInvalidScore
	}

	now := time.Now()
	recType := determineRecommendationType(score)
	confidence := calculateConfidence(score)

	return &Recommendation{
		StockID:     stockID,
		StockSymbol: stockSymbol,
		Type:        recType,
		Score:       score,
		Confidence:  confidence,
		Strategy:    strategy,
		ValidUntil:  now.Add(24 * time.Hour), // Válida por 24 horas
		CreatedAt:   now,
		UpdatedAt:   now,
		Metrics:     make(map[string]float64),
	}, nil
}

// determineRecommendationType determina el tipo de recomendación basado en el score
// Lógica de negocio: score > 80 = STRONG_BUY, > 60 = BUY, > 40 = HOLD, <= 40 = SELL
func determineRecommendationType(score float64) RecommendationType {
	switch {
	case score >= 80:
		return RecommendationStrongBuy
	case score >= 60:
		return RecommendationBuy
	case score >= 40:
		return RecommendationHold
	default:
		return RecommendationSell
	}
}

// calculateConfidence calcula el nivel de confianza basado en qué tan cerca está el score de los umbrales
func calculateConfidence(score float64) float64 {
	// Mapeo simple: scores cercanos a umbrales tienen menor confianza
	thresholds := []float64{40, 60, 80}
	minDistance := 100.0

	for _, threshold := range thresholds {
		distance := abs(score - threshold)
		if distance < minDistance {
			minDistance = distance
		}
	}

	// Confianza máxima cuando está lejos de umbrales (distancia > 15)
	if minDistance >= 15 {
		return 0.95
	}
	// Confianza mínima cuando está muy cerca de umbrales (distancia < 5)
	if minDistance <= 5 {
		return 0.60
	}
	// Interpolar linealmente entre 5 y 15
	return 0.60 + ((minDistance-5)/10)*0.35
}

// AddMetric agrega una métrica al cálculo de la recomendación
func (r *Recommendation) AddMetric(name string, value float64) {
	if r.Metrics == nil {
		r.Metrics = make(map[string]float64)
	}
	r.Metrics[name] = value
}

// SetConfidence establece el nivel de confianza de la recomendación
func (r *Recommendation) SetConfidence(confidence float64) {
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}
	r.Confidence = confidence
	r.UpdatedAt = time.Now()
}

// SetReason establece la razón de la recomendación
func (r *Recommendation) SetReason(reason string) error {
	if reason == "" {
		return ErrInvalidReason
	}
	r.Reason = reason
	r.UpdatedAt = time.Now()
	return nil
}

// IsValid verifica si la recomendación aún es válida (no ha expirado)
func (r *Recommendation) IsValid() bool {
	return time.Now().Before(r.ValidUntil)
}

// IsBuyRecommendation verifica si es una recomendación de compra
func (r *Recommendation) IsBuyRecommendation() bool {
	return r.Type == RecommendationBuy || r.Type == RecommendationStrongBuy
}

// Validate verifica que la recomendación tenga datos válidos
func (r *Recommendation) Validate() error {
	if r.StockID == "" {
		return ErrStockNotFound
	}
	if r.Score < 0 || r.Score > 100 {
		return ErrInvalidScore
	}
	if r.Reason == "" {
		return ErrInvalidReason
	}
	return nil
}

// abs retorna el valor absoluto de un float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
