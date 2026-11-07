package postgres

import (
	"context"
	db "football-live-score/internal/database/sqlc"
	"football-live-score/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

type EventRepository struct {
	q *db.Queries
}

func NewEventRepository(q *db.Queries) *EventRepository {
	return &EventRepository{q: q}
}

func (r *EventRepository) Create(ctx context.Context, event *domain.MatchEvent) (*domain.MatchEvent, error) {
	var detail pgtype.Text
	if event.Detail != nil {
		detail = pgtype.Text{String: *event.Detail, Valid: true}
	}

	row, err := r.q.CreateMatchEvent(ctx, db.CreateMatchEventParams{
		ID:      event.ID,
		MatchID: event.MatchID,
		Minute:  int32(event.Minute),
		Type:    string(event.Type),
		TeamID:  event.TeamID,
		Player:  event.Player,
		Detail:  detail,
	})
	if err != nil {
		return nil, err
	}

	return toDomainEvent(&row), nil
}

func (r *EventRepository) ListByMatchID(ctx context.Context, matchID string) ([]*domain.MatchEvent, error) {
	rows, err := r.q.ListMatchEventsByMatchID(ctx, matchID)
	if err != nil {
		return nil, err
	}

	events := make([]*domain.MatchEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, toDomainEvent(&row))
	}
	return events, nil
}

func toDomainEvent(e *db.MatchEvent) *domain.MatchEvent {
	var detail *string
	if e.Detail.Valid {
		detail = &e.Detail.String
	}

	return &domain.MatchEvent{
		ID:      e.ID,
		MatchID: e.MatchID,
		Minute:  int(e.Minute),
		Type:    domain.EventType(e.Type),
		TeamID:  e.TeamID,
		Player:  e.Player,
		Detail:  detail,
	}
}