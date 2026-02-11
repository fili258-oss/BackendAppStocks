package dto

import "time"

// StockDTO representa un stock en la capa de aplicación
type StockDTO struct {
	ID             string    `json:"id"`
	Symbol         string    `json:"symbol"`
	Name           string    `json:"name"`
	Exchange       string    `json:"exchange"`
	Currency       string    `json:"currency"`
	Price          float64   `json:"price"`
	OpenPrice      float64   `json:"open_price"`
	HighPrice      float64   `json:"high_price"`
	LowPrice       float64   `json:"low_price"`
	ClosePrice     float64   `json:"close_price"`
	Change         float64   `json:"change"`
	ChangePercent  float64   `json:"change_percent"`
	Volume         int64     `json:"volume"`
	MarketCap      float64   `json:"market_cap"`
	PERatio        float64   `json:"pe_ratio"`
	DividendYield  float64   `json:"dividend_yield"`
	Week52High     float64   `json:"week_52_high"`
	Week52Low      float64   `json:"week_52_low"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// StockSearchRequest request para búsqueda de stocks
type StockSearchRequest struct {
	Query         string  `json:"query"`
	Exchange      string  `json:"exchange,omitempty"`
	MinPrice      float64 `json:"min_price,omitempty"`
	MaxPrice      float64 `json:"max_price,omitempty"`
	MinVolume     int64   `json:"min_volume,omitempty"`
	MinMarketCap  float64 `json:"min_market_cap,omitempty"`
	SortBy        string  `json:"sort_by,omitempty"`
	SortDirection string  `json:"sort_direction,omitempty"`
	Limit         int     `json:"limit,omitempty"`
	Offset        int     `json:"offset,omitempty"`
}

// StockListResponse respuesta con lista de stocks
type StockListResponse struct {
	Stocks []StockDTO `json:"stocks"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// FetchStocksRequest request para traer stocks de la API
type FetchStocksRequest struct {
	Symbols []string `json:"symbols"`
	Save    bool     `json:"save"` // Si es true, guarda en BD
}

// FetchStocksResponse respuesta de fetch stocks
type FetchStocksResponse struct {
	Stocks    []StockDTO `json:"stocks"`
	Success   int        `json:"success"`
	Failed    int        `json:"failed"`
	Errors    []string   `json:"errors,omitempty"`
}

// UpdateStocksRequest request para actualizar stocks
type UpdateStocksRequest struct {
	Symbols []string `json:"symbols,omitempty"` // Si está vacío, actualiza todos
	Force   bool     `json:"force"`             // Forzar actualización aunque sean recientes
}

// UpdateStocksResponse respuesta de actualización
type UpdateStocksResponse struct {
	Updated int      `json:"updated"`
	Skipped int      `json:"skipped"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

// SyncStocksRequest request para sincronización batch
type SyncStocksRequest struct {
	Symbols      []string      `json:"symbols"`
	UpdateOlder  time.Duration `json:"update_older,omitempty"` // Actualizar solo si es más viejo que esto
	DeleteMissing bool         `json:"delete_missing"`         // Eliminar stocks que ya no existen
}

// SyncStocksResponse respuesta de sincronización
type SyncStocksResponse struct {
	Added   int      `json:"added"`
	Updated int      `json:"updated"`
	Deleted int      `json:"deleted"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}
