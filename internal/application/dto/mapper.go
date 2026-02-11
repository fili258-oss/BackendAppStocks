package dto

import (
	"github.com/marino/stock-analyzer/internal/domain/entity"
)

// ToStockDTO convierte una entidad Stock a DTO
func ToStockDTO(stock *entity.Stock) StockDTO {
	return StockDTO{
		ID:            stock.ID,
		Symbol:        stock.Symbol,
		Name:          stock.Name,
		Exchange:      stock.Exchange,
		Currency:      stock.Currency,
		Price:         stock.Price,
		OpenPrice:     stock.OpenPrice,
		HighPrice:     stock.HighPrice,
		LowPrice:      stock.LowPrice,
		ClosePrice:    stock.ClosePrice,
		Change:        stock.Change,
		ChangePercent: stock.ChangePercent,
		Volume:        stock.Volume,
		MarketCap:     stock.MarketCap,
		PERatio:       stock.PERatio,
		DividendYield: stock.DividendYield,
		Week52High:    stock.Week52High,
		Week52Low:     stock.Week52Low,
		CreatedAt:     stock.CreatedAt,
		UpdatedAt:     stock.UpdatedAt,
	}
}

// ToStockDTOList convierte una lista de entidades a DTOs
func ToStockDTOList(stocks []*entity.Stock) []StockDTO {
	dtos := make([]StockDTO, len(stocks))
	for i, stock := range stocks {
		dtos[i] = ToStockDTO(stock)
	}
	return dtos
}

// ToStockEntity convierte un DTO a entidad Stock
func ToStockEntity(dto StockDTO) (*entity.Stock, error) {
	stock, err := entity.NewStock(dto.Symbol, dto.Name, dto.Exchange)
	if err != nil {
		return nil, err
	}

	stock.ID = dto.ID
	stock.Currency = dto.Currency
	stock.Price = dto.Price
	stock.OpenPrice = dto.OpenPrice
	stock.HighPrice = dto.HighPrice
	stock.LowPrice = dto.LowPrice
	stock.ClosePrice = dto.ClosePrice
	stock.Change = dto.Change
	stock.ChangePercent = dto.ChangePercent
	stock.Volume = dto.Volume
	stock.MarketCap = dto.MarketCap
	stock.PERatio = dto.PERatio
	stock.DividendYield = dto.DividendYield
	stock.Week52High = dto.Week52High
	stock.Week52Low = dto.Week52Low
	stock.CreatedAt = dto.CreatedAt
	stock.UpdatedAt = dto.UpdatedAt

	return stock, nil
}
