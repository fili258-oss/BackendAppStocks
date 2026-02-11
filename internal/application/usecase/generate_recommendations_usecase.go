package usecase

import (
	"context"
	"fmt"

	"github.com/marino/stock-analyzer/internal/domain/entity"
	"github.com/marino/stock-analyzer/internal/domain/repository"
	"github.com/marino/stock-analyzer/internal/domain/service"
)

// GenerateRecommendationsRequest request para generar recomendaciones
type GenerateRecommendationsRequest struct {
	Symbols   []string // Símbolos específicos (opcional)
	Strategies []string // Estrategias a usar (opcional, default: todas)
	SaveToDB  bool     // Guardar en base de datos
}

// GenerateRecommendationsResponse respuesta con recomendaciones generadas
type GenerateRecommendationsResponse struct {
	Recommendations []*entity.Recommendation
	Generated       int
	Failed          int
	Errors          []string
}

// GenerateRecommendationsUseCase genera recomendaciones usando estrategias
type GenerateRecommendationsUseCase struct {
	stockRepo repository.StockRepository
	recRepo   repository.RecommendationRepository
	strategies map[string]service.RecommendationStrategy
}

// NewGenerateRecommendationsUseCase crea una nueva instancia
func NewGenerateRecommendationsUseCase(
	stockRepo repository.StockRepository,
	recRepo repository.RecommendationRepository,
) *GenerateRecommendationsUseCase {
	// Inicializar todas las estrategias
	strategies := make(map[string]service.RecommendationStrategy)
	strategies[string(entity.StrategyBalanced)] = service.NewBalancedStrategy()
	strategies[string(entity.StrategyMomentum)] = service.NewMomentumStrategy()
	strategies[string(entity.StrategyValue)] = service.NewValueStrategy()
	strategies[string(entity.StrategyDividend)] = service.NewDividendStrategy()
	strategies[string(entity.StrategyGrowth)] = service.NewGrowthStrategy()

	return &GenerateRecommendationsUseCase{
		stockRepo:  stockRepo,
		recRepo:    recRepo,
		strategies: strategies,
	}
}

// Execute ejecuta el caso de uso
func (uc *GenerateRecommendationsUseCase) Execute(
	ctx context.Context,
	request GenerateRecommendationsRequest,
) (*GenerateRecommendationsResponse, error) {
	response := &GenerateRecommendationsResponse{
		Recommendations: make([]*entity.Recommendation, 0),
		Errors:          make([]string, 0),
	}

	// Obtener stocks a analizar
	var stocks []*entity.Stock
	var err error

	if len(request.Symbols) > 0 {
		// Obtener stocks específicos
		for _, symbol := range request.Symbols {
			stock, err := uc.stockRepo.FindBySymbol(ctx, symbol)
			if err != nil {
				response.Errors = append(response.Errors, 
					fmt.Sprintf("Stock %s not found: %v", symbol, err))
				response.Failed++
				continue
			}
			stocks = append(stocks, stock)
		}
	} else {
		// Obtener todos los stocks
		stocks, err = uc.stockRepo.FindAll(ctx, 1000, 0)
		if err != nil {
			return nil, fmt.Errorf("error fetching stocks: %w", err)
		}
	}

	if len(stocks) == 0 {
		return response, nil
	}

	// Determinar qué estrategias usar
	strategiesToUse := uc.getStrategiesToUse(request.Strategies)

	// Generar recomendaciones para cada stock con cada estrategia
	for _, stock := range stocks {
		for _, strategyName := range strategiesToUse {
			strategy, exists := uc.strategies[strategyName]
			if !exists {
				continue
			}

			rec, err := strategy.Analyze(ctx, stock)
			if err != nil {
				response.Errors = append(response.Errors,
					fmt.Sprintf("Error analyzing %s with %s: %v", stock.Symbol, strategyName, err))
				response.Failed++
				continue
			}

			// Guardar en BD si se solicita
			if request.SaveToDB {
				if err := uc.recRepo.Create(ctx, rec); err != nil {
					response.Errors = append(response.Errors,
						fmt.Sprintf("Error saving recommendation: %v", err))
					response.Failed++
					continue
				}
			}

			response.Recommendations = append(response.Recommendations, rec)
			response.Generated++
		}
	}

	return response, nil
}

// getStrategiesToUse determina qué estrategias usar
func (uc *GenerateRecommendationsUseCase) getStrategiesToUse(requested []string) []string {
	if len(requested) == 0 {
		// Por defecto, usar todas
		return []string{
			string(entity.StrategyBalanced),
			string(entity.StrategyMomentum),
			string(entity.StrategyValue),
			string(entity.StrategyDividend),
			string(entity.StrategyGrowth),
		}
	}
	return requested
}

// GetTopRecommendations obtiene las mejores recomendaciones
func (uc *GenerateRecommendationsUseCase) GetTopRecommendations(
	ctx context.Context,
	strategy string,
	limit int,
) ([]*entity.Recommendation, error) {
	if strategy == "" {
		return uc.recRepo.FindTopRecommendations(ctx, limit)
	}

	// Filtrar por estrategia específica
	// Nota: Necesitaríamos un método en el repositorio para esto
	// Por ahora, retornamos las top generales
	return uc.recRepo.FindTopRecommendations(ctx, limit)
}
