package pubsub

import (
	"context"
	"sync"

	"football-live-score/internal/domain"
)

// Broker mengelola subscription untuk match update dan match events
type Broker struct {
	mu        sync.RWMutex
	matchSubs map[string][]chan *domain.Match
	eventSubs map[string][]chan *domain.MatchEvent
}

// NewBroker membuat instance broker baru yang siap pakai
func NewBroker() *Broker {
	return &Broker{
		matchSubs: make(map[string][]chan *domain.Match),
		eventSubs: make(map[string][]chan *domain.MatchEvent),
	}
}

// ─────────────────────────────────────────────
// Match Subscriptions
// ─────────────────────────────────────────────

// SubscribeMatch mendaftarkan listener untuk suatu matchID.
// Jika matchID kosong (""), listener akan menerima update dari SEMUA pertandingan.
func (b *Broker) SubscribeMatch(ctx context.Context, matchID string) <-chan *domain.Match {
	ch := make(chan *domain.Match, 10)

	b.mu.Lock()
	b.matchSubs[matchID] = append(b.matchSubs[matchID], ch)
	b.mu.Unlock()

	// Goroutine pembersih: berjalan di background.
	// Saat koneksi user terputus (ctx.Done()), channel dihapus dari daftar.
	go func() {
		<-ctx.Done()
		b.unsubscribeMatch(matchID, ch)
	}()

	return ch
}

// PublishMatch menyiarkan update match ke semua subscriber yang relevan
func (b *Broker) PublishMatch(match *domain.Match) {
	if match == nil {
		return
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	// 1. Kirim ke subscriber spesifik pertandingan ini
	for _, ch := range b.matchSubs[match.ID] {
		select {
		case ch <- match:
		default: // Lewati jika buffer subscriber penuh agar tidak memblokir server
		}
	}

	// 2. Kirim ke subscriber global (yang mendengarkan semua match, matchID "")
	for _, ch := range b.matchSubs[""] {
		select {
		case ch <- match:
		default:
		}
	}
}

// unsubscribeMatch menghapus channel dari daftar subscriber
func (b *Broker) unsubscribeMatch(matchID string, targetCh chan *domain.Match) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs := b.matchSubs[matchID]
	for i, ch := range subs {
		if ch == targetCh {
			b.matchSubs[matchID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}

	if len(b.matchSubs[matchID]) == 0 {
		delete(b.matchSubs, matchID)
	}
}

// ─────────────────────────────────────────────
// Event Subscriptions
// ─────────────────────────────────────────────

// SubscribeEvent mendaftarkan listener untuk event (gol, kartu, dsb) pada suatu matchID
func (b *Broker) SubscribeEvent(ctx context.Context, matchID string) <-chan *domain.MatchEvent {
	ch := make(chan *domain.MatchEvent, 10)

	b.mu.Lock()
	b.eventSubs[matchID] = append(b.eventSubs[matchID], ch)
	b.mu.Unlock()

	go func() {
		<-ctx.Done()
		b.unsubscribeEvent(matchID, ch)
	}()

	return ch
}

// PublishEvent menyiarkan event baru ke semua subscriber pertandingan tersebut
func (b *Broker) PublishEvent(event *domain.MatchEvent) {
	if event == nil {
		return
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.eventSubs[event.MatchID] {
		select {
		case ch <- event:
		default:
		}
	}
}

// unsubscribeEvent menghapus channel event dari daftar subscriber
func (b *Broker) unsubscribeEvent(matchID string, targetCh chan *domain.MatchEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs := b.eventSubs[matchID]
	for i, ch := range subs {
		if ch == targetCh {
			b.eventSubs[matchID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}

	if len(b.eventSubs[matchID]) == 0 {
		delete(b.eventSubs, matchID)
	}
}
