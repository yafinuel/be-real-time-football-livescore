package domain

import "time"

type MatchStatus string

const (
	MatchStatusNotStarted MatchStatus = "NOT_STARTED"
	MatchStatusLive       MatchStatus = "LIVE"
	MatchStatusHalftime   MatchStatus = "HALFTIME"
	MatchStatusFinished   MatchStatus = "FINISHED"
	MatchStatusPostponed  MatchStatus = "POSTPONED"
	MatchStatusCancelled  MatchStatus = "CANCELLED"
)

type Match struct {
	ID         string
	HomeTeamID string
	AwayTeamID string
	HomeScore  int
	AwayScore  int
	Status     MatchStatus
	Minute     int
	StartTime  time.Time

	// Relasi tim yang dapat dirakit (hydrated) oleh service
	HomeTeam *Team
	AwayTeam *Team
}