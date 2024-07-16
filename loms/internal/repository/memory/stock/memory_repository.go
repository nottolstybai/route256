package stock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"route256.ozon.ru/project/loms/internal/entity"
	"sync"
)

var (
	ErrReserveNotFound        = errors.New("couldn't find reservation with specified sku")
	ErrReserveTooLarge        = errors.New("reservation amount is larger thant total amount")
	ErrCancelOrPayNumTooLarge = errors.New("number of items to pay/cancel is larger than available in stock")
)

type reservation struct {
	totalCount    uint32
	reservedCount uint32
}

type MemoryStocksRepository struct {
	storage map[uint32]reservation
	mutex   sync.RWMutex
}

func NewMemoryStocksRepository() *MemoryStocksRepository {
	storage := make(map[uint32]reservation)
	return &MemoryStocksRepository{
		storage: storage,
		mutex:   sync.RWMutex{},
	}
}

func NewStocksRepositoryFromJSON() *MemoryStocksRepository {
	// bullshit code, but it's ok since we will use postgres as repo
	const defaultResourcesPath = ""

	resourcesPath := os.Getenv("RESOURCES_PATH")
	if resourcesPath == "" {
		resourcesPath = defaultResourcesPath
	}

	type reserve struct {
		Sku           uint32 `json:"sku"`
		TotalCount    uint32 `json:"total_count"`
		ReservedCount uint32 `json:"reserved"`
	}

	stocksRepo := MemoryStocksRepository{
		storage: make(map[uint32]reservation),
		mutex:   sync.RWMutex{},
	}

	relativePath := filepath.Join(resourcesPath, "stock-data.json")
	file, err := os.Open(relativePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var reserves []reserve
	if err := json.NewDecoder(file).Decode(&reserves); err != nil {
		fmt.Println(err)
	}

	for _, r := range reserves {
		stocksRepo.storage[r.Sku] = reservation{totalCount: r.TotalCount, reservedCount: r.ReservedCount}
	}

	return &stocksRepo
}

func (m *MemoryStocksRepository) Reserve(_ context.Context, items []entity.Item) error {
	for _, item := range items {
		if _, ok := m.storage[item.SKU]; !ok {
			return ErrReserveNotFound
		}

		reserve := m.storage[item.SKU]
		if item.Count > reserve.totalCount-reserve.reservedCount {
			return ErrReserveTooLarge
		}
		reserve.reservedCount += item.Count
		m.storage[item.SKU] = reserve
	}
	return nil
}

func (m *MemoryStocksRepository) RemoveReservation(_ context.Context, items []entity.Item) error {
	for _, item := range items {
		if _, ok := m.storage[item.SKU]; !ok {
			return ErrReserveNotFound
		}

		reserve := m.storage[item.SKU]
		if item.Count > reserve.reservedCount || item.Count > reserve.totalCount {
			return ErrCancelOrPayNumTooLarge
		}

		reserve.reservedCount -= item.Count
		reserve.totalCount -= item.Count
		m.storage[item.SKU] = reserve
	}
	return nil
}

func (m *MemoryStocksRepository) CancelReservation(_ context.Context, items []entity.Item) error {
	for _, item := range items {
		if _, ok := m.storage[item.SKU]; !ok {
			return ErrReserveNotFound
		}

		reserve := m.storage[item.SKU]
		if item.Count > reserve.reservedCount || item.Count > reserve.totalCount {
			return ErrCancelOrPayNumTooLarge
		}

		reserve.reservedCount -= item.Count
		m.storage[item.SKU] = reserve
	}
	return nil
}

func (m *MemoryStocksRepository) GetBySku(_ context.Context, sku uint32) (entity.Count, error) {
	if _, ok := m.storage[sku]; !ok {
		return 0, ErrReserveNotFound
	}
	return entity.Count(m.storage[sku].totalCount - m.storage[sku].reservedCount), nil
}
