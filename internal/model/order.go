package model

var NextID = 3

type Order struct {
	ID       int    `json:"id"`
	Customer string `json:"customer"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

var Orders = []Order{
	{
		ID:       1,
		Customer: "Andi",
		Product:  "Laptop",
		Quantity: 1,
		Status:   "Pending",
	},
	{
		ID:       2,
		Customer: "Budi",
		Product:  "Mouse",
		Quantity: 2,
		Status:   "Paid",
	},
}