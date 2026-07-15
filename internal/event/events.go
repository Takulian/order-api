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
type CheckoutEvent struct {
	Customer string `json:"customer"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

const (
	RoutingKeyOrderCreated = "order.created"
	RoutingKeyCheckout     = "order.checkout"
	ExchangeOrderEvents    = "order.events"
)
