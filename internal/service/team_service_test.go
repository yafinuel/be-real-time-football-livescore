package service_test

import (
	"context"
	"testing"

	"football-live-score/internal/service"
	"football-live-score/tests/testhelper"
)

func TestTeamService_ListTeams_Success(t *testing.T) {
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewTeamService(teamRepo)

	teams, err := svc.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if len(teams) == 0 {
		t.Fatal("ekspektasi setidaknya 1 tim, got 0")
	}
}

func TestTeamService_ListTeams_RepoFail(t *testing.T) {
	teamRepo := testhelper.NewMockTeamRepository()
	teamRepo.ShouldFail = true
	svc := service.NewTeamService(teamRepo)

	_, err := svc.ListTeams(context.Background())
	if err == nil {
		t.Fatal("ekspektasi error ketika repository gagal")
	}
}

func TestTeamService_GetTeamByID_Success(t *testing.T) {
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewTeamService(teamRepo)

	team, err := svc.GetTeamByID(context.Background(), "team-1")
	if err != nil {
		t.Fatalf("ekspektasi sukses, got error: %v", err)
	}

	if team.ID != "team-1" {
		t.Errorf("ekspektasi ID team-1, got: %s", team.ID)
	}
	if team.Name != "Arsenal" {
		t.Errorf("ekspektasi nama Arsenal, got: %s", team.Name)
	}
}

func TestTeamService_GetTeamByID_NotFound(t *testing.T) {
	teamRepo := testhelper.NewMockTeamRepository()
	svc := service.NewTeamService(teamRepo)

	_, err := svc.GetTeamByID(context.Background(), "team-tidak-ada")
	if err == nil {
		t.Fatal("ekspektasi error ketika tim tidak ditemukan")
	}
}
