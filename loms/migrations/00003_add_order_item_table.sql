-- +goose Up
-- +goose StatementBegin
create table if not exists order_stock
(
    order_id int not null references "order"(order_id),
    sku_id   int not null references stock(sku),
    count    int not null,
    primary key (order_id, sku_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists order_stock cascade;
-- +goose StatementEnd
