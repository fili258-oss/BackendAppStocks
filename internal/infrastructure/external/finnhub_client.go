package external

import (
	"context"
	"fmt"
	"time"
)

// FinnhubClient cliente para interactuar con Finnhub API
type FinnhubClient struct {
	httpClient *HTTPClient
	apiKey     string
	rateLimiter *RateLimiter
}

// FinnhubConfig configuración para Finnhub
type FinnhubConfig struct {
	APIKey      string
	BaseURL     string
	Timeout     time.Duration
	RateLimit   int // Requests por segundo
}

// NewFinnhubClient crea una nueva instancia del cliente Finnhub
func NewFinnhubClient(config FinnhubConfig) *FinnhubClient {
	if config.BaseURL == "" {
		config.BaseURL = "https://finnhub.io/api/v1"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RateLimit == 0 {
		config.RateLimit = 5 // Free tier: 60 calls/min = ~1/sec, usamos 5 por seguridad
	}

	httpClient := NewHTTPClient(HTTPClientConfig{
		BaseURL: config.BaseURL,
		Timeout: config.Timeout,
	})

	return &FinnhubClient{
		httpClient: httpClient,
		apiKey:     config.APIKey,
		rateLimiter: NewRateLimiter(config.RateLimit, time.Second),
	}
}

// GetQuote obtiene el precio actual y datos del día para un símbolo
func (c *FinnhubClient) GetQuote(ctx context.Context, symbol string) (*FinnhubQuoteResponse, error) {
	// Aplicar rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	params := map[string]string{
		"symbol": symbol,
		"token":  c.apiKey,
	}

	var quote FinnhubQuoteResponse
	if err := c.httpClient.GetJSON(ctx, "/quote", params, &quote); err != nil {
		return nil, fmt.Errorf("error getting quote for %s: %w", symbol, err)
	}

	// Verificar si la respuesta es válida
	if quote.C == 0 && quote.O == 0 {
		return nil, fmt.Errorf("invalid or non-existent symbol: %s", symbol)
	}

	return &quote, nil
}

// GetProfile obtiene información de la compañía
func (c *FinnhubClient) GetProfile(ctx context.Context, symbol string) (*FinnhubProfileResponse, error) {
	// Aplicar rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	params := map[string]string{
		"symbol": symbol,
		"token":  c.apiKey,
	}

	var profile FinnhubProfileResponse
	if err := c.httpClient.GetJSON(ctx, "/stock/profile2", params, &profile); err != nil {
		return nil, fmt.Errorf("error getting profile for %s: %w", symbol, err)
	}

	// Verificar si la respuesta es válida
	if profile.Ticker == "" {
		return nil, fmt.Errorf("profile not found for symbol: %s", symbol)
	}

	return &profile, nil
}

// GetMetrics obtiene métricas fundamentales del stock
func (c *FinnhubClient) GetMetrics(ctx context.Context, symbol string) (*FinnhubMetricsResponse, error) {
	// Aplicar rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	params := map[string]string{
		"symbol": symbol,
		"metric": "all",
		"token":  c.apiKey,
	}

	var metrics FinnhubMetricsResponse
	if err := c.httpClient.GetJSON(ctx, "/stock/metric", params, &metrics); err != nil {
		return nil, fmt.Errorf("error getting metrics for %s: %w", symbol, err)
	}

	return &metrics, nil
}

// SearchSymbol busca stocks por nombre o símbolo
func (c *FinnhubClient) SearchSymbol(ctx context.Context, query string) (*FinnhubSymbolLookupResponse, error) {
	// Aplicar rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	params := map[string]string{
		"q":     query,
		"token": c.apiKey,
	}

	var result FinnhubSymbolLookupResponse
	if err := c.httpClient.GetJSON(ctx, "/search", params, &result); err != nil {
		return nil, fmt.Errorf("error searching symbol %s: %w", query, err)
	}

	return &result, nil
}

// GetCompleteStockData obtiene todos los datos necesarios para un stock
// Combina Quote + Profile + Metrics en una sola operación
func (c *FinnhubClient) GetCompleteStockData(ctx context.Context, symbol string) (*FinnhubQuoteResponse, *FinnhubProfileResponse, *FinnhubMetricsResponse, error) {
	// Obtener Quote (precio actual)
	quote, err := c.GetQuote(ctx, symbol)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting quote: %w", err)
	}

	// Obtener Profile (info de compañía)
	profile, err := c.GetProfile(ctx, symbol)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting profile: %w", err)
	}

	// Obtener Metrics (métricas fundamentales) - opcional
	metrics, err := c.GetMetrics(ctx, symbol)
	if err != nil {
		// Las métricas son opcionales, si fallan continuamos sin ellas
		return quote, profile, nil, nil
	}

	return quote, profile, metrics, nil
}
