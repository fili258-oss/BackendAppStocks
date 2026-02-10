package main

import (
	"context"
	"fmt"
	"log"
	

	"github.com/joho/godotenv"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/infrastructure/config"
	"github.com/marino/stock-analyzer/internal/infrastructure/database"
	"github.com/marino/stock-analyzer/internal/infrastructure/repository"
)

func main() {
	fmt.Println("=== Test de Repositorios ===\n")

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  No se pudo cargar .env: %v\n", err)
	}

	// Cargar configuración de base de datos
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("❌ Error al cargar configuración: %v\n", err)
	}

	// Conectar a la base de datos
	conn, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("✅ Conexión a base de datos exitosa\n")

	// Crear repositorios
	stockRepo := repository.NewStockRepository(conn.DB)
	recRepo := repository.NewRecommendationRepository(conn.DB)

	ctx := context.Background()

	// ==========================================
	// TEST 1: StockRepository - CRUD
	// ==========================================
	fmt.Println("📊 TEST 1: StockRepository - CRUD")
	fmt.Println("----------------------------------")

	// Crear stock de prueba
	stock, err := entity.NewStock("AAPL", "Apple Inc.", "NASDAQ")
	if err != nil {
		log.Fatalf("Error al crear entidad stock: %v\n", err)
	}

	stock.Currency = "USD"
	stock.Price = 175.50
	stock.OpenPrice = 174.00
	stock.HighPrice = 176.00
	stock.LowPrice = 173.50
	stock.ClosePrice = 172.00
	stock.Volume = 50000000
	stock.MarketCap = 2750000000000 // 2.75 trillion
	stock.PERatio = 28.5
	stock.DividendYield = 0.52
	stock.Week52High = 180.00
	stock.Week52Low = 140.00
	stock.CalculateChangePercent()

	// Crear en BD
	if err := stockRepo.Create(ctx, stock); err != nil {
		log.Fatalf("❌ Error al crear stock: %v\n", err)
	}
	fmt.Printf("✅ Stock creado: %s (ID: %s)\n", stock.Symbol, stock.ID)

	// Leer por ID
	stockRead, err := stockRepo.FindByID(ctx, stock.ID)
	if err != nil {
		log.Fatalf("❌ Error al leer stock: %v\n", err)
	}
	fmt.Printf("✅ Stock leído: %s - $%.2f\n", stockRead.Name, stockRead.Price)

	// Actualizar
	stock.Price = 178.00
	stock.CalculateChangePercent()
	if err := stockRepo.Update(ctx, stock); err != nil {
		log.Fatalf("❌ Error al actualizar stock: %v\n", err)
	}
	fmt.Printf("✅ Stock actualizado: nuevo precio $%.2f\n", stock.Price)

	// Buscar por símbolo
	stockBySymbol, err := stockRepo.FindBySymbol(ctx, "AAPL")
	if err != nil {
		log.Fatalf("❌ Error al buscar por símbolo: %v\n", err)
	}
	fmt.Printf("✅ Stock encontrado por símbolo: %s\n", stockBySymbol.Symbol)

	// ==========================================
	// TEST 2: StockRepository - Bulk Create
	// ==========================================
	fmt.Println("\n📊 TEST 2: StockRepository - Bulk Create")
	fmt.Println("----------------------------------------")

	stocks := []*entity.Stock{}
	
	// Crear varios stocks de ejemplo
	symbols := []struct {
		symbol   string
		name     string
		exchange string
		price    float64
		marketCap float64
	}{
		{"GOOGL", "Alphabet Inc.", "NASDAQ", 140.50, 1800000000000},
		{"MSFT", "Microsoft Corporation", "NASDAQ", 380.00, 2850000000000},
		{"TSLA", "Tesla Inc.", "NASDAQ", 245.00, 750000000000},
		{"AMZN", "Amazon.com Inc.", "NASDAQ", 155.00, 1600000000000},
	}

	for _, s := range symbols {
		st, _ := entity.NewStock(s.symbol, s.name, s.exchange)
		st.Currency = "USD"
		st.Price = s.price
		st.MarketCap = s.marketCap
		st.Volume = 10000000
		stocks = append(stocks, st)
	}

	if err := stockRepo.BulkCreate(ctx, stocks); err != nil {
		log.Fatalf("❌ Error en bulk create: %v\n", err)
	}
	fmt.Printf("✅ %d stocks creados en lote\n", len(stocks))

	// Contar total
	count, err := stockRepo.Count(ctx)
	if err != nil {
		log.Fatalf("❌ Error al contar: %v\n", err)
	}
	fmt.Printf("✅ Total de stocks en BD: %d\n", count)

	// ==========================================
	// TEST 3: StockRepository - Búsqueda
	// ==========================================
	fmt.Println("\n📊 TEST 3: StockRepository - Búsqueda")
	fmt.Println("-------------------------------------")

	// Top por market cap
	topByMarketCap, err := stockRepo.GetTopByMarketCap(ctx, 3)
	if err != nil {
		log.Fatalf("❌ Error al obtener top por market cap: %v\n", err)
	}
	fmt.Println("✅ Top 3 por capitalización:")
	for i, s := range topByMarketCap {
		fmt.Printf("   %d. %s - $%.2fB\n", i+1, s.Symbol, s.MarketCap/1000000000)
	}

	// ==========================================
	// TEST 4: RecommendationRepository - CRUD
	// ==========================================
	fmt.Println("\n🎯 TEST 4: RecommendationRepository - CRUD")
	fmt.Println("------------------------------------------")

	// Crear recomendación
	recommendation, err := entity.NewRecommendation(
		stock.ID,
		stock.Symbol,
		"BALANCED",
		85.5,
	)
	if err != nil {
		log.Fatalf("❌ Error al crear entidad recommendation: %v\n", err)
	}

	recommendation.SetReason("Strong momentum and solid fundamentals")
	recommendation.AddMetric("volume_increase", 15.5)
	recommendation.AddMetric("price_momentum", 12.3)

	// Crear en BD
	if err := recRepo.Create(ctx, recommendation); err != nil {
		log.Fatalf("❌ Error al crear recommendation: %v\n", err)
	}
	fmt.Printf("✅ Recomendación creada: %s para %s (Score: %.1f)\n", 
		recommendation.Type, recommendation.StockSymbol, recommendation.Score)

	// Leer por ID
	recRead, err := recRepo.FindByID(ctx, recommendation.ID)
	if err != nil {
		log.Fatalf("❌ Error al leer recommendation: %v\n", err)
	}
	fmt.Printf("✅ Recomendación leída: %s - Confianza: %.2f\n", 
		recRead.Type, recRead.Confidence)

	// ==========================================
	// TEST 5: RecommendationRepository - Búsqueda
	// ==========================================
	fmt.Println("\n🎯 TEST 5: RecommendationRepository - Búsqueda")
	fmt.Println("----------------------------------------------")

	// Buscar por stock
	recsByStock, err := recRepo.FindByStockID(ctx, stock.ID)
	if err != nil {
		log.Fatalf("❌ Error al buscar por stock: %v\n", err)
	}
	fmt.Printf("✅ Recomendaciones para %s: %d\n", stock.Symbol, len(recsByStock))

	// Top recommendations
	topRecs, err := recRepo.FindTopRecommendations(ctx, 5)
	if err != nil {
		log.Fatalf("❌ Error al obtener top recommendations: %v\n", err)
	}
	fmt.Printf("✅ Top recomendaciones: %d\n", len(topRecs))
	for i, r := range topRecs {
		fmt.Printf("   %d. %s (%s) - Score: %.1f\n", 
			i+1, r.StockSymbol, r.Type, r.Score)
	}

	// ==========================================
	// LIMPIEZA
	// ==========================================
	fmt.Println("\n🧹 LIMPIEZA")
	fmt.Println("-----------")

	// Eliminar recomendación
	if err := recRepo.Delete(ctx, recommendation.ID); err != nil {
		log.Printf("⚠️  Error al eliminar recommendation: %v\n", err)
	} else {
		fmt.Println("✅ Recomendación eliminada")
	}

	// Eliminar stocks (esto eliminará recomendaciones en cascada)
	allStocks := append([]*entity.Stock{stock}, stocks...)
	for _, s := range allStocks {
		if err := stockRepo.Delete(ctx, s.ID); err != nil {
			log.Printf("⚠️  Error al eliminar stock %s: %v\n", s.Symbol, err)
		}
	}
	fmt.Printf("✅ %d stocks eliminados\n", len(allStocks))

	// ==========================================
	// RESUMEN
	// ==========================================
	fmt.Println("\n✨ ================================")
	fmt.Println("✨ TODOS LOS TESTS PASARON")
	fmt.Println("✨ Repositorios funcionando correctamente")
	fmt.Println("✨ ================================")
}
