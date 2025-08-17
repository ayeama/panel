package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

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
			// slog.String("host", r.Host),
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
