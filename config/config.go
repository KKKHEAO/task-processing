package config

import "os"

type Config struct {
	Postgres PostgresConfig
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
	}
}
