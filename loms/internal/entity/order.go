package entity

type OrderID int64

type Count uint32

type Item struct {
	SKU   uint32
	Count uint32
}

type OrderInfo struct {
	Status Status
	User   int64
	Items  []Item
}
