package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Wundagor/high-throughput-data-consumer/internal/config"
	"github.com/Wundagor/high-throughput-data-consumer/internal/consumer"
	"github.com/Wundagor/high-throughput-data-consumer/internal/mq"
)

func main() {
	cfg := config.LoadConfig()

	rabbitMQ, err := mq.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ instance: %v", err)
	}
	defer rabbitMQ.Close()

	workerPool := consumer.NewWorkerPool(cfg.Consumer.WorkerCount, rabbitMQ, cfg.Database)
	workerPool.Start()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	log.Println("Shutting down workers...")
	workerPool.Stop()
	log.Println("Workers stopped.")
}
