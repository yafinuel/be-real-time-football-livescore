package main

import (
	"context"
	"log"
	"net/http"

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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	ctx := context.Background()

	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer dbPool.Close()

	if err := database.RunMigrations(ctx, dbPool); err != nil {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}

	queries := db.New(dbPool)
	teamRepo := postgres.NewTeamRepository(queries)
	matchRepo := postgres.NewMatchRepository(queries)
	eventRepo := postgres.NewEventRepository(queries)

	matchService := service.NewMatchService(matchRepo, teamRepo)
	eventService := service.NewEventService(eventRepo, matchRepo)
	teamService := service.NewTeamService(teamRepo)
	broker := pubsub.NewBroker()

	sim := engine.NewSimulator(matchService, eventService, broker, cfg.SimulatorInterval)
	go sim.Start(ctx)

	http.HandleFunc("/healthy", HealthHandler(dbPool))
	http.HandleFunc("/healthy/", HealthHandler(dbPool))

	resolver := &graph.Resolver{
		MatchService: matchService,
		EventService: eventService,
		TeamService:  teamService,
		Broker:       broker,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Server siap di port :%s (env: %s, simulator interval: %v)", cfg.Port, cfg.AppEnv, cfg.SimulatorInterval)
	log.Printf("Health check: http://localhost:%s/healthy", cfg.Port)
	log.Printf("GraphQL Playground: http://localhost:%s/", cfg.Port)

	corsHandler := middleware.Cors(http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsHandler))
}
