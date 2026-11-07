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

func (s *MatchService) UpdateStatus(ctx context.Context, matchID string, status domain.MatchStatus, minute *int) (*domain.Match, error) {
	// 1. Cek apakah pertandingan ada di database
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("pertandingan dengan ID %s tidak ditemukan", matchID)
	}

	// 2. Validasi transisi status (pertandingan selesai tidak boleh diubah lagi)
	if match.Status == domain.MatchStatusFinished {
		return nil, errors.New("pertandingan sudah selesai (FINISHED) dan tidak dapat diubah lagi")
	}

	// 3. Validasi menit jika diisi
	if minute != nil && *minute < 0 {
		return nil, errors.New("menit pertandingan tidak boleh bernilai negatif")
	}

	// 4. Update status ke database
	updatedMatch, err := s.matchRepo.UpdateStatus(ctx, matchID, status, minute)
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui status pertandingan: %w", err)
	}

	// 5. Hydrate data tim
	if homeTeam, err := s.teamRepo.GetByID(ctx, updatedMatch.HomeTeamID); err == nil {
		updatedMatch.HomeTeam = homeTeam
	}
	if awayTeam, err := s.teamRepo.GetByID(ctx, updatedMatch.AwayTeamID); err == nil {
		updatedMatch.AwayTeam = awayTeam
	}

	return updatedMatch, nil
}


func (s *MatchService) GetMatchByID(ctx context.Context, id string) (*domain.Match, error) {
	match, err := s.matchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("pertandingan dengan ID %s tidak ditemukan: %w", id, err)
	}

	if homeTeam, err := s.teamRepo.GetByID(ctx, match.HomeTeamID); err == nil {
		match.HomeTeam = homeTeam
	}
	if awayTeam, err := s.teamRepo.GetByID(ctx, match.AwayTeamID); err == nil {
		match.AwayTeam = awayTeam
	}

	return match, nil
}

func (s *MatchService) GetLiveMatches(ctx context.Context) ([]*domain.Match, error) {
	matches, err := s.ListMatches(ctx)
	if err != nil {
		return nil, err
	}

	live := make([]*domain.Match, 0)
	for _, m := range matches {
		if m.Status == domain.MatchStatusLive {
			live = append(live, m)
		}
	}
	return live, nil
}
