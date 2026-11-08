package engine

import (
	"context"
	"log"
	"math/rand"
	"time"

	"football-live-score/internal/domain"
	"football-live-score/internal/pubsub"
	"football-live-score/internal/service"
)

// Simulator mengontrol jalannya pertandingan secara otomatis di background
type Simulator struct {
	matchService *service.MatchService
	eventService *service.EventService
	broker       *pubsub.Broker
	interval     time.Duration
}

// NewSimulator membuat instance simulator baru
func NewSimulator(
	matchService *service.MatchService,
	eventService *service.EventService,
	broker *pubsub.Broker,
	interval time.Duration,
) *Simulator {
	return &Simulator{
		matchService: matchService,
		eventService: eventService,
		broker:       broker,
		interval:     interval,
	}
}

// Start menjalankan loop simulasi di background sampai context dibatalkan
func (s *Simulator) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("[Simulator] Match Engine dimulai (interval: %v)...", s.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Simulator] Match Engine dihentikan (graceful shutdown).")
			return
		case <-ticker.C:
			s.simulateTick(ctx)
		}
	}
}

// simulateTick mengeksekusi 1 langkah simulasi untuk semua pertandingan LIVE
func (s *Simulator) simulateTick(ctx context.Context) {
	liveMatches, err := s.matchService.GetLiveMatches(ctx)
	if err != nil {
		log.Printf("[Simulator] Gagal mengambil live matches: %v", err)
		return
	}

	if len(liveMatches) == 0 {
		return // Tidak ada pertandingan LIVE saat ini
	}

	for _, m := range liveMatches {
		// 1. Jika sudah menit >= 90, selesaikan pertandingan
		if m.Minute >= 90 {
			finished, err := s.matchService.UpdateStatus(ctx, m.ID, domain.MatchStatusFinished, &m.Minute)
			if err == nil && s.broker != nil {
				s.broker.PublishMatch(finished)
				log.Printf("[Simulator] Peluit panjang! Match %s SELESAI (FT).", m.ID)
			}
			continue
		}

		// 2. Tambah menit pertandingan (+1 atau +2 menit)
		newMinute := m.Minute + 1
		updated, err := s.matchService.UpdateStatus(ctx, m.ID, domain.MatchStatusLive, &newMinute)
		if err != nil {
			log.Printf("[Simulator] Gagal update menit match %s: %v", m.ID, err)
			continue
		}
		m = updated

		// 3. Peluang gol secara acak (30% peluang per tick)
		if rand.Float32() < 0.30 {
			s.handleRandomGoal(ctx, m)
		} else {
			// Jika tidak ada gol, tetap siarkan update menit ke subscriber
			if s.broker != nil {
				s.broker.PublishMatch(m)
			}
			log.Printf("[Simulator] Match %s berjalan di menit %d'", m.ID, m.Minute)
		}
	}
}

// handleRandomGoal memilih tim dan pencetak gol secara acak, lalu menyiarkan skor dan event baru
func (s *Simulator) handleRandomGoal(ctx context.Context, m *domain.Match) {
	isHome := rand.Intn(2) == 0

	var teamID string
	var player string
	newHomeScore := m.HomeScore
	newAwayScore := m.AwayScore

	if isHome {
		newHomeScore++
		teamID = m.HomeTeamID
		player = "Pemain Tuan Rumah"
		if m.HomeTeam != nil {
			player = m.HomeTeam.ShortName + " Striker"
		}
	} else {
		newAwayScore++
		teamID = m.AwayTeamID
		player = "Pemain Tamu"
		if m.AwayTeam != nil {
			player = m.AwayTeam.ShortName + " Striker"
		}
	}

	// Update skor di database
	scoredMatch, err := s.matchService.UpdateScore(ctx, m.ID, newHomeScore, newAwayScore)
	if err != nil {
		log.Printf("[Simulator] Gagal update skor: %v", err)
		return
	}

	// Simpan event gol di database
	detailMsg := "Gol hasil tembakan terarah!"
	event, err := s.eventService.AddEvent(ctx, m.ID, m.Minute, domain.EventTypeGoal, teamID, player, &detailMsg)
	if err != nil {
		log.Printf("[Simulator] Gagal mencatat event gol: %v", err)
	}

	// Broadcast ke WebSocket real-time!
	if s.broker != nil {
		s.broker.PublishMatch(scoredMatch)
		if event != nil {
			s.broker.PublishEvent(event)
		}
	}

	log.Printf("[Simulator] ⚽ GOL! %s mencetak gol di menit %d! Skor sekarang: %d - %d",
		player, m.Minute, newHomeScore, newAwayScore)
}
