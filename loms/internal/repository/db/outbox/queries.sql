-- name: CreateEvent :exec
insert into outbox (orderID, order_status, event_status)
values ($1, $2, 'new');

-- name: MarkEventAsSent :exec
update outbox
set event_status='sent'
where id=$1;

-- name: GetNextEvent :one
select * from outbox
where event_status='new'
order by dttm_inserted
limit 1;
