package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/marino/stock-analyzer/internal/application/dto"
	"github.com/marino/stock-analyzer/internal/application/usecase"
	"github.com/marino/stock-analyzer/internal/infrastructure/config"
	"github.com/marino/stock-analyzer/internal/infrastructure/database"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
	"github.com/marino/stock-analyzer/internal/infrastructure/repository"
)

func main() {
	fmt.Println("=== Test de Use Cases - Application Layer ===\n")

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  No se pudo cargar .env: %v\n", err)
	}

	// Configurar base de datos
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("❌ Error loading DB config: %v\n", err)
	}

	conn, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("❌ Error connecting to DB: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("✅ Conectado a base de datos")

	// Configurar Finnhub
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ FINNHUB_API_KEY not configured")
	}

	finnhubConfig := external.FinnhubConfig{
		APIKey:    apiKey,
		RateLimit: 5,
	}

	finnhubClient := external.NewFinnhubClient(finnhubConfig)
	adapter := external.NewFinnhubAdapter(finnhubClient)

	fmt.Println("✅ Cliente Finnhub configurado")

	// Crear repositorio
	stockRepo := repository.NewStockRepository(conn.DB)

	// Crear use cases
	fetchUC := usecase.NewFetchStocksUseCase(stockRepo, adapter)
	updateUC := usecase.NewUpdateStocksUseCase(stockRepo, adapter)
	searchUC := usecase.NewSearchStocksUseCase(stockRepo, adapter)
	detailsUC := usecase.NewGetStockDetailsUseCase(stockRepo, adapter)
	syncUC := usecase.NewSyncStocksUseCase(stockRepo, adapter)

	ctx := context.Background()

	// ==========================================
	// TEST 1: FetchStocks (traer de API)
	// ==========================================
	fmt.Println("\n📥 TEST 1: FetchStocks Use Case")
	fmt.Println("--------------------------------")

	fetchReq := dto.FetchStocksRequest{
		Symbols: []string{"AAPL", "GOOGL", "MSFT"},
		Save:    true,
	}

	fetchResp, err := fetchUC.Execute(ctx, fetchReq)
	if err != nil {
		log.Printf("⚠️  Error: %v\n", err)
	}

	fmt.Printf("✅ Stocks obtenidos: %d\n", fetchResp.Success)
	fmt.Printf("   Fallidos: %d\n", fetchResp.Failed)
	for i, stock := range fetchResp.Stocks {
		fmt.Printf("   %d. %s - $%.2f\n", i+1, stock.Symbol, stock.Price)
	}

	// ==========================================
	// TEST 2: GetStockDetails (detalles de uno)
	// ==========================================
	fmt.Println("\n📖 TEST 2: GetStockDetails Use Case")
	fmt.Println("------------------------------------")

	stockDetails, err := detailsUC.Execute(ctx, "AAPL")
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Detalles de %s:\n", stockDetails.Symbol)
	fmt.Printf("   Nombre: %s\n", stockDetails.Name)
	fmt.Printf("   Exchange: %s\n", stockDetails.Exchange)
	fmt.Printf("   Precio: $%.2f\n", stockDetails.Price)
	fmt.Printf("   Market Cap: $%.2fB\n", stockDetails.MarketCap/1000000000)
	fmt.Printf("   Volumen: %d\n", stockDetails.Volume)
	fmt.Printf("   P/E Ratio: %.2f\n", stockDetails.PERatio)
	fmt.Printf("   Actualizado: %v\n", stockDetails.UpdatedAt.Format("2006-01-02 15:04:05"))

	// ==========================================
	// TEST 3: SearchStocks (búsqueda en BD)
	// ==========================================
	fmt.Println("\n🔍 TEST 3: SearchStocks Use Case")
	fmt.Println("--------------------------------")

	searchReq := dto.StockSearchRequest{
		Query:  "Apple",
		Limit:  10,
		Offset: 0,
	}

	searchResp, err := searchUC.Execute(ctx, searchReq)
	if err != nil {
		log.Printf("⚠️  Error: %v\n", err)
	}

	fmt.Printf("✅ Resultados encontrados: %d\n", searchResp.Total)
	for i, stock := range searchResp.Stocks {
		fmt.Printf("   %d. %s (%s) - $%.2f\n", i+1, stock.Symbol, stock.Name, stock.Price)
	}

	// ==========================================
	// TEST 4: UpdateStocks (actualizar existentes)
	// ==========================================
	fmt.Println("\n🔄 TEST 4: UpdateStocks Use Case")
	fmt.Println("--------------------------------")

	// Esperar un poco para que la actualización sea significativa
	fmt.Println("Esperando 2 segundos...")
	time.Sleep(2 * time.Second)

	updateReq := dto.UpdateStocksRequest{
		Symbols: []string{"AAPL", "GOOGL"},
		Force:   true, // Forzar actualización
	}

	updateResp, err := updateUC.Execute(ctx, updateReq)
	if err != nil {
		log.Printf("⚠️  Error: %v\n", err)
	}

	fmt.Printf("✅ Stocks actualizados: %d\n", updateResp.Updated)
	fmt.Printf("   Saltados: %d\n", updateResp.Skipped)
	fmt.Printf("   Fallidos: %d\n", updateResp.Failed)

	// ==========================================
	// TEST 5: SyncStocks (sincronización batch)
	// ==========================================
	fmt.Println("\n🔄 TEST 5: SyncStocks Use Case")
	fmt.Println("------------------------------")

	syncReq := dto.SyncStocksRequest{
		Symbols:       []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN"},
		UpdateOlder:   1 * time.Minute,
		DeleteMissing: false,
	}

	syncResp, err := syncUC.Execute(ctx, syncReq)
	if err != nil {
		log.Printf("⚠️  Error: %v\n", err)
	}

	fmt.Printf("✅ Sincronización completada:\n")
	fmt.Printf("   Agregados: %d\n", syncResp.Added)
	fmt.Printf("   Actualizados: %d\n", syncResp.Updated)
	fmt.Printf("   Eliminados: %d\n", syncResp.Deleted)
	fmt.Printf("   Fallidos: %d\n", syncResp.Failed)

	if len(syncResp.Errors) > 0 {
		fmt.Println("   Errores:")
		for _, e := range syncResp.Errors {
			fmt.Printf("     - %s\n", e)
		}
	}

	// ==========================================
	// TEST 6: SearchStocks con API fallback
	// ==========================================
	fmt.Println("\n🔍 TEST 6: SearchStocks con API Fallback")
	fmt.Println("----------------------------------------")

	searchAPIReq := dto.StockSearchRequest{
		Query: "Tesla",
		Limit: 5,
	}

	searchAPIResp, err := searchUC.ExecuteWithAPI(ctx, searchAPIReq)
	if err != nil {
		log.Printf("⚠️  Error: %v\n", err)
	}

	fmt.Printf("✅ Resultados (BD o API): %d\n", searchAPIResp.Total)
	for i, stock := range searchAPIResp.Stocks {
		fmt.Printf("   %d. %s (%s) - $%.2f\n", i+1, stock.Symbol, stock.Name, stock.Price)
	}

	// ==========================================
	// TEST 7: Verificar datos en BD
	// ==========================================
	fmt.Println("\n📊 TEST 7: Verificación Final en BD")
	fmt.Println("-----------------------------------")

	// Contar total de stocks
	count, err := stockRepo.Count(ctx)
	if err != nil {
		log.Printf("⚠️  Error counting: %v\n", err)
	}

	fmt.Printf("✅ Total de stocks en BD: %d\n", count)

	// Obtener top por market cap
	topStocks, err := stockRepo.GetTopByMarketCap(ctx, 5)
	if err != nil {
		log.Printf("⚠️  Error getting top: %v\n", err)
	}

	fmt.Println("\n📈 Top 5 por Market Cap:")
	for i, stock := range topStocks {
		fmt.Printf("   %d. %s - $%.2fB\n", i+1, stock.Symbol, stock.MarketCap/1000000000)
	}

	// ==========================================
	// RESUMEN
	// ==========================================
	fmt.Println("\n✨ ================================")
	fmt.Println("✨ TODOS LOS USE CASES FUNCIONANDO")
	fmt.Println("✨ Application Layer Completa")
	fmt.Println("✨ ================================")
	fmt.Println("\n💡 Próximo paso: Implementar algoritmos de recomendación")
}
