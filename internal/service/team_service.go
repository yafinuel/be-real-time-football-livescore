package service

import (
	"context"
	"football-live-score/internal/domain"
	"football-live-score/internal/repository"
)

type TeamService struct {
	teamRepo repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) ListTeams(ctx context.Context) ([]*domain.Team, error) {
	return s.teamRepo.List(ctx)
}

func (s *TeamService) GetTeamByID(ctx context.Context, id string) (*domain.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}
