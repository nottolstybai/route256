package entity

import (
	"time"
)

type EventType string

const (
	EventOrderCreated  EventType = "order-created"
	EventOrderAwaiting EventType = "order-awaiting"
	EventOrderCanceled EventType = "order-canceled"
	EventOrderFailed   EventType = "order-failed"
	EventOrderPayed    EventType = "order-payed"
)

var orderStatusMapping = map[Status]EventType{
	StatusNew:             EventOrderCreated,
	StatusAwaitingPayment: EventOrderAwaiting,
	StatusFailed:          EventOrderFailed,
	StatusCancelled:       EventOrderCanceled,
	StatusPayed:           EventOrderPayed,
}

type Event struct {
	ID              int32     `json:"id"`
	OrderID         int32     `json:"order_id"`
	OrderStatus     EventType `json:"order_status"`
	OperationMoment time.Time `json:"moment"`
}

func NewEvent(id, orderID int32, orderStatus Status) Event {
	return Event{
		ID:              id,
		OrderID:         orderID,
		OrderStatus:     orderStatusMapping[orderStatus],
		OperationMoment: time.Now(),
	}
}
