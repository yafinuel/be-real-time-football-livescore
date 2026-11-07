package service_test

import (
	"context"
	"testing"

	"football-live-score/internal/domain"
	"football-live-score/internal/service"
	"football-live-score/tests/testhelper"
)

func TestEventService_AddEvent_Success(t *testing.T) {
	eventRepo := testhelper.NewMockEventRepository()
	matchRepo := testhelper.NewMockMatchRepository()
	svc := service.NewEventService(eventRepo, matchRepo)

	detail := "Assist by Odegaard"
	event, err := svc.AddEvent(
		context.Background(),
		"match-1",
		23,
		domain.EventTypeGoal,
		"team-1",
		"Bukayo Saka",
		&detail,
	)

	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}
	if event.Player != "Bukayo Saka" {
		t.Errorf("ekspektasi player Bukayo Saka, got: %s", event.Player)
	}
	if event.Type != domain.EventTypeGoal {
		t.Errorf("ekspektasi type GOAL, got: %s", event.Type)
	}
}

func TestEventService_AddEvent_InvalidMinute(t *testing.T) {
	eventRepo := testhelper.NewMockEventRepository()
	matchRepo := testhelper.NewMockMatchRepository()
	svc := service.NewEventService(eventRepo, matchRepo)

	_, err := svc.AddEvent(
		context.Background(),
		"match-1",
		-5, // menit tidak valid
		domain.EventTypeGoal,
		"team-1",
		"Bukayo Saka",
		nil,
	)

	if err == nil {
		t.Fatal("ekspektasi error untuk menit negatif")
	}
}

func TestEventService_AddEvent_MinuteTooHigh(t *testing.T) {
	eventRepo := testhelper.NewMockEventRepository()
	matchRepo := testhelper.NewMockMatchRepository()
	svc := service.NewEventService(eventRepo, matchRepo)

	_, err := svc.AddEvent(
		context.Background(),
		"match-1",
		200, // menit lebih dari 130
		domain.EventTypeGoal,
		"team-1",
		"Bukayo Saka",
		nil,
	)

	if err == nil {
		t.Fatal("ekspektasi error untuk menit > 130")
	}
}

func TestEventService_AddEvent_MatchNotFound(t *testing.T) {
	eventRepo := testhelper.NewMockEventRepository()
	matchRepo := testhelper.NewMockMatchRepository()
	svc := service.NewEventService(eventRepo, matchRepo)

	_, err := svc.AddEvent(
		context.Background(),
		"match-tidak-ada",
		23,
		domain.EventTypeGoal,
		"team-1",
		"Bukayo Saka",
		nil,
	)

	if err == nil {
		t.Fatal("ekspektasi error ketika match tidak ditemukan")
	}
}

func TestEventService_AddEvent_InvalidTeam(t *testing.T) {
	eventRepo := testhelper.NewMockEventRepository()
	matchRepo := testhelper.NewMockMatchRepository()
	svc := service.NewEventService(eventRepo, matchRepo)

	_, err := svc.AddEvent(
		context.Background(),
		"match-1",
		23,
		domain.EventTypeGoal,
		"team-999", // tim tidak bertanding di match ini
		"Random Player",
		nil,
	)

	if err == nil {
		t.Fatal("ekspektasi error ketika tim tidak bertanding di pertandingan ini")
	}
}

func TestEventService_GetEventsByMatchID(t *testing.T) {
	eventRepo := testhelper.NewMockEventRepository()
	matchRepo := testhelper.NewMockMatchRepository()
	svc := service.NewEventService(eventRepo, matchRepo)

	// Tambahkan 2 event terlebih dahulu
	svc.AddEvent(context.Background(), "match-1", 23, domain.EventTypeGoal, "team-1", "Saka", nil)
	svc.AddEvent(context.Background(), "match-1", 50, domain.EventTypeYellowCard, "team-2", "Reece James", nil)

	events, err := svc.GetEventsByMatchID(context.Background(), "match-1")
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("ekspektasi 2 event, got: %d", len(events))
	}
}
