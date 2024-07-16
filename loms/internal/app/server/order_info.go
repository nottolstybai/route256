package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/entity"
	servicepb "route256.ozon.ru/project/loms/pkg/api/loms/v1"
)

func (l *Server) OrderInfo(ctx context.Context, request *servicepb.OrderInfoRequest) (*servicepb.OrderInfoResponse, error) {
	info, err := l.service.OrderInfo(ctx, entity.OrderID(request.OrderId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &servicepb.OrderInfoResponse{
		Status: servicepb.Statuses(info.Status),
		User:   info.User,
	}
	for _, item := range info.Items {
		resp.Items = append(resp.Items, &servicepb.Item{
			Sku:   item.SKU,
			Count: item.Count,
		})
	}

	return resp, nil
}
