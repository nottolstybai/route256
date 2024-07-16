package server

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/internal/repository/db"
	servicepb "route256.ozon.ru/project/loms/pkg/api/loms/v1"
)

func (l *Server) OrderCreate(ctx context.Context, request *servicepb.OrderCreateRequest) (*servicepb.OrderCreateResponse, error) {
	var items []entity.Item

	for _, item := range request.Items {
		items = append(items, entity.Item{
			SKU:   item.Sku,
			Count: item.Count,
		})
	}

	orderID, err := l.service.OrderCreate(ctx, request.User, items)
	if err != nil {
		if errors.Is(err, db.ErrReservation) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &servicepb.OrderCreateResponse{OrderId: int64(orderID)}, nil
}
