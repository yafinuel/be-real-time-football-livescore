package service

import (
	"context"
	"errors"
	"fmt"
	"football-live-score/internal/domain"
	"football-live-score/internal/repository"
	"time"
)

type EventService struct {
	eventRepo repository.EventRepository
	matchRepo repository.MatchRepository
}

func NewEventService(eventRepo repository.EventRepository, matchRepo repository.MatchRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		matchRepo: matchRepo,
	}
}

func (s *EventService) AddEvent(ctx context.Context, matchID string, minute int, eventType domain.EventType, teamID string, player string, detail *string) (*domain.MatchEvent, error) {
	if minute < 0 || minute > 130 {
		return nil, errors.New("menit pertandingan tidak valid (harus antara 0 - 130)")
	}

	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("pertandingan dengan ID %s tidak ditemukan", matchID)
	}

	if teamID != match.HomeTeamID && teamID != match.AwayTeamID {
		return nil, errors.New("tim yang dipilih tidak bertanding di pertandingan ini")
	}

	eventID := fmt.Sprintf("event-%d", time.Now().UnixNano())

	event := &domain.MatchEvent{
		ID:      eventID,
		MatchID: matchID,
		Minute:  minute,
		Type:    eventType,
		TeamID:  teamID,
		Player:  player,
		Detail:  detail,
	}

	return s.eventRepo.Create(ctx, event)
}

func (s *EventService) GetEventsByMatchID(ctx context.Context, matchID string) ([]*domain.MatchEvent, error) {
	return s.eventRepo.ListByMatchID(ctx, matchID)
}
