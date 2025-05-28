package consumer

import (
	"context"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vadim8q258475/store-cart-microservice/config"
	service "github.com/vadim8q258475/store-cart-microservice/iternal/service/cart"
	"go.uber.org/zap"
)

type rabbitMQConsumer struct {
	channel     *amqp.Channel
	queueName   string
	cartService service.CartService
	logger      *zap.Logger
}

func NewRabbitMQConsumer(channel *amqp.Channel, service service.CartService, cfg config.Config, logger *zap.Logger) Consumer {
	return &rabbitMQConsumer{
		channel:     channel,
		cartService: service,
		queueName:   cfg.RabbitMQQueueName,
		logger:      logger,
	}
}

func (c *rabbitMQConsumer) Subscribe(queueName string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queueName,
		"",    // consumer (автогенерация ID)
		false, // autoAck (выкл авто-подтверждение)
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
}

func (c *rabbitMQConsumer) Listen(ctx context.Context) error {
	msgs, err := c.Subscribe(c.queueName)
	if err != nil {
		return err
	}
	for {
		select {
		case msg, ok := <-msgs:
			if ok {
				go func(msg amqp.Delivery) {
					id, err := strconv.Atoi(string(msg.Body))
					if err != nil {
						c.logger.Error(err.Error())
					}
					timeoutCtx, cancel := context.WithTimeout(ctx, createTimeoutMS*time.Millisecond)
					defer cancel()
					_, err = c.cartService.Create(timeoutCtx, uint32(id))
					if err != nil {
						c.logger.Error(err.Error())
					}
				}(msg)
			} else {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}

	}
}
