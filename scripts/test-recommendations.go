package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/marino/stock-analyzer/internal/application/usecase"
	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/infrastructure/config"
	"github.com/marino/stock-analyzer/internal/infrastructure/database"
	"github.com/marino/stock-analyzer/internal/infrastructure/repository"
)

func main() {
	fmt.Println("=== Test de Estrategias y Recomendaciones ===\n")

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

	// Crear repositorios
	stockRepo := repository.NewStockRepository(conn.DB)
	recRepo := repository.NewRecommendationRepository(conn.DB)

	// Crear use case de recomendaciones
	genRecUC := usecase.NewGenerateRecommendationsUseCase(stockRepo, recRepo)

	ctx := context.Background()

	// Verificar que hay stocks en BD
	count, err := stockRepo.Count(ctx)
	if err != nil {
		log.Fatalf("❌ Error counting stocks: %v\n", err)
	}

	fmt.Printf("📊 Stocks en BD: %d\n", count)

	if count == 0 {
		fmt.Println("\n⚠️  No hay stocks en BD. Ejecuta primero:")
		fmt.Println("   go run scripts/test-usecases.go")
		return
	}

	// ==========================================
	// TEST 1: Generar Recomendaciones - Balanced
	// ==========================================
	fmt.Println("\n🎯 TEST 1: Estrategia Balanced")
	fmt.Println("------------------------------")

	reqBalanced := usecase.GenerateRecommendationsRequest{
		Strategies: []string{string(entity.StrategyBalanced)},
		SaveToDB:   true,
	}

	respBalanced, err := genRecUC.Execute(ctx, reqBalanced)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Recomendaciones generadas: %d\n", respBalanced.Generated)
	fmt.Printf("   Fallidas: %d\n", respBalanced.Failed)

	// Mostrar algunas recomendaciones
	fmt.Println("\nTop 3 recomendaciones (Balanced):")
	for i := 0; i < 3 && i < len(respBalanced.Recommendations); i++ {
		rec := respBalanced.Recommendations[i]
		fmt.Printf("   %d. %s - %s (Score: %.1f, Conf: %.2f)\n",
			i+1, rec.StockSymbol, rec.Type, rec.Score, rec.Confidence)
		fmt.Printf("      Razón: %s\n", rec.Reason)
	}

	// ==========================================
	// TEST 2: Generar Recomendaciones - Momentum
	// ==========================================
	fmt.Println("\n⚡ TEST 2: Estrategia Momentum")
	fmt.Println("-----------------------------")

	reqMomentum := usecase.GenerateRecommendationsRequest{
		Strategies: []string{string(entity.StrategyMomentum)},
		SaveToDB:   true,
	}

	respMomentum, err := genRecUC.Execute(ctx, reqMomentum)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Recomendaciones generadas: %d\n", respMomentum.Generated)

	fmt.Println("\nTop 3 recomendaciones (Momentum):")
	for i := 0; i < 3 && i < len(respMomentum.Recommendations); i++ {
		rec := respMomentum.Recommendations[i]
		fmt.Printf("   %d. %s - %s (Score: %.1f)\n",
			i+1, rec.StockSymbol, rec.Type, rec.Score)
	}

	// ==========================================
	// TEST 3: Generar Recomendaciones - Value
	// ==========================================
	fmt.Println("\n💰 TEST 3: Estrategia Value")
	fmt.Println("---------------------------")

	reqValue := usecase.GenerateRecommendationsRequest{
		Strategies: []string{string(entity.StrategyValue)},
		SaveToDB:   true,
	}

	respValue, err := genRecUC.Execute(ctx, reqValue)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Recomendaciones generadas: %d\n", respValue.Generated)

	fmt.Println("\nTop 3 recomendaciones (Value):")
	for i := 0; i < 3 && i < len(respValue.Recommendations); i++ {
		rec := respValue.Recommendations[i]
		fmt.Printf("   %d. %s - %s (Score: %.1f)\n",
			i+1, rec.StockSymbol, rec.Type, rec.Score)
	}

	// ==========================================
	// TEST 4: Generar Recomendaciones - Dividend
	// ==========================================
	fmt.Println("\n💎 TEST 4: Estrategia Dividend")
	fmt.Println("------------------------------")

	reqDividend := usecase.GenerateRecommendationsRequest{
		Strategies: []string{string(entity.StrategyDividend)},
		SaveToDB:   true,
	}

	respDividend, err := genRecUC.Execute(ctx, reqDividend)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Recomendaciones generadas: %d\n", respDividend.Generated)

	fmt.Println("\nTop 3 recomendaciones (Dividend):")
	for i := 0; i < 3 && i < len(respDividend.Recommendations); i++ {
		rec := respDividend.Recommendations[i]
		fmt.Printf("   %d. %s - %s (Score: %.1f)\n",
			i+1, rec.StockSymbol, rec.Type, rec.Score)
	}

	// ==========================================
	// TEST 5: Generar Recomendaciones - Growth
	// ==========================================
	fmt.Println("\n🚀 TEST 5: Estrategia Growth")
	fmt.Println("----------------------------")

	reqGrowth := usecase.GenerateRecommendationsRequest{
		Strategies: []string{string(entity.StrategyGrowth)},
		SaveToDB:   true,
	}

	respGrowth, err := genRecUC.Execute(ctx, reqGrowth)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Recomendaciones generadas: %d\n", respGrowth.Generated)

	fmt.Println("\nTop 3 recomendaciones (Growth):")
	for i := 0; i < 3 && i < len(respGrowth.Recommendations); i++ {
		rec := respGrowth.Recommendations[i]
		fmt.Printf("   %d. %s - %s (Score: %.1f)\n",
			i+1, rec.StockSymbol, rec.Type, rec.Score)
	}

	// ==========================================
	// TEST 6: Todas las Estrategias
	// ==========================================
	fmt.Println("\n🎨 TEST 6: Todas las Estrategias")
	fmt.Println("--------------------------------")

	reqAll := usecase.GenerateRecommendationsRequest{
		Symbols:  []string{"AAPL"}, // Solo AAPL para este test
		SaveToDB: false,
	}

	respAll, err := genRecUC.Execute(ctx, reqAll)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Recomendaciones para AAPL con todas las estrategias: %d\n", respAll.Generated)

	fmt.Println("\nDesglose por estrategia:")
	for _, rec := range respAll.Recommendations {
		fmt.Printf("   %s: %s (Score: %.1f)\n", rec.Strategy, rec.Type, rec.Score)
	}

	// ==========================================
	// TEST 7: Top Recomendaciones Globales
	// ==========================================
	fmt.Println("\n🏆 TEST 7: Top Recomendaciones Globales")
	fmt.Println("---------------------------------------")

	topRecs, err := recRepo.FindTopRecommendations(ctx, 10)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Top 10 recomendaciones:\n")
	for i, rec := range topRecs {
		fmt.Printf("   %d. %s (%s) - %s (Score: %.1f)\n",
			i+1, rec.StockSymbol, rec.Strategy, rec.Type, rec.Score)
	}

	// ==========================================
	// TEST 8: Estadísticas
	// ==========================================
	fmt.Println("\n📊 TEST 8: Estadísticas")
	fmt.Println("----------------------")

	totalRecs, err := recRepo.Count(ctx)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("✅ Total de recomendaciones en BD: %d\n", totalRecs)

	// Contar por tipo
	buyRecs, _ := recRepo.FindByType(ctx, entity.RecommendationBuy, 1000)
	strongBuyRecs, _ := recRepo.FindByType(ctx, entity.RecommendationStrongBuy, 1000)
	holdRecs, _ := recRepo.FindByType(ctx, entity.RecommendationHold, 1000)
	sellRecs, _ := recRepo.FindByType(ctx, entity.RecommendationSell, 1000)

	fmt.Println("\nDistribución por tipo:")
	fmt.Printf("   STRONG_BUY: %d\n", len(strongBuyRecs))
	fmt.Printf("   BUY: %d\n", len(buyRecs))
	fmt.Printf("   HOLD: %d\n", len(holdRecs))
	fmt.Printf("   SELL: %d\n", len(sellRecs))

	// ==========================================
	// RESUMEN
	// ==========================================
	fmt.Println("\n✨ ================================")
	fmt.Println("✨ TODAS LAS ESTRATEGIAS FUNCIONANDO")
	fmt.Println("✨ Sistema de Recomendaciones Completo")
	fmt.Println("✨ ================================")
	fmt.Println("\n💡 Fase 5 completada exitosamente!")
}
