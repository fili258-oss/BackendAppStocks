package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/marino/stock-analyzer/internal/api/common"
	"github.com/marino/stock-analyzer/internal/application/dto"
	"github.com/marino/stock-analyzer/internal/application/usecase"
	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// StockHandler maneja las peticiones relacionadas con stocks
type StockHandler struct {
	fetchUC   *usecase.FetchStocksUseCase
	updateUC  *usecase.UpdateStocksUseCase
	searchUC  *usecase.SearchStocksUseCase
	detailsUC *usecase.GetStockDetailsUseCase
	syncUC    *usecase.SyncStocksUseCase
}

// NewStockHandler crea una nueva instancia
func NewStockHandler(
	fetchUC *usecase.FetchStocksUseCase,
	updateUC *usecase.UpdateStocksUseCase,
	searchUC *usecase.SearchStocksUseCase,
	detailsUC *usecase.GetStockDetailsUseCase,
	syncUC *usecase.SyncStocksUseCase,
) *StockHandler {
	return &StockHandler{
		fetchUC:   fetchUC,
		updateUC:  updateUC,
		searchUC:  searchUC,
		detailsUC: detailsUC,
		syncUC:    syncUC,
	}
}

// GetStock obtiene un stock por símbolo
// GET /api/stocks/:symbol
func (h *StockHandler) GetStock(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		common.ValidationErrorResponse(c, "Symbol is required")
		return
	}

	stock, err := h.detailsUC.Execute(c.Request.Context(), symbol)
	if err != nil {
		if err == entity.ErrStockNotFound {
			common.NotFoundResponse(c, "Stock")
			return
		}
		common.InternalErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, stock)
}

// SearchStocks busca stocks con filtros
// GET /api/stocks?query=...&exchange=...&limit=...
func (h *StockHandler) SearchStocks(c *gin.Context) {
	// Parsear query params
	request := dto.StockSearchRequest{
		Query:         c.Query("query"),
		Exchange:      c.Query("exchange"),
		SortBy:        c.DefaultQuery("sort_by", "symbol"),
		SortDirection: c.DefaultQuery("sort_direction", "asc"),
		Limit:         parseIntQuery(c, "limit", 20),
		Offset:        parseIntQuery(c, "offset", 0),
	}

	// Parsear filtros numéricos
	if minPrice := c.Query("min_price"); minPrice != "" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			request.MinPrice = val
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			request.MaxPrice = val
		}
	}
	if minVolume := c.Query("min_volume"); minVolume != "" {
		if val, err := strconv.ParseInt(minVolume, 10, 64); err == nil {
			request.MinVolume = val
		}
	}
	if minMarketCap := c.Query("min_market_cap"); minMarketCap != "" {
		if val, err := strconv.ParseFloat(minMarketCap, 64); err == nil {
			request.MinMarketCap = val
		}
	}

	// Ejecutar búsqueda
	response, err := h.searchUC.Execute(c.Request.Context(), request)
	if err != nil {
		common.InternalErrorResponse(c, err)
		return
	}

	// Calcular metadata de paginación
	meta := &common.MetaInfo{
		Page:       (request.Offset / request.Limit) + 1,
		PageSize:   request.Limit,
		TotalItems: response.Total,
		TotalPages: (response.Total + request.Limit - 1) / request.Limit,
	}

	common.SuccessResponseWithMeta(c, http.StatusOK, response.Stocks, meta)
}

// FetchStocks trae stocks desde la API externa
// POST /api/stocks/fetch
// Body: { "symbols": ["AAPL", "GOOGL"], "save": true }
func (h *StockHandler) FetchStocks(c *gin.Context) {
	var request dto.FetchStocksRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		common.ValidationErrorResponse(c, err.Error())
		return
	}

	if len(request.Symbols) == 0 {
		common.ValidationErrorResponse(c, "Symbols array is required")
		return
	}

	response, err := h.fetchUC.Execute(c.Request.Context(), request)
	if err != nil {
		common.InternalErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, response)
}

// UpdateStocks actualiza stocks existentes
// PUT /api/stocks/update
// Body: { "symbols": ["AAPL"], "force": true }
func (h *StockHandler) UpdateStocks(c *gin.Context) {
	var request dto.UpdateStocksRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		common.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.updateUC.Execute(c.Request.Context(), request)
	if err != nil {
		common.InternalErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, response)
}

// SyncStocks sincroniza stocks en batch
// POST /api/stocks/sync
// Body: { "symbols": [...], "update_older": "5m", "delete_missing": false }
func (h *StockHandler) SyncStocks(c *gin.Context) {
	var request dto.SyncStocksRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		common.ValidationErrorResponse(c, err.Error())
		return
	}

	if len(request.Symbols) == 0 {
		common.ValidationErrorResponse(c, "Symbols array is required")
		return
	}

	response, err := h.syncUC.Execute(c.Request.Context(), request)
	if err != nil {
		common.InternalErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, response)
}

// parseIntQuery helper para parsear query params enteros
func parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	if val := c.Query(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return defaultValue
}
