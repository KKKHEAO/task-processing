package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Postgres PostgresConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  string
	PgDriver           string
}

type KafkaConfig struct {
	Brokers      []string
	MainTopic    string
	RetryTopics  []RetryTopicConfig
	DLQTopic     string
	MaxRetries   int
	BatchSize    int
	PollInterval time.Duration
}

type RetryTopicConfig struct {
	Name     string
	Delay    time.Duration
	MaxRetry int // Максимальный номер попытки для этого топика
}

func NewConfig() *Config {
	return &Config{
		Postgres: PostgresConfig{
			PostgresqlHost:     os.Getenv("DB_HOST"),
			PostgresqlPort:     os.Getenv("DB_PORT"),
			PostgresqlUser:     os.Getenv("DB_USER"),
			PostgresqlPassword: os.Getenv("DB_PASSWORD"),
			PostgresqlDbname:   os.Getenv("DB_NAME"),
			PostgresqlSSLMode:  os.Getenv("DB_SSL_MODE"),
			PgDriver:           os.Getenv("DB_DRIVER"),
		},
		Kafka: KafkaConfig{
			Brokers:   strings.Split(getEnvOrDefault("KAFKA_BROKERS", "kafka:9092"), ","),
			MainTopic: getEnvOrDefault("KAFKA_MAIN_TOPIC", "tasks.created"),
			RetryTopics: []RetryTopicConfig{
				{Name: "tasks.retry.1m", Delay: 1 * time.Minute, MaxRetry: 1},
				{Name: "tasks.retry.5m", Delay: 5 * time.Minute, MaxRetry: 2},
				{Name: "tasks.retry.30m", Delay: 30 * time.Minute, MaxRetry: 3},
			},
			DLQTopic:     getEnvOrDefault("KAFKA_DLQ_TOPIC", "tasks.dlq"),
			MaxRetries:   getEnvAsInt("KAFKA_MAX_RETRIES", 4), // 3 retry + 1 initial
			BatchSize:    getEnvAsInt("KAFKA_BATCH_SIZE", 100),
			PollInterval: getEnvAsDuration("KAFKA_POLL_INTERVAL", 2*time.Second),
		},
	}
}

func (c *Config) Validate() error {
	if len(c.Kafka.Brokers) == 0 {
		return errors.New("KAFKA_BROKERS is required")
	}
	if c.Postgres.PostgresqlHost == "" {
		return errors.New("DB_HOST is required")
	}
	if c.Postgres.PostgresqlUser == "" {
		return errors.New("DB_USER is required")
	}
	if c.Postgres.PostgresqlPassword == "" {
		return errors.New("DB_PASSWORD is required")
	}
	if c.Postgres.PostgresqlDbname == "" {
		return errors.New("DB_NAME is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
