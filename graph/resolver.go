package graph

import (
	"football-live-score/internal/service"
)

// Resolver berfungsi untuk dependency injection
type Resolver struct {
	MatchService *service.MatchService
	EventService *service.EventService
	TeamService  *service.TeamService
}
