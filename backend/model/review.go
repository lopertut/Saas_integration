package model

import "time"

type Review struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	UserId    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
