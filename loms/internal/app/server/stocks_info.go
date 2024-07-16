package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	servicepb "route256.ozon.ru/project/loms/pkg/api/loms/v1"
)

func (l *Server) StocksInfo(ctx context.Context, request *servicepb.StocksInfoRequest) (*servicepb.StocksInfoResponse, error) {
	count, err := l.service.StocksInfo(ctx, request.Sku)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &servicepb.StocksInfoResponse{Count: uint32(count)}, nil
}
