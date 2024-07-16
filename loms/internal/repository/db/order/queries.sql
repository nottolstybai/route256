-- name: AddOrder :exec
INSERT INTO "order" (order_id, user_id, status)
VALUES ($1, $2, $3);

-- name: AddOrderStock :copyfrom
insert into order_stock (order_id, sku_id, count)
values ($1, $2, $3);

-- name: SetStatusByOrderID :exec
update "order"
set status=$1
where order_id=$2;

-- name: GetByOrderID :many
select status, user_id, sku_id, count from "order"
join order_stock os on "order".order_id = os.order_id
where os.order_id=$1;