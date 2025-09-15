package model

type Order struct {
	ID   string
	User string
}

type Invoice struct {
	ID      string
	OrderID string
	Amount  int
}
