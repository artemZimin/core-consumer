package rabbitmq

import (
	"context"
	"core-consumer/config"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	rmq       *RabbitMQ
	queueName string
}

func NewProducer(rmq *RabbitMQ, cfg *config.Config) (*Producer, error) {
	ch := rmq.GetChannel()

	_, err := ch.QueueDeclare(
		cfg.RabbitMQQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %s", err)
	}

	return &Producer{
		rmq:       rmq,
		queueName: cfg.RabbitMQQueueName,
	}, nil
}

func (p *Producer) PublishJob(ctx context.Context, job *Job) error {
	ch := p.rmq.GetChannel()

	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %s", err)
	}

	err = ch.PublishWithContext(ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err)
	}

	return nil
}
