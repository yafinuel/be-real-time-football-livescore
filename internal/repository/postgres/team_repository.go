package postgres

import (
	"context"
	db "football-live-score/internal/database/sqlc"
	"football-live-score/internal/domain"
)

type TeamRepository struct {
	q *db.Queries
}

func NewTeamRepository(q *db.Queries) *TeamRepository {
	return &TeamRepository{q: q}
}

func (r *TeamRepository) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	row, err := r.q.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var logoURL *string
	if row.LogoUrl.Valid {
		logoURL = &row.LogoUrl.String
	}

	return &domain.Team{
		ID:        row.ID,
		Name:      row.Name,
		ShortName: row.ShortName,
		LogoURL:   logoURL,
	}, nil
}

func (r *TeamRepository) List(ctx context.Context) ([]*domain.Team, error) {
	rows, err := r.q.ListTeams(ctx)
	if err != nil {
		return nil, err
	}

	teams := make([]*domain.Team, 0, len(rows))
	for _, row := range rows {
		var logoURL *string
		if row.LogoUrl.Valid {
			logoURL = &row.LogoUrl.String
		}
		teams = append(teams, &domain.Team{
			ID:        row.ID,
			Name:      row.Name,
			ShortName: row.ShortName,
			LogoURL:   logoURL,
		})
	}

	return teams, nil
}