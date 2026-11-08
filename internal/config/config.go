package config

import (
	"fmt"
	"os"
	"time"
)

// Config menyimpan seluruh konfigurasi aplikasi yang dibaca dari environment variable
type Config struct {
	Port              string
	DatabaseURL       string
	AppEnv            string
	SimulatorInterval time.Duration
}

// Load membaca environment variable dan mengembalikan Config
// Akan mengembalikan error jika variabel wajib tidak tersedia
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable wajib diisi")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	simInterval := 45 * time.Second
	if intervalStr := os.Getenv("SIMULATOR_INTERVAL"); intervalStr != "" {
		if parsed, err := time.ParseDuration(intervalStr); err == nil {
			simInterval = parsed
		}
	}

	return &Config{
		Port:              port,
		DatabaseURL:       dbURL,
		AppEnv:            appEnv,
		SimulatorInterval: simInterval,
	}, nil
}

// IsDevelopment mengembalikan true jika environment adalah development
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// IsProduction mengembalikan true jika environment adalah production
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}
