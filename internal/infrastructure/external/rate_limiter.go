package external

import (
	"context"
	"sync"
	"time"
)

// RateLimiter controla la tasa de peticiones
// Implementa un token bucket algorithm
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter crea un nuevo rate limiter
// requestsPerPeriod: número de requests permitidas
// period: período de tiempo (ej: 1 segundo, 1 minuto)
func NewRateLimiter(requestsPerPeriod int, period time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     requestsPerPeriod,
		maxTokens:  requestsPerPeriod,
		refillRate: period / time.Duration(requestsPerPeriod),
		lastRefill: time.Now(),
	}
}

// Wait espera hasta que haya un token disponible
func (rl *RateLimiter) Wait(ctx context.Context) error {
	for {
		// Verificar si el contexto fue cancelado
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Intentar obtener un token
		if rl.tryAcquire() {
			return nil
		}

		// Esperar antes de reintentar
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(rl.refillRate):
			// Continuar intentando
		}
	}
}

// tryAcquire intenta obtener un token
func (rl *RateLimiter) tryAcquire() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Rellenar tokens basado en el tiempo transcurrido
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}

	// Intentar consumir un token
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// GetAvailableTokens retorna el número de tokens disponibles
func (rl *RateLimiter) GetAvailableTokens() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.tokens
}
