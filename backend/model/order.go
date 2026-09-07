package model

import "time"

type Order struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	Country    string    `json:"country"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	FirstName  string    `json:"firstname"`
	LastName   string    `json:"lastname"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	TotalPrice float64   `json:"total_price"`
}

type OrderItem struct {
	Id          int     `json:"id"`
	OrderId     int     `json:"order_id"`
	ProductId   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
