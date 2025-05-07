package internal

import (
	"log/slog"
	"net/http"
	"time"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO add configuration

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// type responseWriter struct {
// 	http.ResponseWriter
// 	statusCode int
// }

// func Log(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

// 		responseTimeStart := time.Now().UTC()
// 		next.ServeHTTP(wrapped, r)
// 		responseTimeEnd := time.Now().UTC()

// 		slog.Info(
// 			"handled request",
// 			slog.String("host", r.Host),
// 			// slog.Int("status", wrapped.statusCode),
// 			slog.String("method", r.Method),
// 			slog.String("path", r.URL.Path),
// 			slog.String("query", r.URL.RawQuery),
// 			// slog.String("version", r.Proto),
// 			// slog.String("useragent", r.UserAgent()),
// 			slog.Int64("response_time_ms", responseTimeEnd.Sub(responseTimeStart).Milliseconds()),
// 		)
// 	})
// }

// TODO add different logging for websockets
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseTimeStart := time.Now().UTC()
		next.ServeHTTP(w, r)
		responseTimeEnd := time.Now().UTC()

		slog.Info(
			"handled request",
			slog.String("host", r.Host),
			// slog.Int("status", wrapped.statusCode),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("query", r.URL.RawQuery),
			// slog.String("version", r.Proto),
			// slog.String("useragent", r.UserAgent()),
			slog.Int64("response_time_ms", responseTimeEnd.Sub(responseTimeStart).Milliseconds()),
		)
	})
}
