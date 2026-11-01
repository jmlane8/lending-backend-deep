package middleware

import (
	"net/http"

	"lending-backend-deep/internal/rest/response"

	"github.com/rs/zerolog"
)

// PanicRecovery returns a middleware that catches any panic propagating up the
// handler chain, logs it as an ERROR, and responds with HTTP 500 JSON.
//
// This middleware should be registered early in the chain so that it wraps all
// downstream handlers.
func PanicRecovery(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error().
						Interface("panic", rec).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Msg("panic recovered")
					response.WriteErrorWithEnvelope(w, http.StatusInternalServerError,
						"INTERNAL_SERVER_ERROR", "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
