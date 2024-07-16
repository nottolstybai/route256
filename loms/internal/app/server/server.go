package server

import (
	"context"
	"route256.ozon.ru/project/loms/internal/entity"
	servicepb "route256.ozon.ru/project/loms/pkg/api/loms/v1"
)

type Service interface {
	OrderCreate(ctx context.Context, userID int64, items []entity.Item) (entity.OrderID, error)
	OrderInfo(ctx context.Context, orderID entity.OrderID) (*entity.OrderInfo, error)
	OrderPay(ctx context.Context, orderID entity.OrderID) error
	OrderCancel(ctx context.Context, orderID entity.OrderID) error
	StocksInfo(ctx context.Context, sku uint32) (entity.Count, error)
}

type Server struct {
	servicepb.UnimplementedLomsServer
	service Service
}

func NewServer(service Service) *Server {
	return &Server{service: service}
}
