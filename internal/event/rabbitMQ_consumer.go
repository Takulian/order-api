package event

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedHandler func(ctx context.Context, evt OrderCreatedEvent) error

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

func (c *RabbitMQConsumer) ConsumeOrderCreated(ctx context.Context, handler OrderCreatedHandler) error {
	queueName := "order.created.queue"

	q, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("gagal deklarasi queue: %w", err)
	}

	err = c.channel.QueueBind(
		q.Name,
		RuotingKeyOrderCreated,
		ExchangeOrderEvents,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("gagal bind ke queue exchange: %w", err)
	}

	if err := c.channel.Qos(1, 0, false); err != nil {
		return fmt.Errorf("gagal set qos: %w", err)
	}

	msgs, err := c.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("gagal mulai consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel pesan tertutup")
			}
			var evt OrderCreatedEvent
			if err := json.Unmarshal(msg.Body, &evt); err != nil {
				msg.Nack(false, false)
				continue
			}
			if err := handler(ctx, evt); err != nil {
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}
}

func (c *RabbitMQConsumer) Close() error {
	return c.channel.Close()
}
