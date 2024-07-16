package cart_service

import (
	"context"
	"route256.ozon.ru/project/cart/internal/entity"
)

type AddItemTestData struct {
	userID int64
	item   entity.Item
}

func TestAddItemHelper(ctx context.Context, suite *Suit, testData AddItemTestData) {
	err := suite.cartService.AddItem(ctx, testData.userID, testData.item.SkuID, testData.item.Count)
	suite.Require().NoError(err)
}
