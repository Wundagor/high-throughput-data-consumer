package mq

import (
	"github.com/Wundagor/high-throughput-data-consumer/internal/config"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.RabbitMQConfig
}

func NewRabbitMQ(config config.RabbitMQConfig) (*RabbitMQ, error) {
	conn, err := amqp.DialConfig(config.URL+"/"+config.VHost, amqp.Config{})
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		r.config.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}
