package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/marino/stock-analyzer/internal/api"
	"github.com/marino/stock-analyzer/internal/api/handler"
	"github.com/marino/stock-analyzer/internal/application/usecase"
	"github.com/marino/stock-analyzer/internal/infrastructure/config"
	"github.com/marino/stock-analyzer/internal/infrastructure/database"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
	"github.com/marino/stock-analyzer/internal/infrastructure/repository"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Configuración de la base de datos
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Conectar a la base de datos
	conn, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	log.Println("✅ Connected to database")

	// Verificar salud de la BD
	if err := conn.HealthCheck(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	// Configurar Finnhub client
	finnhubConfig := external.FinnhubConfig{
		APIKey:    os.Getenv("FINNHUB_API_KEY"),
		RateLimit: 5,
	}

	finnhubClient := external.NewFinnhubClient(finnhubConfig)
	finnhubAdapter := external.NewFinnhubAdapter(finnhubClient)

	log.Println("✅ Finnhub client configured")

	// Crear repositorios
	stockRepo := repository.NewStockRepository(conn.DB)
	recRepo := repository.NewRecommendationRepository(conn.DB)

	// Crear use cases
	fetchUC := usecase.NewFetchStocksUseCase(stockRepo, finnhubAdapter)
	updateUC := usecase.NewUpdateStocksUseCase(stockRepo, finnhubAdapter)
	searchUC := usecase.NewSearchStocksUseCase(stockRepo, finnhubAdapter)
	detailsUC := usecase.NewGetStockDetailsUseCase(stockRepo, finnhubAdapter)
	syncUC := usecase.NewSyncStocksUseCase(stockRepo, finnhubAdapter)
	generateRecUC := usecase.NewGenerateRecommendationsUseCase(stockRepo, recRepo)

	// Crear handlers
	stockHandler := handler.NewStockHandler(fetchUC, updateUC, searchUC, detailsUC, syncUC)
	recHandler := handler.NewRecommendationHandler(generateRecUC, recRepo)

	// Crear router
	apiRouter := api.NewRouter(stockHandler, recHandler)
	router := apiRouter.Setup()

	// Configurar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		log.Printf("📡 API available at http://localhost:%s", port)
		log.Printf("💊 Health check at http://localhost:%s/health", port)
		log.Printf("📖 API docs at http://localhost:%s/", port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
