package event

import "time"

type OrderCreatedEvent struct {
	OrderID   int       `json:"order_id"`
	Customer  string    `json:"customer"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	RuotingKeyOrderCreated = "order.created"
	ExchangeOrderEvents    = "order.events"
)
