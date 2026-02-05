package entity

import (
	"time"
)

// Stock representa una acción en el mercado de valores
// Esta entidad del dominio contiene toda la información relevante de un stock
type Stock struct {
	ID          string    `json:"id"`
	Symbol      string    `json:"symbol"`      // Ticker symbol (e.g., AAPL, GOOGL)
	Name        string    `json:"name"`        // Nombre completo de la compañía
	Exchange    string    `json:"exchange"`    // Bolsa donde se lista (NYSE, NASDAQ, etc.)
	Currency    string    `json:"currency"`    // Moneda de cotización
	Price       float64   `json:"price"`       // Precio actual
	OpenPrice   float64   `json:"open_price"`  // Precio de apertura del día
	HighPrice   float64   `json:"high_price"`  // Precio más alto del día
	LowPrice    float64   `json:"low_price"`   // Precio más bajo del día
	ClosePrice  float64   `json:"close_price"` // Precio de cierre anterior
	Volume      int64     `json:"volume"`      // Volumen de transacciones
	MarketCap   float64   `json:"market_cap"`  // Capitalización de mercado
	PERatio     float64   `json:"pe_ratio"`    // Price-to-Earnings ratio
	DividendYield float64 `json:"dividend_yield"` // Rendimiento por dividendos
	Week52High  float64   `json:"week_52_high"`   // Máximo de 52 semanas
	Week52Low   float64   `json:"week_52_low"`    // Mínimo de 52 semanas
	Change      float64   `json:"change"`         // Cambio en precio
	ChangePercent float64 `json:"change_percent"` // Cambio en porcentaje
	UpdatedAt   time.Time `json:"updated_at"`     // Última actualización
	CreatedAt   time.Time `json:"created_at"`     // Fecha de creación en BD
}

// NewStock crea una nueva instancia de Stock con validaciones básicas
// Constructor factory que asegura la creación de objetos válidos
func NewStock(symbol, name, exchange string) (*Stock, error) {
	if symbol == "" {
		return nil, ErrInvalidSymbol
	}
	if name == "" {
		return nil, ErrInvalidName
	}

	now := time.Now()
	return &Stock{
		Symbol:    symbol,
		Name:      name,
		Exchange:  exchange,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// CalculateChangePercent calcula el porcentaje de cambio basado en el precio anterior
// Business logic para cálculo de variación porcentual
func (s *Stock) CalculateChangePercent() {
	if s.ClosePrice > 0 {
		s.Change = s.Price - s.ClosePrice
		s.ChangePercent = (s.Change / s.ClosePrice) * 100
	}
}

// IsPositiveChange verifica si el stock tiene un cambio positivo
func (s *Stock) IsPositiveChange() bool {
	return s.Change > 0
}

// IsVolatile determina si el stock es volátil basado en el rango de precios
// Considera un stock volátil si la diferencia entre high y low es > 5%
func (s *Stock) IsVolatile() bool {
	if s.LowPrice == 0 {
		return false
	}
	volatilityPercent := ((s.HighPrice - s.LowPrice) / s.LowPrice) * 100
	return volatilityPercent > 5.0
}

// GetPriceToBookRatio calcula el ratio precio/valor libro si es aplicable
// Este método puede extenderse con más métricas fundamentales
func (s *Stock) GetPriceToBookRatio() float64 {
	// Placeholder - necesitaría más data de la API
	return 0.0
}

// Update actualiza los datos del stock con nueva información del mercado
// Mantiene la integridad de timestamps
func (s *Stock) Update(price, open, high, low, close float64, volume int64) {
	s.Price = price
	s.OpenPrice = open
	s.HighPrice = high
	s.LowPrice = low
	s.ClosePrice = close
	s.Volume = volume
	s.UpdatedAt = time.Now()
	s.CalculateChangePercent()
}

// Validate verifica que el stock tenga datos válidos
func (s *Stock) Validate() error {
	if s.Symbol == "" {
		return ErrInvalidSymbol
	}
	if s.Price < 0 {
		return ErrInvalidPrice
	}
	if s.Volume < 0 {
		return ErrInvalidVolume
	}
	return nil
}
