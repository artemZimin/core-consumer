package rabbitmq

import (
	"core-consumer/config"
	"core-consumer/internal/app/infra/logger"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func Connect(
	cfg *config.Config,
	loggerService *logger.Logger,
) (*RabbitMQ, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		if err = conn.Close(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to open channel: %s", err)
	}

	loggerService.Info("rabbitmq connected")

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.ch
}
