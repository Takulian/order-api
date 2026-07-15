package event

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	channel *amqp.Channel
}

func NewRabbitMQConsumer(conn *amqp.Connection) (*RabbitMQConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("gagal membuka chanel: %w", err)
	}

	err = ch.ExchangeDeclare(
		ExchangeOrderEvents,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("gagal deklarasi exhange: %w", err)
	}

	return &RabbitMQConsumer{
		channel: ch,
	}, nil
}

func (c *RabbitMQConsumer) consume(ctx context.Context, queueName, routingKey string, handler func(ctx context.Context, body []byte) error) error {
	q, err := c.channel.QueueDeclare(
		queueName,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("gagal deklarasi queue: %w", err)
	}

	err = c.channel.QueueBind(
		q.Name, routingKey, ExchangeOrderEvents, false, nil,
	)
	if err != nil {
		return fmt.Errorf("gagal bind queue %s: %w", queueName, err)
	}
	if err := c.channel.Qos(1, 0, false); err != nil {
		return fmt.Errorf("gagal set qos: %w", err)
	}

	msgs, err := c.channel.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("gagal mulai consume: %w", err)
	}

	for msg := range msgs {
		if err := handler(ctx, msg.Body); err != nil {
			msg.Nack(false, true)
			continue
		}
		msg.Ack(false)
	}
	return fmt.Errorf("channel pesan %s tertutup", queueName)
}

func (c *RabbitMQConsumer) ConsumeOrderCreated(ctx context.Context, handler func(ctx context.Context, evt OrderCreatedEvent) error) error {
	return c.consume(ctx, "order.created.queue", RoutingKeyOrderCreated, func(ctx context.Context, body []byte) error {
		var evt OrderCreatedEvent
		if err := json.Unmarshal(body, &evt); err != nil {
			return nil
		}
		return handler(ctx, evt)
	})
}

func (c *RabbitMQConsumer) ConsumeCheckout(ctx context.Context, handler func(ctx context.Context, evt CheckoutEvent) error) error {
	return c.consume(ctx, "order.checkout.queue", RoutingKeyCheckout, func(ctx context.Context, body []byte) error {
		var evt CheckoutEvent
		if err := json.Unmarshal(body, &evt); err != nil {
			return nil
		}
		return handler(ctx, evt)
	})
}

func (c *RabbitMQConsumer) Close() error {
	return c.channel.Close()
}
