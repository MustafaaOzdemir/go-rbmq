package models

type Order struct {
	ID        string  `json:"id"`
	UserEmail string  `json:"user_email"`
	UserPhone string  `json:"user_phone"`
	Amount    float64 `json:"amount"`
	ProductID string  `json:"product_id"`
}
