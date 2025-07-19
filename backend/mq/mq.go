package mq

import (
	"annotate-x/models"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

var (
	conn    *amqp091.Connection
	once    sync.Once
	initErr error
)

func InitMQ(mqURL models.MQUrl) *amqp091.Connection {
	once.Do(func() {
		conn, initErr = amqp091.Dial(mqURL)
		if initErr != nil {
			panic(fmt.Sprintf("connect mq failed: %v", initErr))
		}
	})
	return conn
}
