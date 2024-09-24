package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RabbitMQ RabbitMQConfig
	Database DatabaseConfig
	Consumer ConsumerConfig
}

type RabbitMQConfig struct {
	URL   string
	VHost string
	Queue string
}

type DatabaseConfig struct {
	DSN string
}

type ConsumerConfig struct {
	Workers int
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
