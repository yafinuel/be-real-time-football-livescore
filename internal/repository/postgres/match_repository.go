package postgres

import (
	"context"
	db "football-live-score/internal/database/sqlc"
	"football-live-score/internal/domain"
)

type MatchRepository struct {
	q *db.Queries
}

func NewMatchRepository(q *db.Queries) *MatchRepository {
	return &MatchRepository{q: q}
}

func (r *MatchRepository) GetByID(ctx context.Context, id string) (*domain.Match, error) {
	row, err := r.q.GetMatchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDomainMatch(&row), nil
}

func (r *MatchRepository) List(ctx context.Context) ([]*domain.Match, error) {
	rows, err := r.q.ListMatches(ctx)
	if err != nil {
		return nil, err
	}

	matches := make([]*domain.Match, 0, len(rows))
	for _, row := range rows {
		matches = append(matches, toDomainMatch(&row))
	}

	return matches, nil
}

// Helper untuk konversi struct sqlc -> domain model
func toDomainMatch(m *db.Match) *domain.Match {
	return &domain.Match{
		ID:         m.ID,
		HomeTeamID: m.HomeTeamID,
		AwayTeamID: m.AwayTeamID,
		HomeScore:  int(m.HomeScore),
		AwayScore:  int(m.AwayScore),
		Status:     domain.MatchStatus(m.Status),
		Minute:     int(m.Minute),
		StartTime:  m.StartTime.Time,
	}
}