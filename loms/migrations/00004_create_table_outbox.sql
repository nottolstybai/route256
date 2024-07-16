-- +goose Up
-- +goose StatementBegin
create type event_status as enum ('new', 'sent');

create table if not exists outbox
(
    id            serial       not null,
    orderID       int          not null,
    order_status  int          not null,
    event_status  event_status not null,
    dttm_inserted timestamp default now(),
    primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists outbox;
-- +goose StatementEnd
