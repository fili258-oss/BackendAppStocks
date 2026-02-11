package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/marino/stock-analyzer/internal/api"
	"github.com/marino/stock-analyzer/internal/application/usecase"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/domain/repository"
)

// RecommendationHandler maneja las peticiones de recomendaciones
type RecommendationHandler struct {
	generateUC *usecase.GenerateRecommendationsUseCase
	recRepo    repository.RecommendationRepository
}

// NewRecommendationHandler crea una nueva instancia
func NewRecommendationHandler(
	generateUC *usecase.GenerateRecommendationsUseCase,
	recRepo repository.RecommendationRepository,
) *RecommendationHandler {
	return &RecommendationHandler{
		generateUC: generateUC,
		recRepo:    recRepo,
	}
}

// GenerateRecommendations genera recomendaciones
// POST /api/recommendations/generate
// Body: { "symbols": ["AAPL"], "strategies": ["BALANCED"], "save_to_db": true }
func (h *RecommendationHandler) GenerateRecommendations(c *gin.Context) {
	var request usecase.GenerateRecommendationsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		api.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.generateUC.Execute(c.Request.Context(), request)
	if err != nil {
		api.InternalErrorResponse(c, err)
		return
	}

	api.SuccessResponse(c, http.StatusOK, response)
}

// GetTopRecommendations obtiene las mejores recomendaciones
// GET /api/recommendations/top?limit=10
func (h *RecommendationHandler) GetTopRecommendations(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 10)
	if limit > 100 {
		limit = 100 // Máximo 100
	}

	recommendations, err := h.recRepo.FindTopRecommendations(c.Request.Context(), limit)
	if err != nil {
		api.InternalErrorResponse(c, err)
		return
	}

	api.SuccessResponse(c, http.StatusOK, gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetRecommendationsByStock obtiene recomendaciones de un stock
// GET /api/recommendations/stock/:symbol
func (h *RecommendationHandler) GetRecommendationsByStock(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		api.ValidationErrorResponse(c, "Symbol is required")
		return
	}

	recommendations, err := h.recRepo.FindByStockSymbol(c.Request.Context(), symbol)
	if err != nil {
		api.InternalErrorResponse(c, err)
		return
	}

	api.SuccessResponse(c, http.StatusOK, gin.H{
		"symbol":          symbol,
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetRecommendationsByType obtiene recomendaciones por tipo
// GET /api/recommendations/type/:type?limit=20
func (h *RecommendationHandler) GetRecommendationsByType(c *gin.Context) {
	recType := c.Param("type")
	limit := parseIntQuery(c, "limit", 20)

	// Validar tipo
	validTypes := map[string]bool{
		"BUY":        true,
		"SELL":       true,
		"HOLD":       true,
		"STRONG_BUY": true,
	}

	if !validTypes[recType] {
		api.ValidationErrorResponse(c, "Invalid recommendation type. Valid: BUY, SELL, HOLD, STRONG_BUY")
		return
	}

	recommendations, err := h.recRepo.FindByType(c.Request.Context(), entity.RecommendationType(recType), limit)
	if err != nil {
		api.InternalErrorResponse(c, err)
		return
	}

	api.SuccessResponse(c, http.StatusOK, gin.H{
		"type":            recType,
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetValidRecommendations obtiene recomendaciones válidas (no expiradas)
// GET /api/recommendations/valid?limit=50
func (h *RecommendationHandler) GetValidRecommendations(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 50)

	recommendations, err := h.recRepo.FindValidRecommendations(c.Request.Context(), limit)
	if err != nil {
		api.InternalErrorResponse(c, err)
		return
	}

	api.SuccessResponse(c, http.StatusOK, gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetRecommendationStats obtiene estadísticas de recomendaciones
// GET /api/recommendations/stats
func (h *RecommendationHandler) GetRecommendationStats(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener conteos por tipo
	buyRecs, _ := h.recRepo.FindByType(ctx, entity.RecommendationBuy, 10000)
	strongBuyRecs, _ := h.recRepo.FindByType(ctx, entity.RecommendationStrongBuy, 10000)
	holdRecs, _ := h.recRepo.FindByType(ctx, entity.RecommendationHold, 10000)
	sellRecs, _ := h.recRepo.FindByType(ctx, entity.RecommendationSell, 10000)

	total, _ := h.recRepo.Count(ctx)

	stats := gin.H{
		"total": total,
		"by_type": gin.H{
			"STRONG_BUY": len(strongBuyRecs),
			"BUY":        len(buyRecs),
			"HOLD":       len(holdRecs),
			"SELL":       len(sellRecs),
		},
	}

	api.SuccessResponse(c, http.StatusOK, stats)
}
