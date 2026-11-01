package loanapplication

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"lending-backend-deep/internal/rest/middleware"
	"lending-backend-deep/internal/rest/response"
	pgstore "lending-backend-deep/internal/store/pg"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// ValidationResultResponse is the JSON response body for a validation result.
// ValidatedAt is null while the application is pending.
type ValidationResultResponse struct {
	ApplicationReference string          `json:"application-reference"`
	Status               string          `json:"status"`
	SubmittedAt          string          `json:"submitted-at"`
	ValidatedAt          *string         `json:"validated-at"`
	Reasons              json.RawMessage `json:"reasons,omitempty"`
}

// ValidationResult returns the HTTP handler for
// GET /loan-applications/v1/{application-reference}/validation-result.  It reports
// the validation result only, never a loan decision.
//
// HTTP status codes:
//   - 200 — current validation result
//   - 401 — missing or wrong api-key (enforced by upstream middleware)
//   - 403 — partner agreement is not active (enforced by upstream middleware)
//   - 404 — unknown reference, or one owned by another partner
//   - 500 — any unhandled storage error
func ValidationResult(logger zerolog.Logger, store *pgstore.LoanApplicationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partnerID, ok := middleware.PartnerIDFrom(r.Context())
		if !ok {
			writeInternalError(w)
			return
		}

		result, err := store.FindByReference(r.Context(), chi.URLParam(r, "application-reference"), partnerID)
		if errors.Is(err, pgstore.ErrNotFound) {
			response.WriteErrorWithEnvelope(w, http.StatusNotFound,
				"APPLICATION_NOT_FOUND", "application not found")
			return
		}
		if err != nil {
			logger.Error().Err(err).Msg("failed to read validation result")
			writeInternalError(w)
			return
		}

		resp := ValidationResultResponse{
			ApplicationReference: result.Reference,
			Status:               result.Status,
			SubmittedAt:          result.SubmittedAt.UTC().Format(time.RFC3339),
		}
		if result.ValidatedAt != nil {
			validatedAt := result.ValidatedAt.UTC().Format(time.RFC3339)
			resp.ValidatedAt = &validatedAt
		}
		if result.Status == "REJECTED" {
			resp.Reasons = result.Reasons
		}
		response.WriteJSONWithEnvelope(w, http.StatusOK, resp)
	}
}
