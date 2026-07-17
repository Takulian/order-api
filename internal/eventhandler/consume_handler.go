package eventhandler

import (
	"context"
	"log/slog"
	"order-api/internal/dto"
	"order-api/internal/event"
	"order-api/internal/observability"
	"order-api/internal/service"
)

type ConsumeHandler struct {
	consumer *event.RabbitMQConsumer
	service  *service.OrderService
	logger   *slog.Logger
}

func NewConsumeHandler(consumer *event.RabbitMQConsumer, service *service.OrderService, logger *slog.Logger) *ConsumeHandler {
	return &ConsumeHandler{
		consumer: consumer,
		service:  service,
		logger:   logger,
	}
}

func (h *ConsumeHandler) Start(ctx context.Context) {
	observability.SafeGo(ctx, h.logger, "consume-checkout", func() {
		err := h.consumer.ConsumeCheckout(ctx, func(ctx context.Context, evt event.CheckoutEvent) error {
			_, err := h.service.Create(ctx, dto.CreateOrderRequest{
				Customer: evt.Customer,
				Product:  evt.Product,
				Quantity: evt.Quantity,
			})
			return err
		})
		if err != nil {
			h.logger.Error("consumer order.checkout berhenti karena error", "error", err)
			panic(err)
		}
	})

	observability.SafeGo(ctx, h.logger, "consume-order-created", func() {
		err := h.consumer.ConsumeOrderCreated(ctx, func(ctx context.Context, evt event.OrderCreatedEvent) error {
			h.logger.InfoContext(ctx, "menerima order.created",
				"order_id", evt.OrderID,
				"consumer", evt.Customer,
				"product", evt.Product,
			)
			return nil
		})
		if err != nil {
			h.logger.Error("consumer order.created berhenti karena error", "error", err)
			panic(err)
		}
	})
}
