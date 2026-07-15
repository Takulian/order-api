package event

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	channel *amqp.Channel
}

func NewRabbitMQPublisher(conn *amqp.Connection) (*RabbitMQPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("gagal membuka channel: %w", err)
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

	return &RabbitMQPublisher{
		channel: ch,
	}, nil
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, routingKey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal payload: %w", err)
	}

	return p.channel.PublishWithContext(
		ctx,
		ExchangeOrderEvents,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (p *RabbitMQPublisher) Close() error {
	return p.channel.Close()
}
