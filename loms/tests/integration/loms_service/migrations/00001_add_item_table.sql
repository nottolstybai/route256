-- +goose Up
-- +goose StatementBegin
create table if not exists stock
(
    sku            int not null,
    total_count    int not null check (total_count >= 0),
    reserved_count int not null check (reserved_count >= 0 and reserved_count <= total_count),
    primary key (sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists stock cascade;
-- +goose StatementEnd
