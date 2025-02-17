// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: copyfrom.go

package stock

import (
	"context"
)

// iteratorForAddStocks implements pgx.CopyFromSource.
type iteratorForAddStocks struct {
	rows                 []AddStocksParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddStocks) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddStocks) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Sku,
		r.rows[0].TotalCount,
		r.rows[0].ReservedCount,
	}, nil
}

func (r iteratorForAddStocks) Err() error {
	return nil
}

func (q *Queries) AddStocks(ctx context.Context, arg []AddStocksParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"stock"}, []string{"sku", "total_count", "reserved_count"}, &iteratorForAddStocks{rows: arg})
}
