package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"route256.ozon.ru/project/loms/internal/app"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/pkg/logger"
	"route256.ozon.ru/project/loms/pkg/trace"
	"syscall"
	"time"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Init()
	defer logger.Sync()

	provider, err := trace.InitTracerProvider("http://jaeger:14268/api/traces", "loms")
	if err != nil {
		logger.Fatal("init tracer", zap.Error(err))
	}
	defer func() {
		provider.ForceFlush(context.Background())
		provider.Shutdown(context.Background())
	}()

	cfg := config.NewConfig()

	lomsApp := app.NewApp(cfg)
	lomsApp.Run(mainCtx)

	<-mainCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := lomsApp.Stop(ctx); err != nil {
		logger.Fatal("graceful shutdown failed", zap.Error(err))
	}
	logger.Info("graceful shutdown successful")
}
