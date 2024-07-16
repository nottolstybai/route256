-- +goose Up
-- +goose StatementBegin
insert into stock(sku, total_count, reserved_count)
values (1, 100, 10),
       (2, 100, 20),
       (3, 100, 30);


insert into "order" (order_id, user_id, status)
values (1, 123, 1),
       (2, 123, 1),
       (3, 123, 1);

insert into order_stock(order_id, sku_id, count)
values (1, 1, 10),
       (2, 2, 20),
       (3, 3, 30);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from order_stock;
delete from "order";
delete from stock
-- +goose StatementEnd
