package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors_PreflightOptions(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := Cors(next)

	req := httptest.NewRequest(http.MethodOptions, "/query", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ekspektasi status 200, dapat %d", rec.Code)
	}

	if nextCalled {
		t.Error("preflight OPTIONS tidak boleh memanggil handler berikutnya")
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" && allowOrigin != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin tidak sesuai: %s", allowOrigin)
	}
}

func TestCors_NormalRequest(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := Cors(next)

	req := httptest.NewRequest(http.MethodPost, "/query", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Error("handler berikutnya harus dipanggil untuk request non-OPTIONS")
	}

	allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
	if allowMethods == "" {
		t.Error("Access-Control-Allow-Methods harus terisi")
	}
}
