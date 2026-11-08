package engine

import (
	"context"
	"testing"
	"time"

	"football-live-score/internal/domain"
	"football-live-score/internal/pubsub"
	"football-live-score/internal/service"
	"football-live-score/tests/testhelper"
)

func TestSimulator_SimulateTick_IncrementMinute(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	eventRepo := testhelper.NewMockEventRepository()

	matchService := service.NewMatchService(matchRepo, teamRepo)
	eventService := service.NewEventService(eventRepo, matchRepo)
	broker := pubsub.NewBroker()

	// Pastikan match-1 berstatus LIVE dan menit 10
	matchRepo.Matches["match-1"].Status = domain.MatchStatusLive
	matchRepo.Matches["match-1"].Minute = 10

	sim := NewSimulator(matchService, eventService, broker, 100*time.Millisecond)

	ctx := context.Background()
	sim.simulateTick(ctx)

	updatedMatch, err := matchService.GetMatchByID(ctx, "match-1")
	if err != nil {
		t.Fatalf("gagal get match: %v", err)
	}

	if updatedMatch.Minute <= 10 {
		t.Errorf("menit harusnya bertambah (>10), didapat: %d", updatedMatch.Minute)
	}
}

func TestSimulator_SimulateTick_FinishAt90(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	eventRepo := testhelper.NewMockEventRepository()

	matchService := service.NewMatchService(matchRepo, teamRepo)
	eventService := service.NewEventService(eventRepo, matchRepo)
	broker := pubsub.NewBroker()

	// Set match-1 di menit 90
	matchRepo.Matches["match-1"].Status = domain.MatchStatusLive
	matchRepo.Matches["match-1"].Minute = 90

	sim := NewSimulator(matchService, eventService, broker, 100*time.Millisecond)

	ctx := context.Background()
	sim.simulateTick(ctx)

	updatedMatch, err := matchService.GetMatchByID(ctx, "match-1")
	if err != nil {
		t.Fatalf("gagal get match: %v", err)
	}

	if updatedMatch.Status != domain.MatchStatusFinished {
		t.Errorf("match harusnya berstatus FINISHED saat menit >= 90, got: %s", updatedMatch.Status)
	}
}

func TestSimulator_StartAndStop(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	eventRepo := testhelper.NewMockEventRepository()

	matchService := service.NewMatchService(matchRepo, teamRepo)
	eventService := service.NewEventService(eventRepo, matchRepo)
	broker := pubsub.NewBroker()

	sim := NewSimulator(matchService, eventService, broker, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan bool)
	go func() {
		sim.Start(ctx)
		done <- true
	}()

	select {
	case <-done:
		// Berhasil berhenti saat context selesai
	case <-time.After(500 * time.Millisecond):
		t.Fatal("simulator gagal berhenti setelah timeout")
	}
}
