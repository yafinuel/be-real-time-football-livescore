package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"football-live-score/graph"
	"football-live-score/internal/database"
	db "football-live-score/internal/database/sqlc"
	"football-live-score/internal/repository/postgres"
	"football-live-score/internal/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable wajib diisi")
	}

	ctx := context.Background()

	// 1. Inisialisasi Database Pool
	dbPool, err := database.NewPostgresPool(ctx, dbURL)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer dbPool.Close()

	// 2. Jalankan Migrasi Otomatis
	if err := database.RunMigrations(ctx, dbPool); err != nil {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}

	// 3. Rakit Layer Arsitektur (Clean Architecture)
	queries := db.New(dbPool)
	teamRepo := postgres.NewTeamRepository(queries)
	matchRepo := postgres.NewMatchRepository(queries)
	matchService := service.NewMatchService(matchRepo, teamRepo)

	// 4. Endpoint Health Check
	http.HandleFunc("/healthy", HealthHandler(dbPool))
	http.HandleFunc("/healthy/", HealthHandler(dbPool))

	// 5. GraphQL Server & Playground
	resolver := &graph.Resolver{
		MatchService: matchService,
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Server siap di port :%s", port)
	log.Printf("Health check: http://localhost:%s/healthy", port)
	log.Printf("GraphQL Playground: http://localhost:%s/", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}