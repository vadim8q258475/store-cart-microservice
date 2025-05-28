package consumer

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

const createTimeoutMS = 300

type Consumer interface {
	Subscribe(queueName string) (<-chan amqp.Delivery, error)
	Listen(ctx context.Context) error
}
