package config_test

import (
	"os"
	"testing"
	"time"

	"football-live-score/internal/config"
)

func TestLoad_Success(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Setenv("PORT", "9090")
	os.Setenv("APP_ENV", "production")
	os.Setenv("SIMULATOR_INTERVAL", "30s")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("PORT")
		os.Unsetenv("APP_ENV")
		os.Unsetenv("SIMULATOR_INTERVAL")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("DatabaseURL tidak sesuai: %s", cfg.DatabaseURL)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port tidak sesuai: %s", cfg.Port)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("AppEnv tidak sesuai: %s", cfg.AppEnv)
	}
	if cfg.SimulatorInterval != 30*time.Second {
		t.Errorf("SimulatorInterval tidak sesuai: %v", cfg.SimulatorInterval)
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Unsetenv("PORT")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("SIMULATOR_INTERVAL")
	defer os.Unsetenv("DATABASE_URL")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Port default seharusnya 8080, got: %s", cfg.Port)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("AppEnv default seharusnya development, got: %s", cfg.AppEnv)
	}
	if cfg.SimulatorInterval != 45*time.Second {
		t.Errorf("SimulatorInterval default seharusnya 45s, got: %v", cfg.SimulatorInterval)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")

	_, err := config.Load()
	if err == nil {
		t.Fatal("ekspektasi error ketika DATABASE_URL tidak diisi, tapi tidak ada error")
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	cfg := &config.Config{AppEnv: "development"}
	if !cfg.IsDevelopment() {
		t.Error("IsDevelopment() seharusnya true untuk environment development")
	}
}

func TestConfig_IsProduction(t *testing.T) {
	cfg := &config.Config{AppEnv: "production"}
	if !cfg.IsProduction() {
		t.Error("IsProduction() seharusnya true untuk environment production")
	}
}
