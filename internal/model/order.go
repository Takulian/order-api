package model

type Order struct {
	ID       int    `json:"id"`
	Customer string `json:"customer"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}