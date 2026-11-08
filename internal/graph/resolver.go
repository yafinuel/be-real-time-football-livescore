package graph

import (
	"football-live-score/internal/pubsub"
	"football-live-score/internal/service"
)

type Resolver struct {
	MatchService *service.MatchService
	EventService *service.EventService
	TeamService  *service.TeamService
	Broker       *pubsub.Broker
}
