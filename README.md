# Stock Analyzer - Sistema de Análisis y Recomendación de Inversiones

Sistema backend para análisis de mercados de valores y recomendaciones de inversión, construido con Clean Architecture y principios SOLID.

## Arquitectura

Este proyecto implementa **Clean Architecture** con las siguientes capas:

```
┌─────────────────────────────────────┐
│   Presentation Layer (HTTP/REST)    │
│   - Handlers, Middleware, Routes    │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│     Application Layer (Use Cases)   │
│   - Business orchestration          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Domain Layer (Entities)        │
│   - Business rules & logic          │
│   - Repository interfaces           │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Infrastructure Layer (External)   │
│   - Database, APIs, Implementations │
└─────────────────────────────────────┘
```

## Principios SOLID Aplicados

### Single Responsibility Principle (SRP)
- Cada entidad tiene una única responsabilidad
- Use cases separados por funcionalidad específica

### Open/Closed Principle (OCP)
- Strategy Pattern para algoritmos de recomendación extensibles
- Factory Pattern para creación de clientes API

### Liskov Substitution Principle (LSP)
- Interfaces de repositorios intercambiables
- Estrategias de recomendación intercambiables

### Interface Segregation Principle (ISP)
- Interfaces pequeñas y específicas
- No se fuerzan dependencias innecesarias

### Dependency Inversion Principle (DIP)
- Domain no depende de Infrastructure
- Inyección de dependencias por constructor

## Patrones de Diseño

### Repository Pattern
Abstracción del acceso a datos. Las interfaces están en `domain/repository` y las implementaciones en `infrastructure/repository`.

### Strategy Pattern
Múltiples algoritmos de recomendación intercambiables en `domain/service`.

### Factory Pattern
Creación de clientes API y estrategias de recomendación.

### Dependency Injection
Todas las dependencias se inyectan por constructor.

## Stack Tecnológico

- **Backend**: Go 1.22
- **Database**: CockroachDB
- **Architecture**: Clean Architecture + DDD
- **Patterns**: Repository, Strategy, Factory, Dependency Injection

## Estructura del Proyecto

```
stock-analyzer/
├── cmd/api/                          # Entry point de la aplicación
├── internal/
│   ├── domain/                       # FASE 1 COMPLETADA
│   │   ├── entity/                   # Entidades de negocio
│   │   │   ├── stock.go             # Entidad Stock
│   │   │   ├── recommendation.go    # Entidad Recommendation
│   │   │   ├── value_objects.go     # Value Objects
│   │   │   └── errors.go            # Errores del dominio
│   │   ├── repository/               # Interfaces de repositorios
│   │   │   ├── stock_repository.go
│   │   │   └── recommendation_repository.go
│   │   └── service/                  # Domain Services
│   │       └── recommendation_strategy.go
│   ├── application/                  # Use Cases (Próxima fase)
│   ├── infrastructure/               # Implementaciones externas
│   └── presentation/                 # Capa HTTP
├── pkg/                              # Código reutilizable
└── test/                             # Tests
```

## Estado del Desarrollo

### Fase 1: Domain Layer (COMPLETADA)

**Entidades:**
- Stock: Entidad principal con métodos de negocio
- Recommendation: Recomendaciones con sistema de scoring
- Value Objects: TimeRange, PriceRange
- Errores de dominio

**Repository Interfaces:**
- StockRepository: CRUD + búsqueda + filtros
- RecommendationRepository: Gestión de recomendaciones

**Domain Services:**
- RecommendationStrategy: Interfaz para estrategias de análisis
- RecommendationScorer: Interfaz para cálculo de scores

**Características del Domain Layer:**
- Independiente de frameworks externos
- Lógica de negocio pura
- Validaciones en las entidades
- Tipos enumerados para categorización
- Cálculos de métricas financieras

### Próximas Fases

- **Fase 2**: Infrastructure - Database (CockroachDB connection + migrations)
- **Fase 3**: Infrastructure - External APIs (Alpha Vantage/Finnhub)
- **Fase 4**: Application Layer - Use Cases
- **Fase 5**: Domain Services - Recommendation Algorithms
- **Fase 6**: Presentation Layer - HTTP API
- **Fase 7**: Testing
- **Fase 8**: Documentation & Refinement

## Configuración

Copiar `.env.example` a `.env` y configurar las variables:

```bash
cp .env.example .env
```

Variables principales:
- `DB_HOST`, `DB_PORT`: Configuración de CockroachDB
- `ALPHA_VANTAGE_API_KEY`: API key para Alpha Vantage
- `FINNHUB_API_KEY`: API key para Finnhub

## Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con coverage
go test -cover ./...

# Tests de una capa específica
go test ./internal/domain/...
```

## Convenciones de Código

- **Idioma**: Código en inglés, comentarios en español
- **Naming**: camelCase para variables, PascalCase para tipos/funciones exportadas
- **Comentarios**: Cada función pública documentada
- **Errores**: Errores tipados en el dominio

## Principios del Dominio

### Entidad Stock
- Representa una acción en el mercado de valores
- Incluye métricas fundamentales y técnicas
- Cálculos automáticos de variación porcentual
- Métodos de validación de datos

### Entidad Recommendation
- Sistema de scoring 0-100
- Tipos de recomendación: BUY, SELL, HOLD, STRONG_BUY
- Cálculo de confianza basado en cercanía a umbrales
- Validez temporal (24 horas por defecto)
- Métricas rastreables en el análisis

## Referencias

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

## Autor

Marino Botina - Software Engineer

## Licencia

Este es un proyecto de prueba técnica para proceso de entrevista.
