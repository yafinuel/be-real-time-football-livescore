package service_test

import (
	"context"
	"testing"
	"time"

	"football-live-score/internal/domain"
	"football-live-score/internal/service"
	"football-live-score/tests/testhelper"
)

func TestMatchService_ListMatches_Success(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	matches, err := svc.ListMatches(context.Background())
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if len(matches) == 0 {
		t.Fatal("ekspektasi setidaknya 1 pertandingan, got 0")
	}

	// Pastikan relasi tim ter-hydrate
	for _, m := range matches {
		if m.HomeTeam == nil {
			t.Errorf("HomeTeam seharusnya ter-hydrate untuk match ID: %s", m.ID)
		}
		if m.AwayTeam == nil {
			t.Errorf("AwayTeam seharusnya ter-hydrate untuk match ID: %s", m.ID)
		}
	}
}

func TestMatchService_ListMatches_RepoFail(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	matchRepo.ShouldFail = true
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	_, err := svc.ListMatches(context.Background())
	if err == nil {
		t.Fatal("ekspektasi error ketika repository gagal")
	}
}

func TestMatchService_GetMatchByID_Success(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	match, err := svc.GetMatchByID(context.Background(), "match-1")
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if match.ID != "match-1" {
		t.Errorf("ekspektasi ID match-1, got: %s", match.ID)
	}
	if match.HomeTeam == nil || match.AwayTeam == nil {
		t.Error("HomeTeam dan AwayTeam seharusnya ter-hydrate")
	}
}

func TestMatchService_GetMatchByID_NotFound(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	_, err := svc.GetMatchByID(context.Background(), "match-tidak-ada")
	if err == nil {
		t.Fatal("ekspektasi error ketika match tidak ditemukan")
	}
}

func TestMatchService_UpdateScore_Success(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	match, err := svc.UpdateScore(context.Background(), "match-1", 3, 2)
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if match.HomeScore != 3 || match.AwayScore != 2 {
		t.Errorf("skor tidak sesuai: home=%d, away=%d", match.HomeScore, match.AwayScore)
	}
}

func TestMatchService_UpdateScore_NegativeScore(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	_, err := svc.UpdateScore(context.Background(), "match-1", -1, 0)
	if err == nil {
		t.Fatal("ekspektasi error ketika skor negatif")
	}
}

func TestMatchService_UpdateScore_MatchNotFound(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	_, err := svc.UpdateScore(context.Background(), "match-tidak-ada", 1, 0)
	if err == nil {
		t.Fatal("ekspektasi error ketika match tidak ditemukan")
	}
}

func TestMatchService_UpdateStatus_Success(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	minute := 45
	match, err := svc.UpdateStatus(context.Background(), "match-1", domain.MatchStatusHalftime, &minute)
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if match.Status != domain.MatchStatusHalftime {
		t.Errorf("ekspektasi status HALFTIME, got: %s", match.Status)
	}
	if match.Minute != 45 {
		t.Errorf("ekspektasi menit 45, got: %d", match.Minute)
	}
}

func TestMatchService_UpdateStatus_CannotChangeFinished(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	// Set match menjadi FINISHED
	matchRepo.Matches["match-1"].Status = domain.MatchStatusFinished
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	_, err := svc.UpdateStatus(context.Background(), "match-1", domain.MatchStatusLive, nil)
	if err == nil {
		t.Fatal("ekspektasi error ketika mencoba mengubah status pertandingan yang sudah FINISHED")
	}
}

func TestMatchService_GetLiveMatches(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	// Match-1 sudah berstatus LIVE
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	liveMatches, err := svc.GetLiveMatches(context.Background())
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	for _, m := range liveMatches {
		if m.Status != domain.MatchStatusLive {
			t.Errorf("semua match dari GetLiveMatches seharusnya LIVE, got: %s", m.Status)
		}
	}
}

func TestMatchService_CreateMatch_Success(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	startTime := time.Now().Add(24 * time.Hour)
	match, err := svc.CreateMatch(context.Background(), "team-1", "team-2", startTime)
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if match == nil {
		t.Fatal("ekspektasi match tidak nil")
	}
	if match.HomeTeamID != "team-1" || match.AwayTeamID != "team-2" {
		t.Errorf("ID tim tidak cocok: home=%s, away=%s", match.HomeTeamID, match.AwayTeamID)
	}
	if match.Status != domain.MatchStatusNotStarted {
		t.Errorf("ekspektasi status NOT_STARTED, got: %s", match.Status)
	}
	if match.HomeScore != 0 || match.AwayScore != 0 {
		t.Errorf("ekspektasi skor awal 0-0, got: %d-%d", match.HomeScore, match.AwayScore)
	}
	if match.HomeTeam == nil || match.AwayTeam == nil {
		t.Error("HomeTeam dan AwayTeam seharusnya terisi/ter-hydrate")
	}
}

func TestMatchService_CreateMatch_SameTeam(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	_, err := svc.CreateMatch(context.Background(), "team-1", "team-1", time.Now())
	if err == nil {
		t.Fatal("ekspektasi error ketika tim kandang dan tandang sama")
	}
}

func TestMatchService_CreateMatch_TeamNotFound(t *testing.T) {
	matchRepo := testhelper.NewMockMatchRepository()
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewMatchService(matchRepo, teamRepo)

	// Tim kandang tidak ada
	_, err := svc.CreateMatch(context.Background(), "team-tidak-ada", "team-2", time.Now())
	if err == nil {
		t.Fatal("ekspektasi error ketika tim kandang tidak ditemukan")
	}

	// Tim tandang tidak ada
	_, err = svc.CreateMatch(context.Background(), "team-1", "team-tidak-ada", time.Now())
	if err == nil {
		t.Fatal("ekspektasi error ketika tim tandang tidak ditemukan")
	}
}
