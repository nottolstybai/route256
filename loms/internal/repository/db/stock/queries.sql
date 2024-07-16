-- name: ReserveItems :exec
update stock
set reserved_count=reserved_count+$1
where sku=$2;

-- name: RemoveReservationOfItem :exec
update stock
set reserved_count=reserved_count-$1,
    total_count=total_count-$1
where sku=$2;

-- name: CancelReservationOfItem :exec
update stock
set reserved_count=reserved_count-$1
where sku=$2;

-- name: GetBySku :one
select * from stock
where sku=$1;

-- name: AddStocks :copyfrom
insert into stock (sku, total_count, reserved_count)
values ($1, $2, $3);