package service

import (
	"context"
	"football-live-score/internal/domain"
	"football-live-score/internal/repository"
)

type MatchService struct {
	matchRepo repository.MatchRepository
	teamRepo  repository.TeamRepository
}

func NewMatchService(matchRepo repository.MatchRepository, teamRepo repository.TeamRepository) *MatchService {
	return &MatchService{
		matchRepo: matchRepo,
		teamRepo:  teamRepo,
	}
}

// ListMatches mengambil pertandingan dan merakit informasi HomeTeam & AwayTeam
func (s *MatchService) ListMatches(ctx context.Context) ([]*domain.Match, error) {
	matches, err := s.matchRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Lengkapi (hydrate) data tim untuk setiap pertandingan
	for _, m := range matches {
		if homeTeam, err := s.teamRepo.GetByID(ctx, m.HomeTeamID); err == nil {
			m.HomeTeam = homeTeam
		}
		if awayTeam, err := s.teamRepo.GetByID(ctx, m.AwayTeamID); err == nil {
			m.AwayTeam = awayTeam
		}
	}

	return matches, nil
}