-- +goose Up
-- +goose StatementBegin
create table if not exists "order"
(
    order_id int not null,
    user_id  int not null,
    status   int not null,
    primary key (order_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "order" cascade;
-- +goose StatementEnd
