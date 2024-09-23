package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Wundagor/high-throughput-data-consumer/internal/db"
	"github.com/Wundagor/high-throughput-data-consumer/internal/rabbitmq"
	"github.com/Wundagor/high-throughput-data-consumer/internal/repository"
)

const numWorkers = 5

func main() {
	database, err := db.Connect()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := repository.NewRepository(database)

	rabbitConfigs := []rabbitmq.RabbitMQConfig{
		{URL: os.Getenv("RABBITMQ_URL"), VHost: os.Getenv("RABBITMQ_VHOST"), Queue: os.Getenv("RABBITMQ_QUEUE")},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, config := range rabbitConfigs {
		go func(cfg rabbitmq.RabbitMQConfig) {
			consumer, err := rabbitmq.NewConsumer(cfg, repo)

			if err != nil {
				log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
			}

			defer consumer.Close()

			if err := consumer.Start(ctx, numWorkers); err != nil {
				log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
			}
		}(config)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
}
