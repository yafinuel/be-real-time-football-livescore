package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockPinger untuk menyimulasikan respon database tanpa perlu DB asli saat unit testing
type mockPinger struct {
	shouldFail bool
}

func (m *mockPinger) Ping(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("koneksi database terputus")
	}
	return nil
}

func TestHealthHandler_Connected(t *testing.T) {
	mockDB := &mockPinger{shouldFail: false}
	handler := HealthHandler(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/healthy", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ekspektasi status 200 OK, didapat: %d", rec.Code)
	}

	var res HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("gagal decode response JSON: %v", err)
	}

	if res.Status != "healthy" || res.Database != "connected" {
		t.Errorf("respon tidak sesuai: %+v", res)
	}
}

func TestHealthHandler_Disconnected(t *testing.T) {
	mockDB := &mockPinger{shouldFail: true}
	handler := HealthHandler(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/healthy", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("ekspektasi status 503 Service Unavailable, didapat: %d", rec.Code)
	}

	var res HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("gagal decode response JSON: %v", err)
	}

	if res.Status != "unhealthy" || res.Database != "disconnected" {
		t.Errorf("respon tidak sesuai: %+v", res)
	}
}