package app

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"route256.ozon.ru/project/notifier/internal/handler"
	"route256.ozon.ru/project/notifier/pkg/logger"
	"sync"
	"syscall"
)

type ConsumerGroup interface {
	StartConsume(ctx context.Context, wg *sync.WaitGroup)
	Close() error
	Errors() <-chan error
}

type App struct {
	cg ConsumerGroup
}

func NewApp(cfg Config) *App {
	kafka, err := handler.NewKafkaConsumerGroup(cfg.bootstrapServers, cfg.topics, cfg.groupId, nil)
	if err != nil {
		panic(err)
	}
	return &App{
		cg: kafka,
	}
}

func (a *App) Run() {
	wg := &sync.WaitGroup{}
	ctx := runSignalHandler(context.Background(), wg)

	defer a.cg.Close()
	runCGErrorHandler(ctx, a.cg, wg)
	a.cg.StartConsume(ctx, wg)
	wg.Wait()
}

func runSignalHandler(ctx context.Context, wg *sync.WaitGroup) context.Context {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	sigCtx, cancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		defer logger.Info("signal terminate")
		defer signal.Stop(sigterm)
		defer wg.Done()
		defer cancel()

		for {
			select {
			case sig, ok := <-sigterm:
				if !ok {
					logger.Error("signal chan closed: %s", zap.String("signal", sig.String()))
					return
				}

				logger.Info("signal received: %s", zap.String("signal", sig.String()))
				return
			case _, ok := <-sigCtx.Done():
				if !ok {
					logger.Error("signal context closed")
					return
				}

				logger.Info("signal ctx done: %s", zap.Error(ctx.Err()))
				return
			}
		}
	}()

	return sigCtx
}

func runCGErrorHandler(ctx context.Context, cg ConsumerGroup, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case chErr, ok := <-cg.Errors():
				if !ok {
					logger.Error("[cg-error] chan closed")
					return
				}

				logger.Error("[cg-error] error", zap.Error(chErr))
			case <-ctx.Done():
				logger.Info("[cg-error] ctx closed", zap.Error(ctx.Err()))
				return
			}
		}
	}()
}
