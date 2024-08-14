package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			attr := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
			}

			// Call the next handler in the chain
			next.ServeHTTP(ww, r)

			// Log the time taken to process the request
			latency := time.Since(start)
			attr = append(attr,
				slog.Duration("latency", latency),
				slog.Int("status", ww.Status()),
			)

			logger.Info("Request processed",
				"attr", attr,
			)
		})
	}
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		token := r.Header.Get("Authorization")

		// Check if the token matches the static token
		if token != s.StaticToken {
			s.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// If the token matches, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
