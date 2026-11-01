package healthcheck

import (
	"encoding/json"
	"lending-backend-deep/internal/rest/constant"
	"net/http"
)

// HealthcheckResponse is the JSON response body for GET /healthcheck.
type HealthcheckResponse struct {
	Status string `json:"status"`
}

// Healthcheck returns the HTTP handler for GET /healthcheck.
//
// HTTP status codes:
//   - 200 — service is healthy
func Healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", constant.MimeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthcheckResponse{Status: "ok"})
	}
}
