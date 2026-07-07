package dto

type CreateOrderRequest struct {
	Customer string `json:"customer"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}