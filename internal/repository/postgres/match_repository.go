package postgres

import (
	"context"
	db "football-live-score/internal/database/sqlc"
	"football-live-score/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
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

func (r *MatchRepository) UpdateScore(ctx context.Context, id string, homeScore, awayScore int) (*domain.Match, error) {
	row, err := r.q.UpdateMatchScore(ctx, db.UpdateMatchScoreParams{
		ID:        id,
		HomeScore: int32(homeScore),
		AwayScore: int32(awayScore),
	})
	if err != nil {
		return nil, err
	}

	return toDomainMatch(&row), nil
}

func (r *MatchRepository) UpdateStatus(ctx context.Context, id string, status domain.MatchStatus, minute *int) (*domain.Match, error) {
	var minuteVal pgtype.Int4
	if minute != nil {
		minuteVal = pgtype.Int4{Int32: int32(*minute), Valid: true}
	}

	row, err := r.q.UpdateMatchStatus(ctx, db.UpdateMatchStatusParams{
		ID:     id,
		Status: string(status),
		Minute: minuteVal,
	})
	if err != nil {
		return nil, err
	}

	return toDomainMatch(&row), nil
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
