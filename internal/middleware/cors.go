package middleware

import (
	"net/http"
)

// Cors adalah middleware untuk mengizinkan Cross-Origin Resource Sharing (CORS)
// agar aplikasi Frontend (React, Vue, Flutter, dll) bisa mengakses API dari port manapun.
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Tangkap origin client (misal http://localhost:3000 atau http://localhost:5173)
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 2. Tangani Preflight Request (OPTIONS)
		// Browser selalu mengirim OPTIONS terlebih dahulu sebelum POST/GET lintas port
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 3. Teruskan request ke handler berikutnya (GraphQL server)
		next.ServeHTTP(w, r)
	})
}
