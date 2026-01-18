package rabbitmq

import (
	"context"
	"core-consumer/config"
	"core-consumer/internal/app/infra/logger"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	rmq           *RabbitMQ
	queueName     string
	handlers      map[string]JobHandler
	handlerLock   sync.RWMutex
	workerCount   int
	loggerService *logger.Logger
}

type JobHandler func(ctx context.Context, job *Job) error

func NewConsumer(rmq *RabbitMQ, cfg *config.Config, loggerService *logger.Logger) (*Consumer, error) {
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
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.Qos(cfg.RabbitMQConsumeWorkersCount*2, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &Consumer{
		rmq:           rmq,
		queueName:     cfg.RabbitMQQueueName,
		handlers:      make(map[string]JobHandler),
		workerCount:   cfg.RabbitMQConsumeWorkersCount,
		loggerService: loggerService,
	}, nil
}

func (c *Consumer) RegisterHandler(jobType string, handler JobHandler) {
	c.handlerLock.Lock()
	defer c.handlerLock.Unlock()
	c.handlers[jobType] = handler
}

func (c *Consumer) Start(ctx context.Context) error {
	ch := c.rmq.GetChannel()

	msgs, err := ch.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %s", err)
	}

	for i := range c.workerCount {
		go c.worker(ctx, msgs, i)
	}

	c.loggerService.Info(
		"Started workers for queue",
		slog.Int("workers", c.workerCount),
		slog.String("queue", c.queueName),
	)
	return nil
}

func (c *Consumer) worker(ctx context.Context, msgs <-chan amqp.Delivery, workerID int) {
	for {
		select {
		case <-ctx.Done():
			c.loggerService.Debug(
				fmt.Sprintf("Worker %d stopped", workerID),
			)
			return
		case msg, ok := <-msgs:
			if !ok {
				c.loggerService.Info(
					"Worker channel closed",
					slog.Int("worker_id", workerID),
				)
				return
			}

			err := c.handleMessage(ctx, msg)
			if err != nil {
				c.loggerService.Error(
					fmt.Sprintf("Worker %d: error processing message: %v", workerID, err),
				)
				_ = msg.Nack(false, false)
			} else {
				c.loggerService.Debug(
					fmt.Sprintf("Worker %d: message processed successfully", workerID),
				)
				_ = msg.Ack(false)
			}
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg amqp.Delivery) error {
	var job Job
	err := json.Unmarshal(msg.Body, &job)
	if err != nil {
		return fmt.Errorf("failed to unmarshal job: %s", err)
	}

	c.handlerLock.RLock()
	handler, exists := c.handlers[job.Job]
	c.handlerLock.RUnlock()

	if !exists {
		return fmt.Errorf("no handler for job type: %s", job.Job)
	}

	return handler(ctx, &job)
}
