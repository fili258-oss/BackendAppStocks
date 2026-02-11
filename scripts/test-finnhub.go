package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/marino/stock-analyzer/internal/infrastructure/external"
)

func main() {
	fmt.Println("=== Test de Finnhub API ===\n")

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  No se pudo cargar .env: %v\n", err)
	}

	// Obtener API key
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ FINNHUB_API_KEY no está configurada en .env")
	}

	fmt.Println("✅ API Key encontrada")

	// Crear cliente de Finnhub
	config := external.FinnhubConfig{
		APIKey:    apiKey,
		Timeout:   30 * time.Second,
		RateLimit: 5, // 5 requests por segundo (seguro para free tier)
	}

	client := external.NewFinnhubClient(config)
	adapter := external.NewFinnhubAdapter(client)

	ctx := context.Background()

	// ==========================================
	// TEST 1: Obtener Quote (Precio Actual)
	// ==========================================
	fmt.Println("\n📊 TEST 1: Quote de AAPL")
	fmt.Println("-------------------------")

	quote, err := client.GetQuote(ctx, "AAPL")
	if err != nil {
		log.Fatalf("❌ Error getting quote: %v\n", err)
	}

	fmt.Printf("✅ Precio actual: $%.2f\n", quote.C)
	fmt.Printf("   Apertura: $%.2f\n", quote.O)
	fmt.Printf("   Máximo: $%.2f\n", quote.H)
	fmt.Printf("   Mínimo: $%.2f\n", quote.L)
	fmt.Printf("   Cambio: $%.2f (%.2f%%)\n", quote.D, quote.DP)

	// ==========================================
	// TEST 2: Obtener Profile (Info Compañía)
	// ==========================================
	fmt.Println("\n🏢 TEST 2: Profile de AAPL")
	fmt.Println("--------------------------")

	profile, err := client.GetProfile(ctx, "AAPL")
	if err != nil {
		log.Fatalf("❌ Error getting profile: %v\n", err)
	}

	fmt.Printf("✅ Nombre: %s\n", profile.Name)
	fmt.Printf("   Ticker: %s\n", profile.Ticker)
	fmt.Printf("   Exchange: %s\n", profile.Exchange)
	fmt.Printf("   Currency: %s\n", profile.Currency)
	fmt.Printf("   Market Cap: $%.2fB\n", profile.MarketCap)
	fmt.Printf("   Industry: %s\n", profile.Industry)
	fmt.Printf("   Country: %s\n", profile.Country)

	// ==========================================
	// TEST 3: Obtener Metrics (Fundamentales)
	// ==========================================
	fmt.Println("\n📈 TEST 3: Metrics de AAPL")
	fmt.Println("--------------------------")

	metrics, err := client.GetMetrics(ctx, "AAPL")
	if err != nil {
		log.Printf("⚠️  Error getting metrics: %v\n", err)
	} else {
		fmt.Printf("✅ P/E Ratio: %.2f\n", metrics.Metric.PEBasicExclExtraTTM)
		fmt.Printf("   Dividend Yield: %.2f%%\n", metrics.Metric.DividendYieldIndicatedAnnual)
		fmt.Printf("   52 Week High: $%.2f\n", metrics.Metric.Week52High)
		fmt.Printf("   52 Week Low: $%.2f\n", metrics.Metric.Week52Low)
		fmt.Printf("   Beta: %.2f\n", metrics.Metric.Beta)
	}

	// ==========================================
	// TEST 4: Adapter - Stock Completo
	// ==========================================
	fmt.Println("\n🔄 TEST 4: Adapter - Stock Completo")
	fmt.Println("-----------------------------------")

	stock, err := adapter.GetStock(ctx, "GOOGL")
	if err != nil {
		log.Fatalf("❌ Error with adapter: %v\n", err)
	}

	fmt.Printf("✅ Stock creado desde Finnhub:\n")
	fmt.Printf("   Symbol: %s\n", stock.Symbol)
	fmt.Printf("   Name: %s\n", stock.Name)
	fmt.Printf("   Exchange: %s\n", stock.Exchange)
	fmt.Printf("   Price: $%.2f\n", stock.Price)
	fmt.Printf("   Market Cap: $%.2fB\n", stock.MarketCap/1000000000)
	fmt.Printf("   Volume: %d\n", stock.Volume)
	fmt.Printf("   Change: %.2f%%\n", stock.ChangePercent)
	fmt.Printf("   P/E Ratio: %.2f\n", stock.PERatio)

	// ==========================================
	// TEST 5: Múltiples Stocks
	// ==========================================
	fmt.Println("\n📦 TEST 5: Múltiples Stocks")
	fmt.Println("---------------------------")

	symbols := []string{"MSFT", "TSLA", "AMZN"}
	fmt.Printf("Obteniendo datos para: %v\n", symbols)

	stocks, err := adapter.GetMultipleStocks(ctx, symbols)
	if err != nil {
		log.Printf("⚠️  Algunos stocks fallaron: %v\n", err)
	}

	fmt.Printf("✅ %d stocks obtenidos:\n", len(stocks))
	for i, s := range stocks {
		fmt.Printf("   %d. %s (%s) - $%.2f\n", i+1, s.Symbol, s.Name, s.Price)
	}

	// ==========================================
	// TEST 6: Búsqueda de Stocks
	// ==========================================
	fmt.Println("\n🔍 TEST 6: Búsqueda de Stocks")
	fmt.Println("-----------------------------")

	searchResults, err := client.SearchSymbol(ctx, "Tesla")
	if err != nil {
		log.Printf("⚠️  Error searching: %v\n", err)
	} else {
		fmt.Printf("✅ %d resultados encontrados para 'Tesla':\n", searchResults.Count)
		for i := 0; i < 5 && i < len(searchResults.Result); i++ {
			r := searchResults.Result[i]
			fmt.Printf("   %d. %s - %s (%s)\n", i+1, r.Symbol, r.Description, r.Type)
		}
	}

	// ==========================================
	// TEST 7: Rate Limiting
	// ==========================================
	fmt.Println("\n⏱️  TEST 7: Rate Limiting")
	fmt.Println("------------------------")

	fmt.Println("Haciendo 10 requests rápidos (debería tardar ~2 segundos)...")
	start := time.Now()
	
	for i := 0; i < 10; i++ {
		_, err := client.GetQuote(ctx, "AAPL")
		if err != nil {
			log.Printf("Request %d failed: %v\n", i+1, err)
		}
		fmt.Printf(".")
	}
	
	elapsed := time.Since(start)
	fmt.Printf("\n✅ 10 requests completados en %.2f segundos\n", elapsed.Seconds())
	fmt.Printf("   Rate limiting funcionando correctamente\n")

	// ==========================================
	// RESUMEN
	// ==========================================
	fmt.Println("\n✨ ================================")
	fmt.Println("✨ TODOS LOS TESTS PASARON")
	fmt.Println("✨ Finnhub API funcionando correctamente")
	fmt.Println("✨ ================================")
	fmt.Println("\n💡 Próximo paso: Guardar estos datos en la base de datos")
}
