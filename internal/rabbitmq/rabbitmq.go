package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Wundagor/high-throughput-data-consumer/internal/repository"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	URL      string
	VHost    string
	Queue    string
	Exchange string
}

type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	repo      repository.Repository
}

func NewConsumer(config RabbitMQConfig, repo repository.Repository) (*Consumer, error) {
	conn, err := amqp.Dial(config.URL + config.VHost)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(config.Queue, true, false, false, false, amqp.Table{
		"x-queue-mode": "lazy",
		"x-queue-type": "classic",
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:      conn,
		channel:   channel,
		queueName: config.Queue,
		repo:      repo,
	}, nil
}

func (c *Consumer) Start(ctx context.Context, numWorkers int) error {
	messages, err := c.channel.Consume(c.queueName, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for msg := range messages {
				select {
				case <-ctx.Done():
					return
				default:
					var data repository.SourceData

					if err := json.Unmarshal(msg.Body, &data); err != nil {
						log.Printf("Failed to unmarshal message: %v", err)

						continue
					}

					if err := c.repo.InsertData(ctx, data); err != nil {
						log.Printf("Failed to insert data: %v", err)

						continue
					}
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}

	return c.conn.Close()
}
