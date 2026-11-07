package service

import (
	"context"
	"football-live-score/internal/domain"
	"football-live-score/internal/repository"
	"errors"
	"fmt"
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

// UpdateScore memvalidasi aturan skor dan memperbarui data pertandingan
func (s *MatchService) UpdateScore(ctx context.Context, matchID string, homeScore, awayScore int) (*domain.Match, error) {
	// 1. Validasi Business Rules
	if homeScore < 0 || awayScore < 0 {
		return nil, errors.New("skor tidak boleh bernilai negatif")
	}

	// 2. Cek apakah pertandingan ada di database
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("pertandingan dengan ID %s tidak ditemukan", matchID)
	}

	// 3. Update skor ke database
	updatedMatch, err := s.matchRepo.UpdateScore(ctx, match.ID, homeScore, awayScore)
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui skor: %w", err)
	}

	// 4. Hydrate data tim
	if homeTeam, err := s.teamRepo.GetByID(ctx, updatedMatch.HomeTeamID); err == nil {
		updatedMatch.HomeTeam = homeTeam
	}
	if awayTeam, err := s.teamRepo.GetByID(ctx, updatedMatch.AwayTeamID); err == nil {
		updatedMatch.AwayTeam = awayTeam
	}

	return updatedMatch, nil
}