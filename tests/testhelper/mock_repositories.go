// Package testhelper menyediakan mock repository untuk keperluan unit testing
package testhelper

import (
	"context"
	"errors"
	"football-live-score/internal/domain"
)

// --- Mock Team Repository ---

type MockTeamRepository struct {
	Teams  map[string]*domain.Team
	ShouldFail bool
}

func NewMockTeamRepository() *MockTeamRepository {
	return &MockTeamRepository{
		Teams: map[string]*domain.Team{
			"team-1": {ID: "team-1", Name: "Arsenal", ShortName: "ARS"},
			"team-2": {ID: "team-2", Name: "Chelsea", ShortName: "CHE"},
		},
	}
}

func (m *MockTeamRepository) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	team, ok := m.Teams[id]
	if !ok {
		return nil, errors.New("team not found")
	}
	return team, nil
}

func (m *MockTeamRepository) List(ctx context.Context) ([]*domain.Team, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	result := make([]*domain.Team, 0, len(m.Teams))
	for _, t := range m.Teams {
		result = append(result, t)
	}
	return result, nil
}

// --- Mock Match Repository ---

type MockMatchRepository struct {
	Matches    map[string]*domain.Match
	ShouldFail bool
}

func NewMockMatchRepository() *MockMatchRepository {
	return &MockMatchRepository{
		Matches: map[string]*domain.Match{
			"match-1": {
				ID:         "match-1",
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				HomeScore:  2,
				AwayScore:  1,
				Status:     domain.MatchStatusLive,
				Minute:     75,
			},
		},
	}
}

func (m *MockMatchRepository) GetByID(ctx context.Context, id string) (*domain.Match, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	match, ok := m.Matches[id]
	if !ok {
		return nil, errors.New("match not found")
	}
	return match, nil
}

func (m *MockMatchRepository) List(ctx context.Context) ([]*domain.Match, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	result := make([]*domain.Match, 0, len(m.Matches))
	for _, match := range m.Matches {
		result = append(result, match)
	}
	return result, nil
}

func (m *MockMatchRepository) UpdateScore(ctx context.Context, id string, homeScore, awayScore int) (*domain.Match, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	match, ok := m.Matches[id]
	if !ok {
		return nil, errors.New("match not found")
	}
	match.HomeScore = homeScore
	match.AwayScore = awayScore
	return match, nil
}

func (m *MockMatchRepository) UpdateStatus(ctx context.Context, id string, status domain.MatchStatus, minute *int) (*domain.Match, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	match, ok := m.Matches[id]
	if !ok {
		return nil, errors.New("match not found")
	}
	match.Status = status
	if minute != nil {
		match.Minute = *minute
	}
	return match, nil
}

// --- Mock Event Repository ---

type MockEventRepository struct {
	Events     []*domain.MatchEvent
	ShouldFail bool
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		Events: []*domain.MatchEvent{},
	}
}

func (m *MockEventRepository) Create(ctx context.Context, event *domain.MatchEvent) (*domain.MatchEvent, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	m.Events = append(m.Events, event)
	return event, nil
}

func (m *MockEventRepository) ListByMatchID(ctx context.Context, matchID string) ([]*domain.MatchEvent, error) {
	if m.ShouldFail {
		return nil, errors.New("database error")
	}
	var result []*domain.MatchEvent
	for _, e := range m.Events {
		if e.MatchID == matchID {
			result = append(result, e)
		}
	}
	return result, nil
}
