package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr       string
	DSN            string
	MigrationsPath string

	RabbitMQURL string
	RedisURL    string
}

func Read() Config {
	var config Config

	httpAddr, exists := os.LookupEnv("HTTP_ADDR")
	fmt.Println(httpAddr)
	if exists {
		config.HTTPAddr = httpAddr
	} else {
		config.HTTPAddr = ":8080"
	}

	migrationsPath, exists := os.LookupEnv("MIGRATIONS_PATH")
	if exists {
		config.MigrationsPath = migrationsPath
	} else {
		config.MigrationsPath = "file://./internal/app/migrations"
	}

	dsn, exists := os.LookupEnv("DSN")
	if exists {
		config.DSN = dsn
	} else {
		config.DSN = "postgres://user:password@127.0.0.1:5432/payments?sslmode=disable"
	}

	rabbitUrl, exists := os.LookupEnv("RABBITMQ_URL")

	if exists {
		config.RabbitMQURL = rabbitUrl
	} else {
		config.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	redisUrl, exists := os.LookupEnv("REDIS_URL")
	if exists {
		config.RedisURL = redisUrl
	} else {
		config.RedisURL = "localhost:6379"
	}
	return config
}
