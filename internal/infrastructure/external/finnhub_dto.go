package external



// FinnhubQuoteResponse respuesta del endpoint /quote
// Contiene precio actual y cambios
type FinnhubQuoteResponse struct {
	C  float64 `json:"c"`  // Current price
	D  float64 `json:"d"`  // Change
	DP float64 `json:"dp"` // Percent change
	H  float64 `json:"h"`  // High price of the day
	L  float64 `json:"l"`  // Low price of the day
	O  float64 `json:"o"`  // Open price of the day
	PC float64 `json:"pc"` // Previous close price
	T  int64   `json:"t"`  // Timestamp
}

// FinnhubProfileResponse respuesta del endpoint /stock/profile2
// Contiene información de la compañía
type FinnhubProfileResponse struct {
	Country       string  `json:"country"`
	Currency      string  `json:"currency"`
	Exchange      string  `json:"exchange"`
	IPO           string  `json:"ipo"`
	MarketCap     float64 `json:"marketCapitalization"`
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	ShareOutstanding float64 `json:"shareOutstanding"`
	Ticker        string  `json:"ticker"`
	WebURL        string  `json:"weburl"`
	Logo          string  `json:"logo"`
	Industry      string  `json:"finnhubIndustry"`
}

// FinnhubMetricsResponse respuesta del endpoint /stock/metric
// Contiene métricas fundamentales
type FinnhubMetricsResponse struct {
	Metric MetricData `json:"metric"`
	Series interface{} `json:"series"` // No lo usamos por ahora
}

// MetricData métricas fundamentales del stock
type MetricData struct {
	// PE Ratio
	PEBasicExclExtraTTM float64 `json:"peBasicExclExtraTTM"`
	PENormalizedAnnual  float64 `json:"peNormalizedAnnual"`
	
	// Dividend
	DividendYieldIndicatedAnnual float64 `json:"dividendYieldIndicatedAnnual"`
	
	// 52 Week High/Low
	Week52High float64 `json:"52WeekHigh"`
	Week52Low  float64 `json:"52WeekLow"`
	
	// Price
	Week52HighDate string `json:"52WeekHighDate"`
	Week52LowDate  string `json:"52WeekLowDate"`
	
	// Beta
	Beta float64 `json:"beta"`
}

// FinnhubSymbolLookupResponse respuesta del endpoint /search
// Para buscar stocks por nombre o símbolo
type FinnhubSymbolLookupResponse struct {
	Count  int                    `json:"count"`
	Result []FinnhubSymbolResult `json:"result"`
}

// FinnhubSymbolResult resultado individual de búsqueda
type FinnhubSymbolResult struct {
	Description string `json:"description"`
	DisplaySymbol string `json:"displaySymbol"`
	Symbol      string `json:"symbol"`
	Type        string `json:"type"`
}

// FinnhubErrorResponse respuesta de error de la API
type FinnhubErrorResponse struct {
	Error string `json:"error"`
}
