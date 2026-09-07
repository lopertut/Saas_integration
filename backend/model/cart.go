package model

type Cart struct {
	Id      int `json:"id"`
	User_id int `json:"user_id"`
}

type CartItem struct {
	Id        int      `json:"id"`
	CartId    int      `json:"cart_id"`
	ProductId int      `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Name      string   `json:"name"`
	Price     float64  `json:"price"`
	Images    []string `json:"images"`
	AvgRating float64  `json:"rating"`
}
