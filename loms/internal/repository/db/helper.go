package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/internal/repository/db/order"
	"route256.ozon.ru/project/loms/internal/repository/db/outbox"
	"time"
)

func (ls *LomsStorage) createOrderHelper(ctx context.Context, orderID int32, userID int64, items []entity.Item) error {
	tx, err := ls.pool.Master().Begin(ctx) // create savepoint
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	start := time.Now()
	if err := ls.orderRepo.WithTx(tx).AddOrder(ctx, order.AddOrderParams{
		OrderID: orderID,
		UserID:  int32(userID),
		Status:  int32(entity.StatusNew),
	}); err != nil {
		return err
	}
	measureMetrics(insert, start, err)

	start = time.Now()
	if _, err := ls.orderRepo.WithTx(tx).AddOrderStock(ctx, toAddOrderStock(orderID, items)); err != nil {
		return err
	}
	measureMetrics(insert, start, err)

	return tx.Commit(ctx) // release savepoint
}

func (ls *LomsStorage) setStatusAndEvent(ctx context.Context, tx pgx.Tx, orderID int32, status entity.Status) (err error) {
	start := time.Now()
	if err = ls.orderRepo.WithTx(tx).SetStatusByOrderID(ctx, order.SetStatusByOrderIDParams{
		OrderID: orderID,
		Status:  int32(status),
	}); err != nil {
		return fmt.Errorf("failed to set status with status=%d for orderID=%d: %v", status, orderID, err)
	}
	measureMetrics(update, start, err)

	start = time.Now()
	if err = ls.outboxRepo.WithTx(tx).CreateEvent(ctx, outbox.CreateEventParams{
		Orderid:     orderID,
		OrderStatus: int32(status),
	}); err != nil {
		return fmt.Errorf("failed to create event with status=%d for orderID=%d: %v", status, orderID, err)
	}
	measureMetrics(insert, start, err)

	return nil
}

func (ls *LomsStorage) tryReserveItems(ctx context.Context, items []entity.Item) error {
	tx, err := ls.pool.Master().Begin(ctx) // create savepoint
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		start := time.Now()
		if err := ls.stockRepo.WithTx(tx).ReserveItems(ctx, toReserveItems(item)); err != nil {
			return ErrReservation
		}
		measureMetrics(update, start, err)
	}
	return tx.Commit(ctx) // release savepoint
}
