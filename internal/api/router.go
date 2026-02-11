package api

import (
	"github.com/gin-gonic/gin"
	"github.com/marino/stock-analyzer/internal/api/handler"
	"github.com/marino/stock-analyzer/internal/api/middleware"
)

// Router configura todas las rutas de la API
type Router struct {
	stockHandler          *handler.StockHandler
	recommendationHandler *handler.RecommendationHandler
}

// NewRouter crea una nueva instancia del router
func NewRouter(
	stockHandler *handler.StockHandler,
	recommendationHandler *handler.RecommendationHandler,
) *Router {
	return &Router{
		stockHandler:          stockHandler,
		recommendationHandler: recommendationHandler,
	}
}

// Setup configura todas las rutas
func (r *Router) Setup() *gin.Engine {
	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middlewares globales
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", r.HealthCheck)
	router.GET("/", r.Welcome)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Stocks endpoints
		stocks := v1.Group("/stocks")
		{
			stocks.GET("/:symbol", r.stockHandler.GetStock)           // Obtener un stock
			stocks.GET("", r.stockHandler.SearchStocks)               // Buscar stocks
			stocks.POST("/fetch", r.stockHandler.FetchStocks)         // Traer desde API
			stocks.PUT("/update", r.stockHandler.UpdateStocks)        // Actualizar existentes
			stocks.POST("/sync", r.stockHandler.SyncStocks)           // Sincronizar batch
		}

		// Recommendations endpoints
		recommendations := v1.Group("/recommendations")
		{
			recommendations.POST("/generate", r.recommendationHandler.GenerateRecommendations)     // Generar recomendaciones
			recommendations.GET("/top", r.recommendationHandler.GetTopRecommendations)             // Top recomendaciones
			recommendations.GET("/stock/:symbol", r.recommendationHandler.GetRecommendationsByStock) // Por stock
			recommendations.GET("/type/:type", r.recommendationHandler.GetRecommendationsByType)   // Por tipo
			recommendations.GET("/valid", r.recommendationHandler.GetValidRecommendations)         // Válidas
			recommendations.GET("/stats", r.recommendationHandler.GetRecommendationStats)          // Estadísticas
		}
	}

	return router
}

// HealthCheck endpoint de salud
func (r *Router) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "stock-analyzer-api",
		"version": "1.0.0",
	})
}

// Welcome mensaje de bienvenida
func (r *Router) Welcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to Stock Analyzer API",
		"version": "1.0.0",
		"endpoints": gin.H{
			"health":           "GET /health",
			"stocks":           "GET /api/v1/stocks",
			"stock_detail":     "GET /api/v1/stocks/:symbol",
			"fetch_stocks":     "POST /api/v1/stocks/fetch",
			"update_stocks":    "PUT /api/v1/stocks/update",
			"sync_stocks":      "POST /api/v1/stocks/sync",
			"generate_recs":    "POST /api/v1/recommendations/generate",
			"top_recs":         "GET /api/v1/recommendations/top",
			"stock_recs":       "GET /api/v1/recommendations/stock/:symbol",
			"recs_by_type":     "GET /api/v1/recommendations/type/:type",
			"valid_recs":       "GET /api/v1/recommendations/valid",
			"recs_stats":       "GET /api/v1/recommendations/stats",
		},
	})
}
