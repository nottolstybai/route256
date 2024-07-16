package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"math/rand"
	"os"
	"path/filepath"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/internal/repository/db/order"
	"route256.ozon.ru/project/loms/internal/repository/db/outbox"
	"route256.ozon.ru/project/loms/internal/repository/db/stock"
	"time"
)

var (
	ErrReservation   = errors.New("error while creating reservation")
	ErrOrderNotFound = errors.New("order with specified orderID not found")
	ErrStockNotFound = errors.New("stock item with specified sku not found")
)

type messageSender interface {
	SendMessage(event entity.Event) error
}

type LomsStorage struct {
	pool       *ConnectionPool
	orderRepo  order.Queries
	stockRepo  stock.Queries
	outboxRepo outbox.Queries
	sender     messageSender
}

func NewLomsStorage(masterPool, replicaPool *pgxpool.Pool, sender messageSender) *LomsStorage {
	return &LomsStorage{
		pool:   NewConnectionPool(masterPool, replicaPool),
		sender: sender,
	}
}

func (ls *LomsStorage) CloseConnections() error {
	ls.pool.Close()

	return nil
}

func (ls *LomsStorage) CreateOrder(ctx context.Context, userID int64, items []entity.Item) (entity.OrderID, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsStorage.CreateOrder")
	defer span.End()

	tx, err := ls.pool.Master().Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := tx.Commit(ctx)
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	orderID := genOrderID()
	if err = ls.createOrderHelper(ctx, orderID, userID, items); err != nil {
		return 0, err
	}

	start := time.Now()
	if err = ls.outboxRepo.WithTx(tx).CreateEvent(ctx, outbox.CreateEventParams{
		Orderid:     orderID,
		OrderStatus: int32(entity.StatusNew),
	}); err != nil {
		return entity.OrderID(orderID), fmt.Errorf("failed to create event with status=%d for orderID=%d: %v", entity.StatusNew, orderID, err)
	}
	measureMetrics(insert, start, err)

	reserveErr := ls.tryReserveItems(ctx, items)
	if reserveErr != nil {
		if err = ls.setStatusAndEvent(ctx, tx, orderID, entity.StatusFailed); err != nil {
			return 0, err
		}
		return 0, ErrReservation
	}

	if err = ls.setStatusAndEvent(ctx, tx, orderID, entity.StatusAwaitingPayment); err != nil {
		return 0, err
	}
	return entity.OrderID(orderID), nil
}

func (ls *LomsStorage) GetOrderInfoByID(ctx context.Context, orderID entity.OrderID) (*entity.OrderInfo, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsStorage.GetOrderInfoByID")
	defer span.End()

	start := time.Now()
	rows, err := order.New(ls.pool.Acquire()).GetByOrderID(ctx, int32(orderID))
	if err != nil {
		return nil, err
	}
	measureMetrics(find, start, err)

	orderInfo := toEntityOrderInfo(rows)
	if orderInfo == nil {
		return nil, ErrOrderNotFound
	}
	return orderInfo, nil
}

func (ls *LomsStorage) PayOrder(ctx context.Context, orderID entity.OrderID) error {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsStorage.PayOrder")
	defer span.End()

	tx, err := ls.pool.Master().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	start := time.Now()
	rows, err := ls.orderRepo.WithTx(tx).GetByOrderID(ctx, int32(orderID))
	if err != nil {
		return ErrOrderNotFound
	}
	measureMetrics(find, start, err)

	for _, row := range rows {
		start = time.Now()
		if err = ls.stockRepo.WithTx(tx).RemoveReservationOfItem(ctx, stock.RemoveReservationOfItemParams{
			ReservedCount: row.Count,
			Sku:           row.SkuID,
		}); err != nil {
			return err
		}
		measureMetrics(update, start, err)
	}

	if err := ls.setStatusAndEvent(ctx, tx, int32(orderID), entity.StatusPayed); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (ls *LomsStorage) CancelOrder(ctx context.Context, orderID entity.OrderID) error {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsStorage.CancelOrder")
	defer span.End()

	tx, err := ls.pool.Master().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	start := time.Now()
	rows, err := ls.orderRepo.WithTx(tx).GetByOrderID(ctx, int32(orderID))
	if err != nil {
		return err
	}
	measureMetrics(find, start, err)

	for _, row := range rows {
		start = time.Now()
		if err = ls.stockRepo.WithTx(tx).CancelReservationOfItem(ctx, stock.CancelReservationOfItemParams{
			ReservedCount: row.Count,
			Sku:           row.SkuID,
		}); err != nil {
			return err
		}
		measureMetrics(update, start, err)
	}

	if err := ls.setStatusAndEvent(ctx, tx, int32(orderID), entity.StatusCancelled); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (ls *LomsStorage) GetStocksInfoByID(ctx context.Context, sku uint32) (entity.Count, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsStorage.GetStocksInfoByID")
	defer span.End()

	start := time.Now()
	stockInfo, err := stock.New(ls.pool.Acquire()).GetBySku(ctx, int32(sku))
	if err != nil {
		return 0, ErrStockNotFound
	}
	measureMetrics(find, start, err)

	return entity.Count(stockInfo.TotalCount - stockInfo.ReservedCount), nil
}

func (ls *LomsStorage) FetchAndMark(ctx context.Context) error {
	ctx, span := otel.Tracer("default").Start(ctx, "LomsStorage.FetchAndMark")
	defer span.End()

	tx, err := ls.pool.Master().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	start := time.Now()
	event, err := ls.outboxRepo.WithTx(tx).GetNextEvent(ctx)
	if err != nil {
		return fmt.Errorf("failed to get event: %v", err)
	}
	measureMetrics(find, start, err)

	eventMsg := entity.NewEvent(event.ID, event.Orderid, entity.Status(event.OrderStatus))

	if err = ls.sender.SendMessage(eventMsg); err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	start = time.Now()
	if err = ls.outboxRepo.WithTx(tx).MarkEventAsSent(ctx, event.ID); err != nil {
		return fmt.Errorf("failed to mark event as 'sent': %v", err)
	}
	measureMetrics(update, start, err)

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func genOrderID() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int31()
}

func NewLomsStorageFromJson(masterPool, replicaPool *pgxpool.Pool, sender messageSender) *LomsStorage {
	const defaultResourcesPath = "C:\\Users\\HP\\GolandProjects\\student-project\\loms\\resources"

	resourcesPath := os.Getenv("RESOURCES_PATH")
	if resourcesPath == "" {
		resourcesPath = defaultResourcesPath
	}

	type stockItem struct {
		Sku           int32 `json:"sku"`
		TotalCount    int32 `json:"total_count"`
		ReservedCount int32 `json:"reserved"`
	}

	relativePath := filepath.Join(resourcesPath, "stock-data.json")
	file, err := os.Open(relativePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var stocks []stockItem
	if err := json.NewDecoder(file).Decode(&stocks); err != nil {
		panic(err.Error())
	}

	args := make([]stock.AddStocksParams, 0, len(stocks))
	for _, s := range stocks {
		args = append(args, stock.AddStocksParams{Sku: s.Sku, TotalCount: s.TotalCount, ReservedCount: s.ReservedCount})
	}

	_, err = stock.New(masterPool).AddStocks(context.Background(), args)
	if err != nil {
		panic(err.Error())
	}

	return &LomsStorage{
		pool:   NewConnectionPool(masterPool, replicaPool),
		sender: sender,
	}
}
