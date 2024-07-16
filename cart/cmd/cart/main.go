package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"route256.ozon.ru/project/cart/internal/app/httpserver"
	"route256.ozon.ru/project/cart/internal/client/loms_service"
	"route256.ozon.ru/project/cart/internal/client/product_service"
	"route256.ozon.ru/project/cart/internal/config"
	"route256.ozon.ru/project/cart/internal/repository/memory"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/pkg/logger"
	"route256.ozon.ru/project/cart/pkg/trace"
	"syscall"
	"time"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Init()
	defer logger.Sync()

	provider, err := trace.InitTracerProvider("http://jaeger:14268/api/traces", "cart")
	if err != nil {
		logger.Fatal("init tracer", zap.Error(err))
	}
	defer func() {
		provider.ForceFlush(context.Background())
		provider.Shutdown(context.Background())
	}()

	cfg := config.NewConfig()

	repo := memory.NewMemoryRepo()
	productService := product_service.NewProductServiceClient(cfg.ProductServiceHost, cfg.ProductServiceToken)
	lomsService := loms_service.NewLoms(cfg.LomsServiceHost)
	cartService := service.NewCartService(productService, repo, lomsService)

	cartServer := httpserver.NewServer(cartService, cfg.ServeAddr)
	go func() {
		if err := cartServer.Serve(); err != nil {
			logger.Error("serve error", zap.Error(err))
		}
	}()

	<-mainCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := cartServer.Shutdown(ctx); err != nil {
		logger.Fatal("graceful shutdown failed", zap.Error(err))
	}

	logger.Info("graceful shutdown successful")
}
