package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Pinger adalah interface minimal agar mudah di-mock pada unit test
type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Error    string `json:"error,omitempty"`
}

// HealthHandler mengembalikan status server dan konektivitas database
func HealthHandler(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(HealthResponse{
				Status:   "unhealthy",
				Database: "disconnected",
				Error:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:   "healthy",
			Database: "connected",
		})
	}
}