package memory

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"math/rand"
	"reflect"
	"route256.ozon.ru/project/cart/internal/entity"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestMemoryRepo_Insert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		userID int64
		skuID  int64
		count  uint16
		want   map[int64][]entity.ProductInfo
	}{
		{
			name:   "Insert new item",
			userID: 1,
			skuID:  1001,
			count:  1,
			want: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 1},
				},
			},
		},
		{
			name:   "Insert item with existing userID and skuID",
			userID: 1,
			skuID:  1001,
			count:  2,
			want: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
				},
			},
		},
		{
			name:   "Insert an item with an existing userID but different skuID",
			userID: 1,
			skuID:  1002,
			count:  2,
			want: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
			},
		},
	}

	repo := NewMemoryRepo()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// dont use t.Parallel() since each tc expected result is based on previous tc
			ctx := context.Background()
			if err := repo.Insert(ctx, tc.userID, tc.skuID, tc.count); err != nil {
				t.Errorf("Error inserting item: %v", err)
			}

			if !reflect.DeepEqual(repo.storage, tc.want) {
				t.Errorf("Storage map does not match expected value")
			}
		})
	}
}

func TestMemoryRepo_Remove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  int64
		skuID   int64
		storage map[int64][]entity.ProductInfo
		want    map[int64][]entity.ProductInfo
	}{
		{
			name:   "Remove existing item",
			userID: 1,
			skuID:  1001,
			storage: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
			},
			want: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1002, Count: 2},
				},
			},
		},
		{
			name:   "Remove non-existing item",
			userID: 1,
			skuID:  1003,
			storage: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
			},
			want: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &MemoryRepo{
				storage: tc.storage,
			}

			if err := repo.Remove(context.Background(), tc.userID, tc.skuID); err != nil {
				t.Errorf("Error removing item: %v", err)
			}

			require.Equal(t, tc.want, repo.storage)
		})
	}
}

func TestMemoryRepo_RemoveByUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  int64
		storage map[int64][]entity.ProductInfo
		want    map[int64][]entity.ProductInfo
	}{
		{
			name:   "Remove existing user's items",
			userID: 1,
			storage: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
				2: {
					{SkuID: 2001, Count: 1},
				},
			},
			want: map[int64][]entity.ProductInfo{
				2: {
					{SkuID: 2001, Count: 1},
				},
			},
		},
		{
			name:   "Remove non-existing user's items",
			userID: 3,
			storage: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
				2: {
					{SkuID: 2001, Count: 1},
				},
			},
			want: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
				2: {
					{SkuID: 2001, Count: 1},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &MemoryRepo{
				storage: tc.storage,
			}

			if err := repo.RemoveByUserID(context.Background(), tc.userID); err != nil {
				t.Errorf("Error removing items: %v", err)
			}
			require.Equal(t, tc.want, repo.storage)
		})
	}
}

func TestMemoryRepo_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  int64
		storage map[int64][]entity.ProductInfo
		want    []entity.ProductInfo
	}{
		{
			name:   "List existing user's items",
			userID: 1,
			storage: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
				2: {
					{SkuID: 2001, Count: 1},
				},
			},
			want: []entity.ProductInfo{
				{SkuID: 1001, Count: 3},
				{SkuID: 1002, Count: 2},
			},
		},
		{
			name:   "List non-existing user's items",
			userID: 3,
			storage: map[int64][]entity.ProductInfo{
				1: {
					{SkuID: 1001, Count: 3},
					{SkuID: 1002, Count: 2},
				},
			},
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &MemoryRepo{
				storage: tc.storage,
			}

			items, err := repo.List(context.Background(), tc.userID)
			if err != nil {
				t.Errorf("Error listing items: %v", err)
			}
			require.Equal(t, tc.want, items)
		})
	}
}

func BenchmarkMemoryRepo_Insert_Complex(b *testing.B) {
	repo := NewMemoryRepo()
	rand.New(rand.NewSource(time.Now().Unix()))

	// Simulate multiple users inserting items
	userIDs := generateRandomInt64s(b.N)
	skuIDs := generateRandomInt64s(b.N)
	counts := generateRandomUint16s(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := repo.Insert(context.Background(), userIDs[i], skuIDs[i], counts[i])
		if err != nil {
			b.Errorf("Error inserting item: %v", err)
		}
	}
}

func generateRandomInt64s(n int) []int64 {
	result := make([]int64, n)
	for i := 0; i < n; i++ {
		result[i] = rand.Int63()
	}
	return result
}

func generateRandomUint16s(n int) []uint16 {
	result := make([]uint16, n)
	for i := 0; i < n; i++ {
		result[i] = uint16(rand.Intn(65536))
	}
	return result
}
