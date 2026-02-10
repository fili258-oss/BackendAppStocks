package entity

import (
	"errors"
	"time"
)

// TimeRange representa un rango de tiempo
// Value Object: inmutable y se compara por valor
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// NewTimeRange crea un nuevo rango de tiempo con validación
func NewTimeRange(start, end time.Time) (*TimeRange, error) {
	if end.Before(start) {
		return nil, errors.New("end time cannot be before start time")
	}
	return &TimeRange{
		Start: start,
		End:   end,
	}, nil
}

// Contains verifica si un tiempo está dentro del rango
func (tr TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && !t.After(tr.End)
}

// Duration retorna la duración del rango
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// PriceRange representa un rango de precios
type PriceRange struct {
	Min float64
	Max float64
}

// NewPriceRange crea un nuevo rango de precios con validación
func NewPriceRange(min, max float64) (*PriceRange, error) {
	if min < 0 || max < 0 {
		return nil, errors.New("prices cannot be negative")
	}
	if max < min {
		return nil, errors.New("max price cannot be less than min price")
	}
	return &PriceRange{
		Min: min,
		Max: max,
	}, nil
}

// Contains verifica si un precio está dentro del rango
func (pr PriceRange) Contains(price float64) bool {
	return price >= pr.Min && price <= pr.Max
}

// Range retorna el tamaño del rango
func (pr PriceRange) Range() float64 {
	return pr.Max - pr.Min
}

// Midpoint retorna el punto medio del rango
func (pr PriceRange) Midpoint() float64 {
	return (pr.Min + pr.Max) / 2
}
