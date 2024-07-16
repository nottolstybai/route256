package loms_old

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/internal/repository/memory/order"
	memory2 "route256.ozon.ru/project/loms/internal/repository/memory/stock"
	"testing"
)

func TestLOMSService_OrderCreate(t *testing.T) {
	t.Parallel()

	type input struct {
		userID int64
		items  []entity.Item
	}

	type want struct {
		orderID         entity.OrderID
		status          entity.Status
		orderCreateErr  error
		stockReserveErr error
		setStatusErr    error
		err             error
	}

	ErrOrderCreate := errors.New("order creation failed")
	ErrSetStatusFailed := errors.New("failed to set status")

	mc := minimock.NewController(t)
	mockOrderRepo := NewOrderRepositoryMock(mc)
	mockStocksRepo := NewStocksRepositoryMock(mc)

	lomsService := NewLOMSService(mockOrderRepo, mockStocksRepo)

	testCases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Successful order creation",
			input: input{
				userID: 1,
				items:  []entity.Item{{SKU: 1, Count: 10}},
			},
			want: want{
				orderID: 123,
				status:  entity.StatusAwaitingPayment,
			},
		},
		{
			name: "Order create failed",
			input: input{
				userID: 1,
				items:  []entity.Item{{SKU: 1, Count: 1000}},
			},
			want: want{
				orderCreateErr: ErrOrderCreate,
				err:            ErrOrderCreate,
			},
		},
		{
			name: "Reserve failed with ErrReserveTooLarge",
			input: input{
				userID: 1,
				items:  []entity.Item{{SKU: 1, Count: 1000}},
			},
			want: want{
				status:          entity.StatusFailed,
				orderCreateErr:  nil,
				stockReserveErr: memory2.ErrReserveTooLarge,
				setStatusErr:    ErrSetStatusFailed,
				err:             ErrSetStatusFailed,
			},
		},
		{
			name: "Set status failed",
			input: input{
				userID: 1,
				items:  []entity.Item{{SKU: 1, Count: 1000}},
			},
			want: want{
				status:       entity.StatusAwaitingPayment,
				setStatusErr: ErrSetStatusFailed,
				err:          ErrSetStatusFailed,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			ctx := context.Background()
			mockOrderRepo.CreateMock.
				Expect(ctx, tc.input.userID, tc.input.items).
				Return(tc.want.orderID, tc.want.orderCreateErr)
			mockStocksRepo.ReserveMock.
				Expect(ctx, tc.input.items).
				Return(tc.want.stockReserveErr)
			mockOrderRepo.SetStatusMock.
				Expect(ctx, tc.want.orderID, tc.want.status).
				Return(tc.want.setStatusErr)

			orderID, err := lomsService.OrderCreate(ctx, tc.input.userID, tc.input.items)
			require.ErrorIs(t, err, tc.want.err)
			require.Equal(t, tc.want.orderID, orderID)
		})
	}
}

func TestLOMSService_OrderInfo(t *testing.T) {
	t.Parallel()

	type input struct {
		orderID entity.OrderID
	}

	type want struct {
		orderInfo *entity.OrderInfo
		err       error
	}

	mc := minimock.NewController(t)
	mockOrderRepo := NewOrderRepositoryMock(mc)
	mockStocksRepo := NewStocksRepositoryMock(mc)

	lomsService := NewLOMSService(mockOrderRepo, mockStocksRepo)

	testCases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name:  "Successful order info",
			input: input{orderID: entity.OrderID(123)},
			want: want{
				orderInfo: &entity.OrderInfo{
					Status: entity.StatusNew,
					Items:  []entity.Item{{SKU: 1, Count: 1000}},
				},
			},
		},
		{
			name:  "Invalid orderID",
			input: input{orderID: entity.OrderID(10000)},
			want:  want{err: order.ErrOrderNotFound},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			ctx := context.Background()
			mockOrderRepo.GetByOrderIDMock.Expect(ctx, tc.input.orderID).Return(tc.want.orderInfo, tc.want.err)

			orderInfo, err := lomsService.OrderInfo(ctx, tc.input.orderID)
			require.ErrorIs(t, err, tc.want.err)
			require.Equal(t, tc.want.orderInfo, orderInfo)
		})
	}
}

func TestLOMSService_OrderPay(t *testing.T) {
	t.Parallel()

	ErrSetStatusFailed := errors.New("failed to set status")

	type want struct {
		orderInfo    *entity.OrderInfo
		orderIDErr   error
		setStatusErr error
		err          error
	}

	mc := minimock.NewController(t)
	mockOrderRepo := NewOrderRepositoryMock(mc)
	mockStocksRepo := NewStocksRepositoryMock(mc)

	lomsService := NewLOMSService(mockOrderRepo, mockStocksRepo)

	testCases := []struct {
		name    string
		orderID entity.OrderID
		status  entity.Status
		want    want
	}{
		{
			name:    "Successful order pay",
			orderID: entity.OrderID(123),
			status:  entity.StatusPayed,
			want: want{
				orderInfo: &entity.OrderInfo{
					Status: entity.StatusAwaitingPayment,
					Items:  []entity.Item{{SKU: 1, Count: 1000}},
				},
			},
		},
		{
			name:    "invalid OrderID",
			orderID: entity.OrderID(10000),
			want: want{
				orderInfo: &entity.OrderInfo{
					Items: nil,
				},
				orderIDErr: order.ErrOrderNotFound,
				err:        order.ErrOrderNotFound},
		},
		{
			name:    "set status Payed failed",
			orderID: entity.OrderID(123),
			status:  entity.StatusPayed,
			want: want{
				orderInfo: &entity.OrderInfo{
					Status: entity.StatusAwaitingPayment,
					Items:  []entity.Item{{SKU: 1, Count: 1000}},
				},
				setStatusErr: ErrSetStatusFailed,
				err:          ErrSetStatusFailed,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			ctx := context.Background()
			mockOrderRepo.GetByOrderIDMock.Expect(ctx, tc.orderID).Return(tc.want.orderInfo, tc.want.orderIDErr)
			mockStocksRepo.RemoveReservationMock.Expect(ctx, tc.want.orderInfo.Items).Return(nil)
			mockOrderRepo.SetStatusMock.Expect(ctx, tc.orderID, tc.status).Return(tc.want.setStatusErr)
			err := lomsService.OrderPay(ctx, tc.orderID)
			require.ErrorIs(t, err, tc.want.err)
		})
	}
}

func TestLOMSService_OrderCancel(t *testing.T) {
	t.Parallel()

	ErrSetStatusFailed := errors.New("failed to set status")
	ErrCancelReserveFailed := errors.New("cancel failed")

	type input struct {
		orderID entity.OrderID
		status  entity.Status
	}

	type want struct {
		orderInfo    *entity.OrderInfo
		orderIDErr   error
		cancelErr    error
		setStatusErr error
		err          error
	}

	mc := minimock.NewController(t)
	mockOrderRepo := NewOrderRepositoryMock(mc)
	mockStocksRepo := NewStocksRepositoryMock(mc)

	lomsService := NewLOMSService(mockOrderRepo, mockStocksRepo)

	testCases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Successful order cancel",
			input: input{
				orderID: entity.OrderID(123),
				status:  entity.StatusCancelled,
			},
			want: want{
				orderInfo: &entity.OrderInfo{
					Status: entity.StatusAwaitingPayment,
					Items:  []entity.Item{{SKU: 1, Count: 1000}},
				},
			},
		},
		{
			name: "invalid OrderID",
			input: input{
				orderID: entity.OrderID(10000),
			},
			want: want{
				orderInfo: &entity.OrderInfo{
					Items: nil,
				},
				orderIDErr: order.ErrOrderNotFound,
				err:        order.ErrOrderNotFound},
		},
		{
			name: "reservation cancel failed",
			input: input{
				orderID: entity.OrderID(123),
			},
			want: want{
				orderInfo: &entity.OrderInfo{
					Status: entity.StatusAwaitingPayment,
					Items:  []entity.Item{{SKU: 1, Count: 1000}},
				},
				cancelErr: ErrCancelReserveFailed,
				err:       ErrCancelReserveFailed,
			},
		},
		{
			name: "set status Payed failed",
			input: input{
				orderID: entity.OrderID(123),
				status:  entity.StatusCancelled,
			},
			want: want{
				orderInfo: &entity.OrderInfo{
					Status: entity.StatusAwaitingPayment,
					Items:  []entity.Item{{SKU: 1, Count: 1000}},
				},
				setStatusErr: ErrSetStatusFailed,
				err:          ErrSetStatusFailed,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			ctx := context.Background()
			mockOrderRepo.GetByOrderIDMock.Expect(ctx, tc.input.orderID).Return(tc.want.orderInfo, tc.want.orderIDErr)
			mockStocksRepo.CancelReservationMock.Expect(ctx, tc.want.orderInfo.Items).Return(tc.want.cancelErr)
			mockOrderRepo.SetStatusMock.Expect(ctx, tc.input.orderID, tc.input.status).Return(tc.want.setStatusErr)
			err := lomsService.OrderCancel(ctx, tc.input.orderID)
			require.ErrorIs(t, err, tc.want.err)
		})
	}
}
