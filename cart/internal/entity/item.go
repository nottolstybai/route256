package entity

type Item struct {
	SkuID int64
	Name  string
	Count uint16
	Price uint32
}

type ListItems struct {
	Items      []Item
	TotalPrice uint32
}
