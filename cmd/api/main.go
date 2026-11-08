package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"football-live-score/internal/config"
	"football-live-score/internal/database"
	db "football-live-score/internal/database/sqlc"
	"football-live-score/internal/graph"
	"football-live-score/internal/middleware"
	"football-live-score/internal/pubsub"
	"football-live-score/internal/repository/postgres"
	"football-live-score/internal/service"
	"football-live-score/internal/engine" 

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	// 1. Load konfigurasi dari environment variable
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	ctx := context.Background()

	// 2. Inisialisasi Database Pool
	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer dbPool.Close()

	// 3. Jalankan Migrasi Otomatis
	if err := database.RunMigrations(ctx, dbPool); err != nil {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}

	// 4. Rakit Layer Arsitektur (Clean Architecture)
	queries := db.New(dbPool)
	teamRepo := postgres.NewTeamRepository(queries)
	matchRepo := postgres.NewMatchRepository(queries)
	eventRepo := postgres.NewEventRepository(queries)

	matchService := service.NewMatchService(matchRepo, teamRepo)
	eventService := service.NewEventService(eventRepo, matchRepo)
	teamService := service.NewTeamService(teamRepo)
	broker := pubsub.NewBroker()

	// 5. Nyalakan Mock Match Simulator di background
	sim := engine.NewSimulator(matchService, eventService, broker, 5*time.Second)
	go sim.Start(ctx)

	// 6. Endpoint Health Check
	http.HandleFunc("/healthy", HealthHandler(dbPool))
	http.HandleFunc("/healthy/", HealthHandler(dbPool))

	// 7. GraphQL Server & Playground
	resolver := &graph.Resolver{
		MatchService: matchService,
		EventService: eventService,
		TeamService:  teamService,
		Broker:       broker,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Server siap di port :%s (env: %s)", cfg.Port, cfg.AppEnv)
	log.Printf("Health check: http://localhost:%s/healthy", cfg.Port)
	log.Printf("GraphQL Playground: http://localhost:%s/", cfg.Port)

	// Bungkus seluruh router HTTP dengan CORS middleware
	corsHandler := middleware.Cors(http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsHandler))
}
