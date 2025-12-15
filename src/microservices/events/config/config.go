package config

import (
	"os"
)

type Config struct {
	Port               string
	KafkaBrokers       string
	MovieEventsTopic   string
	UserEventsTopic    string
	PaymentEventsTopic string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8000"),
		KafkaBrokers:       getEnv("KAFKA_BROKERS", ""),
		MovieEventsTopic:   getEnv("MOVIE_TOPIC", "movie-events"),
		UserEventsTopic:    getEnv("USER_TOPIC", "user-events"),
		PaymentEventsTopic: getEnv("PAYMENT_TOPIC", "payment-events"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
