package config

import (
	"log"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Port             string
	MonolithURL      *url.URL
	MoviesURL        *url.URL
	EventsURL        *url.URL
	GradualMigration bool
	MoviesPercent    int
}

func Load() *Config {
	port := getEnv("PORT", "8000")
	monolithURL := parseURL("MONOLITH_URL")
	moviesURL := parseURL("MOVIES_SERVICE_URL")
	eventsURL := parseURL("EVENTS_SERVICE_URL")
	migration := getEnv("GRADUAL_MIGRATION", "false") == "true"
	percentStr := getEnv("MOVIES_MIGRATION_PERCENT", "10")
	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		log.Fatalf("Invalid MOVIES_MIGRATION_PERCENT: %v", err)
	}

	return &Config{
		Port:             port,
		MonolithURL:      monolithURL,
		MoviesURL:        moviesURL,
		EventsURL:        eventsURL,
		GradualMigration: migration,
		MoviesPercent:    percent,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func parseURL(envKey string) *url.URL {
	raw := os.Getenv(envKey)
	parsed, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("Unable to parse %s: %v", envKey, err)
	}
	return parsed
}
