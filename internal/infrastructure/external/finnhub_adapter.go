package external

import (
	"context"
	"fmt"
	"time"

	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// FinnhubAdapter adapta las respuestas de Finnhub a entidades del dominio
// Implementa el Adapter Pattern
type FinnhubAdapter struct {
	client *FinnhubClient
}

// NewFinnhubAdapter crea una nueva instancia del adapter
func NewFinnhubAdapter(client *FinnhubClient) *FinnhubAdapter {
	return &FinnhubAdapter{
		client: client,
	}
}

// GetStock obtiene un stock completo desde Finnhub y lo convierte a entidad del dominio
func (a *FinnhubAdapter) GetStock(ctx context.Context, symbol string) (*entity.Stock, error) {
	// Obtener todos los datos de Finnhub
	quote, profile, metrics, err := a.client.GetCompleteStockData(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("error getting stock data: %w", err)
	}

	// Crear entidad del dominio
	stock, err := entity.NewStock(symbol, profile.Name, profile.Exchange)
	if err != nil {
		return nil, fmt.Errorf("error creating stock entity: %w", err)
	}

	// Mapear datos de Quote
	stock.Currency = profile.Currency
	stock.Price = quote.C
	stock.OpenPrice = quote.O
	stock.HighPrice = quote.H
	stock.LowPrice = quote.L
	stock.ClosePrice = quote.PC
	stock.Change = quote.D
	stock.ChangePercent = quote.DP

	// Mapear datos de Profile
	stock.MarketCap = profile.MarketCap * 1000000 // Finnhub retorna en millones

	// Mapear datos de Metrics (si están disponibles)
	if metrics != nil {
		stock.PERatio = metrics.Metric.PEBasicExclExtraTTM
		if stock.PERatio == 0 {
			stock.PERatio = metrics.Metric.PENormalizedAnnual
		}
		stock.DividendYield = metrics.Metric.DividendYieldIndicatedAnnual
		stock.Week52High = metrics.Metric.Week52High
		stock.Week52Low = metrics.Metric.Week52Low
	}

	// Calcular cambio porcentual (por si acaso)
	stock.CalculateChangePercent()

	// Timestamps
	stock.UpdatedAt = time.Now()

	return stock, nil
}

// GetMultipleStocks obtiene múltiples stocks
func (a *FinnhubAdapter) GetMultipleStocks(ctx context.Context, symbols []string) ([]*entity.Stock, error) {
	stocks := make([]*entity.Stock, 0, len(symbols))
	errors := make([]error, 0)

	for _, symbol := range symbols {
		stock, err := a.GetStock(ctx, symbol)
		if err != nil {
			errors = append(errors, fmt.Errorf("error getting %s: %w", symbol, err))
			continue
		}
		stocks = append(stocks, stock)
	}

	// Si hubo errores pero obtuvimos algunos stocks, retornar ambos
	if len(errors) > 0 && len(stocks) > 0 {
		return stocks, fmt.Errorf("partial success: %d errors occurred", len(errors))
	}

	// Si todos fallaron, retornar error
	if len(stocks) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("all requests failed: %v", errors)
	}

	return stocks, nil
}

// SearchStocks busca stocks y los convierte a entidades del dominio
func (a *FinnhubAdapter) SearchStocks(ctx context.Context, query string, limit int) ([]*entity.Stock, error) {
	// Buscar símbolos
	result, err := a.client.SearchSymbol(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error searching stocks: %w", err)
	}

	if result.Count == 0 {
		return []*entity.Stock{}, nil
	}

	// Limitar resultados
	count := result.Count
	if limit > 0 && limit < count {
		count = limit
	}

	// Obtener datos completos para cada símbolo encontrado
	stocks := make([]*entity.Stock, 0, count)
	for i := 0; i < count && i < len(result.Result); i++ {
		symbolResult := result.Result[i]
		
		// Solo procesar stocks (ignorar otros tipos como crypto, forex, etc)
		if symbolResult.Type != "Common Stock" && symbolResult.Type != "EQS" {
			continue
		}

		stock, err := a.GetStock(ctx, symbolResult.Symbol)
		if err != nil {
			// Log error pero continuar con los demás
			continue
		}

		stocks = append(stocks, stock)
	}

	return stocks, nil
}
