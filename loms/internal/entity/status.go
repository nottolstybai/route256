package entity

type Status int

const (
	StatusNew Status = iota
	StatusAwaitingPayment
	StatusFailed
	StatusPayed
	StatusCancelled
)
