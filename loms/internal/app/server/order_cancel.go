package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256.ozon.ru/project/loms/internal/entity"
	servicepb "route256.ozon.ru/project/loms/pkg/api/loms/v1"
)

func (l *Server) OrderCancel(ctx context.Context, request *servicepb.OrderCancelRequest) (*emptypb.Empty, error) {
	if err := l.service.OrderCancel(ctx, entity.OrderID(request.OrderId)); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
