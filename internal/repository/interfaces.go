package repository

import (
	"context"
	"football-live-score/internal/domain"
)

type TeamRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	List(ctx context.Context) ([]*domain.Team, error)
}

type MatchRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Match, error)
	List(ctx context.Context) ([]*domain.Match, error)
	UpdateScore(ctx context.Context, id string, homeScore, awayScore int) (*domain.Match, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *domain.MatchEvent) (*domain.MatchEvent, error)
	ListByMatchID(ctx context.Context, matchID string) ([]*domain.MatchEvent, error)
}