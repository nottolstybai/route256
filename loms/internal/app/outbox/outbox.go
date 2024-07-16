package outbox

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"route256.ozon.ru/project/loms/pkg/logger"
	"time"
)

type Dispatcher interface {
	Dispatch(ctx context.Context) error
}

type OutboxSender struct {
	dispatcher Dispatcher
}

func NewOutboxSender(dispatcher Dispatcher) *OutboxSender {
	return &OutboxSender{dispatcher: dispatcher}
}

func (o *OutboxSender) RunDispatcher(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second):
			if err := o.dispatcher.Dispatch(ctx); err != nil {
				if !errors.As(err, &pgx.ErrNoRows) {
					logger.Error("error while running outbox loop", zap.Error(err))
				}
			}
		}
	}
}
