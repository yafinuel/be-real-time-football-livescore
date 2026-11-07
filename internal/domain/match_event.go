package domain

type EventType string

const (
	EventTypeGoal         EventType = "GOAL"
	EventTypeYellowCard   EventType = "YELLOW_CARD"
	EventTypeRedCard      EventType = "RED_CARD"
	EventTypeSubstitution EventType = "SUBSTITUTION"
)

type MatchEvent struct {
	ID      string
	MatchID string
	Minute  int
	Type    EventType
	TeamID  string
	Player  string
	Detail  *string
}