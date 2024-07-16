package order

import (
	"context"
	"errors"
	"math/rand"
	"route256.ozon.ru/project/loms/internal/entity"
	"sync"
	"time"
)

var ErrOrderNotFound = errors.New("couldn't find order with specified orderID")

type (
	storageItem struct {
		sku   uint32
		count uint32
	}
	order struct {
		user   int64
		status entity.Status
		items  []storageItem
	}
)

type MemoryOrderRepository struct {
	storage map[int64]order
	mutex   sync.RWMutex
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	storage := make(map[int64]order)
	return &MemoryOrderRepository{
		storage: storage,
		mutex:   sync.RWMutex{},
	}
}

func (m *MemoryOrderRepository) Create(_ context.Context, userID int64, items []entity.Item) (entity.OrderID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	status := entity.StatusNew
	orderID := genOrderID()

	storageItems := make([]storageItem, 0, len(items))
	for _, item := range items {
		storageItems = append(storageItems, storageItem{
			sku:   item.SKU,
			count: item.Count,
		})
	}

	m.storage[orderID] = order{
		user:   userID,
		status: status,
		items:  storageItems,
	}
	return entity.OrderID(orderID), nil
}

func (m *MemoryOrderRepository) SetStatus(_ context.Context, orderID entity.OrderID, status entity.Status) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if orderInfo, ok := m.storage[int64(orderID)]; ok {
		orderInfo.status = status
		m.storage[int64(orderID)] = orderInfo
		return nil
	}
	return ErrOrderNotFound
}

func (m *MemoryOrderRepository) GetByOrderID(_ context.Context, orderID entity.OrderID) (*entity.OrderInfo, error) {
	if order, ok := m.storage[int64(orderID)]; ok {
		items := make([]entity.Item, 0, len(order.items))

		for _, storageItem := range order.items {
			items = append(items, entity.Item{
				SKU:   storageItem.sku,
				Count: storageItem.count,
			})
		}

		orderInfo := entity.OrderInfo{
			Status: order.status,
			User:   order.user,
			Items:  items,
		}
		return &orderInfo, nil
	}

	return nil, ErrOrderNotFound
}

func genOrderID() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63()
}
