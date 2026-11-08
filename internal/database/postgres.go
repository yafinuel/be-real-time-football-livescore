package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// NewPostgresPool membuat pool koneksi dan memverifikasi koneksi database
func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("gagal parse connection string: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pool database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("database tidak dapat dijangkau: %w", err)
	}

	log.Println("Berhasil terhubung ke database PostgreSQL!")
	return pool, nil
}

// RunMigrations mengeksekusi semua file .sql dari folder migrations
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationFiles := []string{
		"migrations/001_create_teams.sql",
		"migrations/002_create_matches.sql",
		"migrations/003_create_match_events.sql",
		"migrations/004_seed_data.sql",
	}

	for _, file := range migrationFiles {
		content, err := migrationFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("gagal membaca file migrasi %s: %w", file, err)
		}

		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("gagal mengeksekusi migrasi %s: %w", file, err)
		}
		log.Printf("Migrasi sukses: %s", file)
	}

	return nil
}
