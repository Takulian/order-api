package dto

type UpdateOrderRequest struct {
	Customer string `json:"customer"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}