package pubsub

import (
	"context"
	"sync"
	"testing"
	"time"

	"football-live-score/internal/domain"
)

func TestBroker_SubscribeAndPublishMatch(t *testing.T) {
	broker := NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	matchChan := broker.SubscribeMatch(ctx, "match-1")

	expectedMatch := &domain.Match{
		ID:        "match-1",
		HomeScore: 2,
		AwayScore: 1,
	}

	broker.PublishMatch(expectedMatch)

	select {
	case received := <-matchChan:
		if received.ID != expectedMatch.ID {
			t.Errorf("ekspektasi ID %s, dapat %s", expectedMatch.ID, received.ID)
		}
		if received.HomeScore != 2 || received.AwayScore != 1 {
			t.Errorf("skor tidak sesuai: %v", received)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout menunggu match event")
	}
}

func TestBroker_SubscribeGlobalMatch(t *testing.T) {
	broker := NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	globalChan := broker.SubscribeMatch(ctx, "")

	match := &domain.Match{
		ID:        "match-99",
		HomeScore: 1,
		AwayScore: 0,
	}

	broker.PublishMatch(match)

	select {
	case received := <-globalChan:
		if received.ID != "match-99" {
			t.Errorf("ekspektasi match-99, dapat %s", received.ID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout menunggu update di global subscriber")
	}
}

func TestBroker_SubscribeAndPublishEvent(t *testing.T) {
	broker := NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan := broker.SubscribeEvent(ctx, "match-1")

	expectedEvent := &domain.MatchEvent{
		ID:      "event-1",
		MatchID: "match-1",
		Player:  "Bukayo Saka",
		Minute:  25,
		Type:    domain.EventTypeGoal,
	}

	broker.PublishEvent(expectedEvent)

	select {
	case received := <-eventChan:
		if received.Player != "Bukayo Saka" {
			t.Errorf("ekspektasi player Bukayo Saka, dapat %s", received.Player)
		}
		if received.Type != domain.EventTypeGoal {
			t.Errorf("ekspektasi GOAL, dapat %s", received.Type)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout menunggu event")
	}
}

func TestBroker_UnsubscribeOnContextCancel(t *testing.T) {
	broker := NewBroker()
	ctx, cancel := context.WithCancel(context.Background())

	matchChan := broker.SubscribeMatch(ctx, "match-1")

	cancel()

	time.Sleep(50 * time.Millisecond)

	broker.mu.RLock()
	subsLen := len(broker.matchSubs["match-1"])
	broker.mu.RUnlock()

	if subsLen != 0 {
		t.Errorf("subscriber harusnya sudah dihapus dari map, tersisa: %d", subsLen)
	}

	_, ok := <-matchChan
	if ok {
		t.Error("channel harusnya tertutup setelah unsubscribe")
	}
}

func TestBroker_ConcurrentSafety(t *testing.T) {
	broker := NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch := broker.SubscribeMatch(ctx, "match-1")
			select {
			case <-ch:
			case <-time.After(2 * time.Second):
			}
		}()
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(score int) {
			defer wg.Done()
			broker.PublishMatch(&domain.Match{
				ID:        "match-1",
				HomeScore: score,
			})
		}(i)
	}

	wg.Wait()
}
