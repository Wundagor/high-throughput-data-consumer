package consumer

import (
	"log"
	"sync"

	"github.com/Wundagor/high-throughput-data-consumer/internal/config"
	"github.com/Wundagor/high-throughput-data-consumer/internal/database"
	"github.com/Wundagor/high-throughput-data-consumer/internal/mq"

	"github.com/streadway/amqp"
)

type WorkerPool struct {
	workerCount int
	rabbitMQ    *mq.RabbitMQ
	dbConfig    config.DatabaseConfig
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewWorkerPool(workerCount int, rabbitMQ *mq.RabbitMQ, dbConfig config.DatabaseConfig) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		rabbitMQ:    rabbitMQ,
		dbConfig:    dbConfig,
		stopChan:    make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.stopChan)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	msgs, err := wp.rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Worker %d: failed to start consuming: %v", id, err)
	}

	db, err := database.Connect(wp.dbConfig)
	if err != nil {
		log.Fatalf("Worker %d: failed to connect to the database: %v", id, err)
	}
	defer db.Close()

	for {
		select {
		case msg := <-msgs:
			if err := handleMessage(db, msg); err != nil {
				log.Printf("Worker %d: error handling message: %v", id, err)
				msg.Nack(false, false)
			} else {
				msg.Ack(false)
			}
		case <-wp.stopChan:
			log.Printf("Worker %d stopping...", id)
			return
		}
	}
}

func handleMessage(db *database.DB, msg amqp.Delivery) error {
	data, err := decodeMessage(msg.Body)
	if err != nil {
		return err
	}

	if err := db.Save(data); err != nil {
		return err
	}

	return nil
}
