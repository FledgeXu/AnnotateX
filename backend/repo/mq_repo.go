package repo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type IMqRepo interface {
	DeclareExchange(name string, kind string, durable bool) error
	Publish(exchange, routingKey string, body any) error
}

type MqRepo struct {
	Conn *amqp091.Connection
}

func NewMqRepo(client *amqp091.Connection) *MqRepo {
	return &MqRepo{
		Conn: client,
	}
}

func (c *MqRepo) DeclareExchange(name, kind string, durable bool) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel failed: %w", err)
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		name,    // exchange name
		kind,    // "direct", "fanout", "topic", "headers"
		durable, // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
}

func (r *MqRepo) Publish(exchange, routingKey string, msg any) error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel failed: %w", err)
	}
	defer ch.Close()

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}

	return ch.Publish(
		exchange, routingKey, false, false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
			Timestamp:    time.Now(),
		})
}
