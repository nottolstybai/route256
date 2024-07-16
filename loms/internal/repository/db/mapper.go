package db

import (
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/internal/repository/db/order"
	"route256.ozon.ru/project/loms/internal/repository/db/stock"
)

func toAddOrderStock(orderID int32, items []entity.Item) (params []order.AddOrderStockParams) {
	for _, item := range items {
		params = append(params, order.AddOrderStockParams{
			OrderID: orderID,
			SkuID:   int32(item.SKU),
			Count:   int32(item.Count),
		})
	}
	return
}

func toReserveItems(item entity.Item) stock.ReserveItemsParams {
	return stock.ReserveItemsParams{Sku: int32(item.SKU), ReservedCount: int32(item.Count)}
}

func toEntityOrderInfo(rows []order.GetByOrderIDRow) *entity.OrderInfo {
	if len(rows) == 0 { //if retrieved rows are empty return nil
		return nil
	}

	orderInfo := entity.OrderInfo{
		Status: entity.Status(rows[0].Status),
		User:   int64(rows[0].UserID),
		Items:  make([]entity.Item, 0, len(rows)),
	}
	for _, row := range rows {
		orderInfo.Items = append(orderInfo.Items, entity.Item{
			SKU:   uint32(row.SkuID),
			Count: uint32(row.Count),
		})
	}
	return &orderInfo
}
