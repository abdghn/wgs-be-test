package model

type OrderItem struct {
	ID          int     `json:"id"`
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	OrderID     uint    `json:"-"`
}
