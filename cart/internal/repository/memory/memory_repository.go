package memory

import (
	"context"
	"go.opentelemetry.io/otel"
	"route256.ozon.ru/project/cart/internal/entity"
	"sync"
)

type MemoryRepo struct {
	storage map[int64][]entity.ProductInfo
	mutex   sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	storage := make(map[int64][]entity.ProductInfo, 10)
	return &MemoryRepo{storage: storage, mutex: sync.RWMutex{}}
}

func (m *MemoryRepo) Insert(ctx context.Context, userID, skuID int64, count uint16) error {
	ctx, span := otel.Tracer("default").Start(ctx, "MemoryRepo.Insert")
	defer span.End()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.storage[userID]
	if !ok {
		items := make([]entity.ProductInfo, 0)
		items = append(items, entity.ProductInfo{
			SkuID: skuID,
			Count: count,
		})
		m.storage[userID] = items
	} else {
		for i := 0; i < len(m.storage[userID]); i++ {
			if m.storage[userID][i].SkuID == skuID {
				m.storage[userID][i].Count += count
				return nil
			}
		}
		m.storage[userID] = append(m.storage[userID], entity.ProductInfo{SkuID: skuID, Count: count})
	}
	return nil
}

func (m *MemoryRepo) Remove(ctx context.Context, userID, skuID int64) error {
	ctx, span := otel.Tracer("default").Start(ctx, "MemoryRepo.Remove")
	defer span.End()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	items, ok := m.storage[userID]
	if ok {
		for i := 0; i < len(m.storage[userID]); i++ {
			if m.storage[userID][i].SkuID == skuID {
				m.storage[userID][i] = m.storage[userID][len(items)-1]
				m.storage[userID] = m.storage[userID][:len(items)-1]
				return nil
			}
		}
	}
	return nil
}

func (m *MemoryRepo) RemoveByUserID(ctx context.Context, userID int64) error {
	ctx, span := otel.Tracer("default").Start(ctx, "MemoryRepo.RemoveByUserID")
	defer span.End()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.storage[userID]
	if ok {
		delete(m.storage, userID)
	}
	return nil
}

func (m *MemoryRepo) List(ctx context.Context, userID int64) ([]entity.ProductInfo, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "MemoryRepo.RemoveByUserID")
	defer span.End()

	items, ok := m.storage[userID]
	if ok {
		return items, nil
	}
	return nil, nil
}
