package cart_service

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/cart/internal/client/product_service"
	"route256.ozon.ru/project/cart/internal/config"
	"route256.ozon.ru/project/cart/internal/entity"
	"route256.ozon.ru/project/cart/internal/repository/memory"
	"route256.ozon.ru/project/cart/internal/service"
)

type Suit struct {
	suite.Suite
	prodClient  *product_service.ProductServiceClient
	storage     *memory.MemoryRepo
	cartService *service.CartService
}

func (s *Suit) SetupSuite() {
	var err error
	cfg := config.NewConfig()

	s.prodClient = product_service.NewProductServiceClient(cfg.ProductServiceHost, cfg.ProductServiceToken)
	require.NoError(s.T(), err)

	s.storage = memory.NewMemoryRepo()
	s.cartService = service.NewCartService(s.prodClient, s.storage, nil) //todo add mock for loms client
}

func (s *Suit) TestAdd() {
	ctx := context.Background()

	addItemData := AddItemTestData{
		userID: 1,
		item: entity.Item{
			SkuID: 773297411,
			Count: 10,
		},
	}
	TestAddItemHelper(ctx, s, addItemData)
}

func (s *Suit) TestDelete() {
	ctx := context.Background()

	addItemData := AddItemTestData{
		userID: 1,
		item: entity.Item{
			SkuID: 773297411,
			Count: 10,
		},
	}
	TestAddItemHelper(ctx, s, addItemData)

	err := s.cartService.DeleteItem(ctx, addItemData.userID, addItemData.item.SkuID)
	s.Require().NoError(err)
}

func (s *Suit) TestList() {
	ctx := context.Background()

	userID := int64(1)

	addItemData1 := AddItemTestData{
		userID: userID,
		item: entity.Item{
			SkuID: 773297411,
			Count: 10,
		},
	}

	addItemData2 := AddItemTestData{
		userID: userID,
		item: entity.Item{
			SkuID: 2958025,
			Count: 10,
		},
	}
	TestAddItemHelper(ctx, s, addItemData1)
	TestAddItemHelper(ctx, s, addItemData1)
	TestAddItemHelper(ctx, s, addItemData2)

	actualItemList, err := s.cartService.GetItemsByUserID(ctx, userID)
	s.Require().NoError(err)

	expectedItemList := &entity.ListItems{
		Items: []entity.Item{
			{
				SkuID: 773297411,
				Count: 20,
				Name:  "Кроссовки Nike JORDAN",
				Price: 2202,
			},
			{
				SkuID: 2958025,
				Count: 10,
				Name:  "Roxy Music. Stranded. Remastered Edition",
				Price: 1028,
			},
		},
		TotalPrice: 20*2202 + 10*1028,
	}

	s.Equal(expectedItemList, actualItemList)
}
