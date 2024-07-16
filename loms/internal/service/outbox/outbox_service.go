package outbox

import (
	"context"
)

type Storage interface {
	FetchAndMark(ctx context.Context) error
}

type OutboxService struct {
	storage Storage
}

func NewOutboxService(storage Storage) *OutboxService {
	return &OutboxService{storage: storage}
}

func (os *OutboxService) Dispatch(ctx context.Context) error {
	err := os.storage.FetchAndMark(ctx)
	if err != nil {
		return err
	}

	return nil
}
