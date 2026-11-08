// Package integration berisi integration test untuk semua endpoint GraphQL API.
// Test ini berjalan terhadap server HTTP asli menggunakan httptest.
// Menggunakan mock repository sehingga tidak memerlukan koneksi database asli.
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"football-live-score/internal/graph"
	"football-live-score/internal/graph/model"
	"football-live-score/internal/pubsub"
	"football-live-score/internal/service"
	"football-live-score/tests/testhelper"

	"github.com/99designs/gqlgen/graphql/handler"
)

// graphQLResponse adalah struktur umum respons GraphQL
type graphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// newTestResolver membuat resolver lengkap dengan mock repo dan broker
func newTestResolver() *graph.Resolver {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	eventRepo := testhelper.NewMockEventRepository()

	matchService := service.NewMatchService(matchRepo, teamRepo)
	eventService := service.NewEventService(eventRepo, matchRepo)
	teamService := service.NewTeamService(teamRepo)
	broker := pubsub.NewBroker()

	return &graph.Resolver{
		MatchService: matchService,
		EventService: eventService,
		TeamService:  teamService,
		Broker:       broker,
	}
}

// newTestServer membuat HTTP test server GraphQL dengan mock repository
func newTestServer() *httptest.Server {
	resolver := newTestResolver()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	return httptest.NewServer(srv)
}

// doGraphQL mengirim request GraphQL ke server dan mengembalikan respons
func doGraphQL(t *testing.T, server *httptest.Server, query string) graphQLResponse {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	var result graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("gagal decode respons: %v", err)
	}
	return result
}

// --- QUERY TESTS ---

func TestQuery_Matches(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `query {
		matches {
			id
			status
			minute
			score { home away }
			homeTeam { id name shortName }
			awayTeam { id name shortName }
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("query matches gagal: %v", resp.Errors[0].Message)
	}

	matches, ok := resp.Data["matches"].([]interface{})
	if !ok {
		t.Fatal("respons matches bukan array")
	}
	if len(matches) == 0 {
		t.Error("ekspektasi setidaknya 1 pertandingan")
	}
}

func TestQuery_Match_Found(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `query {
		match(id: "match-1") {
			id
			status
			minute
			score { home away }
			homeTeam { name }
			awayTeam { name }
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("query match(id) gagal: %v", resp.Errors[0].Message)
	}

	match := resp.Data["match"]
	if match == nil {
		t.Fatal("ekspektasi data match, got null")
	}
}

func TestQuery_Match_NotFound(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `query {
		match(id: "match-tidak-ada") {
			id
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("query tidak boleh error untuk match tidak ditemukan: %v", resp.Errors[0].Message)
	}

	// Sesuai PRD Q-02: mengembalikan null jika tidak ditemukan
	if resp.Data["match"] != nil {
		t.Error("ekspektasi null untuk match yang tidak ditemukan")
	}
}

func TestQuery_LiveMatches(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `query {
		liveMatches {
			id
			status
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("query liveMatches gagal: %v", resp.Errors[0].Message)
	}

	liveMatches, ok := resp.Data["liveMatches"].([]interface{})
	if !ok {
		t.Fatal("respons liveMatches bukan array")
	}

	// Semua hasil harus berstatus LIVE
	for _, m := range liveMatches {
		match := m.(map[string]interface{})
		if match["status"] != "LIVE" {
			t.Errorf("semua liveMatches harus berstatus LIVE, got: %v", match["status"])
		}
	}
}

func TestQuery_Teams(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `query {
		teams {
			id
			name
			shortName
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("query teams gagal: %v", resp.Errors[0].Message)
	}

	teams, ok := resp.Data["teams"].([]interface{})
	if !ok {
		t.Fatal("respons teams bukan array")
	}
	if len(teams) == 0 {
		t.Error("ekspektasi setidaknya 1 tim")
	}
}

// --- MUTATION TESTS ---

func TestMutation_UpdateScore_Success(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `mutation {
		updateScore(input: {
			matchId: "match-1"
			homeScore: 3
			awayScore: 1
		}) {
			id
			score { home away }
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("mutation updateScore gagal: %v", resp.Errors[0].Message)
	}

	result, ok := resp.Data["updateScore"].(map[string]interface{})
	if !ok {
		t.Fatal("respons updateScore bukan object")
	}
	score := result["score"].(map[string]interface{})
	if score["home"].(float64) != 3 || score["away"].(float64) != 1 {
		t.Errorf("skor tidak sesuai: %v", score)
	}
}

func TestMutation_UpdateScore_NegativeScore(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `mutation {
		updateScore(input: {
			matchId: "match-1"
			homeScore: -1
			awayScore: 0
		}) {
			id
		}
	}`)

	if len(resp.Errors) == 0 {
		t.Fatal("ekspektasi error untuk skor negatif")
	}
}

func TestMutation_UpdateMatchStatus_Success(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `mutation {
		updateMatchStatus(matchId: "match-1", status: HALFTIME, minute: 45) {
			id
			status
			minute
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("mutation updateMatchStatus gagal: %v", resp.Errors[0].Message)
	}

	result, ok := resp.Data["updateMatchStatus"].(map[string]interface{})
	if !ok {
		t.Fatal("respons updateMatchStatus bukan object")
	}
	if result["status"] != "HALFTIME" {
		t.Errorf("ekspektasi status HALFTIME, got: %v", result["status"])
	}
	if result["minute"].(float64) != 45 {
		t.Errorf("ekspektasi menit 45, got: %v", result["minute"])
	}
}

func TestMutation_AddMatchEvent_Success(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `mutation {
		addMatchEvent(matchId: "match-1", input: {
			minute: 23
			type: GOAL
			teamId: "team-1"
			player: "Bukayo Saka"
			detail: "Assist oleh Odegaard"
		}) {
			id
			minute
			type
			player
			detail
		}
	}`)

	if len(resp.Errors) > 0 {
		t.Fatalf("mutation addMatchEvent gagal: %v", resp.Errors[0].Message)
	}

	result, ok := resp.Data["addMatchEvent"].(map[string]interface{})
	if !ok {
		t.Fatal("respons addMatchEvent bukan object")
	}
	if result["player"] != "Bukayo Saka" {
		t.Errorf("ekspektasi player Bukayo Saka, got: %v", result["player"])
	}
	if result["type"] != "GOAL" {
		t.Errorf("ekspektasi type GOAL, got: %v", result["type"])
	}
}

func TestMutation_AddMatchEvent_InvalidMinute(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `mutation {
		addMatchEvent(matchId: "match-1", input: {
			minute: 999
			type: GOAL
			teamId: "team-1"
			player: "Bukayo Saka"
		}) {
			id
		}
	}`)

	if len(resp.Errors) == 0 {
		t.Fatal("ekspektasi error untuk menit tidak valid")
	}
}

func TestMutation_AddMatchEvent_InvalidTeam(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp := doGraphQL(t, server, `mutation {
		addMatchEvent(matchId: "match-1", input: {
			minute: 23
			type: GOAL
			teamId: "team-999"
			player: "Unknown Player"
		}) {
			id
		}
	}`)

	if len(resp.Errors) == 0 {
		t.Fatal("ekspektasi error ketika tim tidak bertanding di pertandingan ini")
	}
}

// --- SUBSCRIPTION TESTS ---

func TestSubscription_MatchUpdated(t *testing.T) {
	resolver := newTestResolver()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	matchID := "match-1"
	subChan, err := resolver.Subscription().MatchUpdated(ctx, &matchID)
	if err != nil {
		t.Fatalf("gagal subscribe matchUpdated: %v", err)
	}

	// Jalankan mutation UpdateScore
	_, err = resolver.Mutation().UpdateScore(ctx, model.UpdateScoreInput{
		MatchID:   "match-1",
		HomeScore: 3,
		AwayScore: 2,
	})
	if err != nil {
		t.Fatalf("gagal update score: %v", err)
	}

	// Pastikan subscription menerima match yang telah di-update
	select {
	case updated := <-subChan:
		if updated.ID != "match-1" {
			t.Errorf("ekspektasi ID match-1, dapat: %s", updated.ID)
		}
		if updated.Score.Home != 3 || updated.Score.Away != 2 {
			t.Errorf("skor update tidak sesuai: %v-%v", updated.Score.Home, updated.Score.Away)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout menunggu matchUpdated broadcast")
	}
}

func TestSubscription_NewEventAdded(t *testing.T) {
	resolver := newTestResolver()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	matchID := "match-1"
	subChan, err := resolver.Subscription().NewEventAdded(ctx, matchID)
	if err != nil {
		t.Fatalf("gagal subscribe newEventAdded: %v", err)
	}

	// Jalankan mutation AddMatchEvent
	_, err = resolver.Mutation().AddMatchEvent(ctx, "match-1", model.AddMatchEventInput{
		Minute: 55,
		Type:   model.EventTypeGoal,
		TeamID: "team-1",
		Player: "Gabriel Martinelli",
	})
	if err != nil {
		t.Fatalf("gagal add match event: %v", err)
	}

	// Pastikan subscription menerima event baru
	select {
	case event := <-subChan:
		if event.Player != "Gabriel Martinelli" {
			t.Errorf("ekspektasi Gabriel Martinelli, dapat: %s", event.Player)
		}
		if event.Type != model.EventTypeGoal {
			t.Errorf("ekspektasi GOAL, dapat: %s", event.Type)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout menunggu newEventAdded broadcast")
	}
}
