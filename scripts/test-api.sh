#!/bin/bash

# Script de prueba de la API
# Asume que el servidor está corriendo en localhost:8000

BASE_URL="http://localhost:8000"
API_URL="$BASE_URL/api/v1"

echo "=== Test de Stock Analyzer API ==="
echo ""

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
curl -s "$BASE_URL/health" | jq .
echo ""

# Test 2: Welcome
echo -e "${YELLOW}Test 2: Welcome / Endpoints${NC}"
curl -s "$BASE_URL/" | jq .
echo ""

# Test 3: Buscar stocks
echo -e "${YELLOW}Test 3: Search Stocks${NC}"
curl -s "$API_URL/stocks?query=Apple&limit=5" | jq .
echo ""

# Test 4: Obtener stock específico
echo -e "${YELLOW}Test 4: Get Stock Details (AAPL)${NC}"
curl -s "$API_URL/stocks/AAPL" | jq .
echo ""

# Test 5: Fetch stocks desde Finnhub
echo -e "${YELLOW}Test 5: Fetch Stocks from API${NC}"
curl -s -X POST "$API_URL/stocks/fetch" \
  -H "Content-Type: application/json" \
  -d '{
    "symbols": ["MSFT"],
    "save": true
  }' | jq .
echo ""

# Test 6: Generar recomendaciones
echo -e "${YELLOW}Test 6: Generate Recommendations${NC}"
curl -s -X POST "$API_URL/recommendations/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "symbols": ["AAPL"],
    "strategies": ["BALANCED", "MOMENTUM"],
    "save_to_db": true
  }' | jq .
echo ""

# Test 7: Top recomendaciones
echo -e "${YELLOW}Test 7: Get Top Recommendations${NC}"
curl -s "$API_URL/recommendations/top?limit=5" | jq .
echo ""

# Test 8: Recomendaciones por stock
echo -e "${YELLOW}Test 8: Get Recommendations by Stock (AAPL)${NC}"
curl -s "$API_URL/recommendations/stock/AAPL" | jq .
echo ""

# Test 9: Recomendaciones por tipo
echo -e "${YELLOW}Test 9: Get Recommendations by Type (STRONG_BUY)${NC}"
curl -s "$API_URL/recommendations/type/STRONG_BUY?limit=3" | jq .
echo ""

# Test 10: Estadísticas
echo -e "${YELLOW}Test 10: Get Recommendation Stats${NC}"
curl -s "$API_URL/recommendations/stats" | jq .
echo ""

echo -e "${GREEN}✅ All tests completed!${NC}"
