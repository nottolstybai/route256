package mw

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
	"route256.ozon.ru/project/loms/pkg/logger"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message)) // для превращения protbuf структур в json используем google.golang.org/protobuf/encoding/protojson пакет а не encoding/json

	logger.Debug(
		"got grpc request",
		zap.String("method", info.FullMethod),
		zap.String("body", string(raw)),
	)

	return handler(ctx, req)
}

func WithHTTPLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			return
		}

		logger.Debug(
			"got http request",
			zap.String("method", r.Method),
			zap.String("PATH", r.URL.Path),
		)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
