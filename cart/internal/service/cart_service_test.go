package service

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/cart/internal/client/product_service"
	"route256.ozon.ru/project/cart/internal/entity"
	"route256.ozon.ru/project/cart/pkg/logger"
	"testing"
)

func TestMain(m *testing.M) {
	logger.Init()
	goleak.VerifyTestMain(m)
}

func TestCartService_AddItem(t *testing.T) {
	t.Parallel()

	type inputData struct {
		userID int64
		item   entity.Item
	}

	type wantData struct {
		product    *entity.Product
		productErr error
		insertErr  error
	}

	mc := minimock.NewController(t)

	productServiceMock := NewClientMock(mc)
	repositoryMock := NewRepositoryMock(mc)
	lomsServiceMock := NewLomsClientMock(mc)

	cartService := NewCartService(productServiceMock, repositoryMock, lomsServiceMock)

	testCases := []struct {
		name  string
		input inputData
		want  wantData
	}{
		{
			name: "Valid product and successful insert",
			input: inputData{
				userID: 1,
				item:   entity.Item{SkuID: 123, Count: 2},
			},
			want: wantData{
				product: &entity.Product{},
			},
		},
		{
			name: "Invalid product",
			input: inputData{
				userID: 1,
				item:   entity.Item{SkuID: 123, Count: 2},
			},
			want: wantData{
				productErr: product_service.ErrProductNotFound,
			},
		},
		//TODO: add more testcases when we start using DB
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := context.Background()

			productServiceMock.GetProductMock.
				Expect(ctx, tc.input.item.SkuID).
				Return(tc.want.product, tc.want.productErr)

			repositoryMock.InsertMock.
				Expect(ctx, tc.input.userID, tc.input.item.SkuID, tc.input.item.Count).
				Return(tc.want.insertErr)

			lomsServiceMock.StocksInfoMock.
				Expect(ctx, uint32(tc.input.item.SkuID)).
				Return(tc.input.item.Count, tc.want.productErr)

			err := cartService.AddItem(ctx, tc.input.userID, tc.input.item.SkuID, tc.input.item.Count)
			require.ErrorIs(t, err, tc.want.productErr)
		})
	}
}

func TestCartService_DeleteItem(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	productServiceMock := NewClientMock(mc)
	repositoryMock := NewRepositoryMock(mc)
	lomsServiceMock := NewLomsClientMock(mc)

	cartService := NewCartService(productServiceMock, repositoryMock, lomsServiceMock)

	testCases := []struct {
		name    string
		userID  int64
		skuID   int64
		wantErr error
	}{
		{
			name:   "Item deleted successfully",
			userID: 1,
			skuID:  123,
		},
		{
			name:    "Invalid userID",
			userID:  1,
			skuID:   123,
			wantErr: errors.New("user with specified id doesn't exist"),
		},
		{
			name:    "Invalid skuID",
			userID:  1,
			skuID:   123,
			wantErr: errors.New("product with specified id doesn't exist"),
		},
		//TODO: add more testcases when we start using DB
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			repositoryMock.RemoveMock.Expect(ctx, tc.userID, tc.skuID).Return(tc.wantErr)

			err := cartService.DeleteItem(context.Background(), tc.userID, tc.skuID)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCartService_DeleteAll(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	productServiceMock := NewClientMock(mc)
	repositoryMock := NewRepositoryMock(mc)
	lomsServiceMock := NewLomsClientMock(mc)

	cartService := NewCartService(productServiceMock, repositoryMock, lomsServiceMock)

	testCases := []struct {
		name    string
		userID  int64
		wantErr error
	}{
		{
			name:   "All items deleted successfully",
			userID: 1,
		},
		{
			name:    "Invalid userID",
			userID:  1,
			wantErr: errors.New("user with specified id doesn't exist"),
		},
		//TODO: add more testcases when we start using DB
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()
			ctx := context.Background()
			repositoryMock.RemoveByUserIDMock.Expect(ctx, tc.userID).Return(tc.wantErr)

			err := cartService.DeleteItemsByUserID(context.Background(), tc.userID)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCartService_GetAllItems(t *testing.T) {
	t.Parallel()

	var (
		ErrRepo   = errors.New("repo error")
		ErrClient = errors.New("client error")
	)

	type inputData struct {
		userID int64
		skuId  int64
	}

	type wantData struct {
		repoErr           error
		repoProducts      []entity.ProductInfo
		clientErr         error
		clientProductInfo *entity.Product
		expectedResp      *entity.ListItems
		expectedErr       error
	}

	mc := minimock.NewController(t)

	productServiceMock := NewClientMock(mc)
	repositoryMock := NewRepositoryMock(mc)
	lomsServiceMock := NewLomsClientMock(mc)

	cartService := NewCartService(productServiceMock, repositoryMock, lomsServiceMock)

	testCases := []struct {
		name  string
		input inputData
		want  wantData
	}{
		{
			name: "Items retrieved successfully",
			input: inputData{
				userID: 10,
				skuId:  123,
			},
			want: wantData{
				expectedErr: nil,
				repoProducts: []entity.ProductInfo{
					{SkuID: 123, Count: 2},
				},
				clientProductInfo: &entity.Product{
					Name:  "Product 123",
					Price: 10,
				},
				expectedResp: &entity.ListItems{
					Items: []entity.Item{
						{
							SkuID: 123,
							Count: 2,
							Name:  "Product 123",
							Price: 10,
						},
					},
					TotalPrice: 20,
				},
			},
		},
		{
			name: "Repository error",
			input: inputData{
				userID: 10,
				skuId:  123,
			},
			want: wantData{
				repoErr:     ErrRepo,
				expectedErr: ErrRepo,
			},
		},
		{
			name: "Client error",
			input: inputData{
				userID: 10,
				skuId:  123,
			},
			want: wantData{
				repoProducts: []entity.ProductInfo{
					{SkuID: 123, Count: 2},
				},
				clientErr:   ErrClient,
				expectedErr: ErrClient,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repositoryMock.ListMock.
				Expect(minimock.AnyContext, tc.input.userID).
				Return(tc.want.repoProducts, tc.want.repoErr)
			productServiceMock.GetProductMock.
				Expect(minimock.AnyContext, tc.input.skuId).
				Return(tc.want.clientProductInfo, tc.want.clientErr)

			actualResp, err := cartService.GetItemsByUserID(context.Background(), tc.input.userID)

			require.ErrorIs(t, err, tc.want.expectedErr)
			require.Equal(t, tc.want.expectedResp, actualResp)
		})
	}
}
