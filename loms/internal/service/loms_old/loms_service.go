package loms_old

import (
	"context"
	"errors"
	"log"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/internal/repository/memory/stock"
)

var ErrReservation = errors.New("error while creating reservation")

//go:generate minimock [-i Repository, OrderRepository, StocksRepository] [-o mock_order_repo_test.go, mock_stock_repo_test.go]
type OrderRepository interface {
	Create(ctx context.Context, userID int64, items []entity.Item) (entity.OrderID, error)
	SetStatus(ctx context.Context, orderID entity.OrderID, status entity.Status) error
	GetByOrderID(ctx context.Context, orderID entity.OrderID) (*entity.OrderInfo, error)
}

type StocksRepository interface {
	Reserve(ctx context.Context, items []entity.Item) error
	RemoveReservation(ctx context.Context, items []entity.Item) error
	CancelReservation(ctx context.Context, items []entity.Item) error
	GetBySku(ctx context.Context, sku uint32) (entity.Count, error)
}

type LOMSService struct {
	orderRepo  OrderRepository
	stocksRepo StocksRepository
}

func NewLOMSService(orderRepo OrderRepository, stocksRepo StocksRepository) *LOMSService {
	return &LOMSService{orderRepo: orderRepo, stocksRepo: stocksRepo}
}

func (ls *LOMSService) OrderCreate(ctx context.Context, userID int64, items []entity.Item) (entity.OrderID, error) {
	orderID, err := ls.orderRepo.Create(ctx, userID, items)
	if err != nil {
		return 0, err
	}

	if err := ls.stocksRepo.Reserve(ctx, items); err != nil {
		log.Printf("Error creating reservation=%v", err)

		if errors.Is(err, stock.ErrReserveTooLarge) {
			if err := ls.orderRepo.SetStatus(ctx, orderID, entity.StatusFailed); err != nil {
				return 0, err
			}
		}
		return 0, ErrReservation
	}

	if err := ls.orderRepo.SetStatus(ctx, orderID, entity.StatusAwaitingPayment); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (ls *LOMSService) OrderInfo(ctx context.Context, orderID entity.OrderID) (*entity.OrderInfo, error) {
	orderInfo, err := ls.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("Error fetching order by id=%d: %v", orderID, err)
		return nil, err
	}
	return orderInfo, nil
}

func (ls *LOMSService) OrderPay(ctx context.Context, orderID entity.OrderID) error {
	orderInfo, err := ls.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("Error fetching order by id=%d: %v", orderID, err)
		return err
	}

	if err := ls.stocksRepo.RemoveReservation(ctx, orderInfo.Items); err != nil {
		log.Printf("Error removing reservation by order_id=%d: %v", orderID, err)
		return err
	}

	if err := ls.orderRepo.SetStatus(ctx, orderID, entity.StatusPayed); err != nil {
		log.Printf("Error setting status: %v", err)
		return err
	}
	return nil
}

func (ls *LOMSService) OrderCancel(ctx context.Context, orderID entity.OrderID) error {
	orderInfo, err := ls.orderRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("Error fetching order by id=%d: %v", orderID, err)
		return err
	}

	if err := ls.stocksRepo.CancelReservation(ctx, orderInfo.Items); err != nil {
		log.Printf("Error cancelling reservation by order_id=%d: %v", orderID, err)
		return err
	}

	if err := ls.orderRepo.SetStatus(ctx, orderID, entity.StatusCancelled); err != nil {
		log.Printf("Error setting status: %v", err)
		return err
	}
	return nil
}

func (ls *LOMSService) StocksInfo(ctx context.Context, sku uint32) (entity.Count, error) {
	count, err := ls.stocksRepo.GetBySku(ctx, sku)
	if err != nil {
		log.Printf("Error getting info about stocks with sku=%d: %v", sku, err)
		return 0, err
	}
	return count, nil
}
