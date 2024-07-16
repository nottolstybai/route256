package loms

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/pkg/logger"
)

type Storage interface {
	CreateOrder(ctx context.Context, userID int64, items []entity.Item) (entity.OrderID, error)
	GetOrderInfoByID(ctx context.Context, orderID entity.OrderID) (*entity.OrderInfo, error)
	PayOrder(ctx context.Context, orderID entity.OrderID) error
	CancelOrder(ctx context.Context, orderID entity.OrderID) error
	GetStocksInfoByID(ctx context.Context, sku uint32) (entity.Count, error)
}

type LOMSService struct {
	storage Storage
}

func NewLOMSService(storage Storage) *LOMSService {
	return &LOMSService{storage: storage}
}

func (ls *LOMSService) OrderCreate(ctx context.Context, userID int64, items []entity.Item) (entity.OrderID, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LOMSService.OrderCreate")
	defer span.End()

	orderID, err := ls.storage.CreateOrder(ctx, userID, items)
	if err != nil {
		return 0, logger.WithError(err, "create order failed")
	}
	logger.Debug("create order successfully", zap.Int("orderID", int(orderID)))
	return orderID, nil
}

func (ls *LOMSService) OrderInfo(ctx context.Context, orderID entity.OrderID) (*entity.OrderInfo, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LOMSService.OrderInfo")
	defer span.End()

	orderInfo, err := ls.storage.GetOrderInfoByID(ctx, orderID)
	if err != nil {
		return nil, logger.WithError(err, "getting order info failed")
	}

	logger.Debug("create order successfully", zap.Any("orderInfo", orderInfo))
	return orderInfo, nil
}

func (ls *LOMSService) OrderPay(ctx context.Context, orderID entity.OrderID) error {
	ctx, span := otel.Tracer("default").Start(ctx, "LOMSService.OrderPay")
	defer span.End()

	if err := ls.storage.PayOrder(ctx, orderID); err != nil {
		return logger.WithError(err, "order pay failed")
	}
	logger.Debug("pay order success using orderID", zap.Int("orderID", int(orderID)))
	return nil
}

func (ls *LOMSService) OrderCancel(ctx context.Context, orderID entity.OrderID) error {
	ctx, span := otel.Tracer("default").Start(ctx, "LOMSService.OrderCancel")
	defer span.End()

	if err := ls.storage.CancelOrder(ctx, orderID); err != nil {
		return logger.WithError(err, "order cancel failed")
	}
	logger.Debug("cancel order success using orderID", zap.Int("orderID", int(orderID)))
	return nil
}

func (ls *LOMSService) StocksInfo(ctx context.Context, sku uint32) (entity.Count, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LOMSService.StocksInfo")
	defer span.End()

	count, err := ls.storage.GetStocksInfoByID(ctx, sku)
	if err != nil {
		return 0, logger.WithError(err, "getting stock item info failed")
	}
	logger.Debug("getting stocks info success", zap.Uint32("sku", sku), zap.Int("count", int(count)))
	return count, err
}
