package loms_service

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	_ "google.golang.org/grpc/metadata"
	"route256.ozon.ru/project/cart/internal/entity"
	"route256.ozon.ru/project/cart/pkg/api/loms/v1"
)

type LomsGRPCClient struct {
	client loms.LomsClient
}

func NewLoms(address string) *LomsGRPCClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := loms.NewLomsClient(conn)
	return &LomsGRPCClient{client: client}
}

func (lc *LomsGRPCClient) OrderCreate(ctx context.Context, userID int64, itemList *entity.ListItems) (int64, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsGRPCClient.OrderCreate")
	defer span.End()

	mdCtx := Trace(ctx)

	reqItems := make([]*loms.Item, 0, len(itemList.Items))
	for _, item := range itemList.Items {
		reqItems = append(reqItems, &loms.Item{
			Sku:   uint32(item.SkuID),
			Count: uint32(item.Count),
		})
	}

	resp, err := lc.client.OrderCreate(mdCtx, &loms.OrderCreateRequest{User: userID, Items: reqItems})
	if err != nil {
		return 0, err
	}

	return resp.OrderId, nil
}

func (lc *LomsGRPCClient) StocksInfo(ctx context.Context, sku uint32) (uint16, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsGRPCClient.StocksInfo", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	mdCtx := Trace(ctx)

	resp, err := lc.client.StocksInfo(mdCtx, &loms.StocksInfoRequest{Sku: sku})
	if err != nil {
		return 0, err
	}
	return uint16(resp.Count), nil
}

func Trace(ctx context.Context) context.Context {
	spanCtx := trace.SpanContextFromContext(ctx)

	md := metadata.New(map[string]string{
		"x-trace-id": spanCtx.TraceID().String(),
		"x-span-id":  spanCtx.SpanID().String(),
	})
	mdCtx := metadata.NewOutgoingContext(ctx, md)

	return mdCtx
}
