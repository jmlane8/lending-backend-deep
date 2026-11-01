package loanapplication

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"lending-backend-deep/internal/rest/middleware"
	"lending-backend-deep/internal/rest/response"
	pgstore "lending-backend-deep/internal/store/pg"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const maxBodyBytes = 16 << 20

// SubmitRequest is the JSON request body for POST /loan-applications/v1.  Sections
// stay raw and fields are not typed here: a field problem is reported after
// acknowledgment, never as a 400 (FR-012, R-002).
type SubmitRequest struct {
	Customer   json.RawMessage `json:"customer"`
	Collateral json.RawMessage `json:"collateral"`
	Loan       json.RawMessage `json:"loan"`
}

// SubmitResponse is the JSON response body for a recorded application.
type SubmitResponse struct {
	ApplicationReference string `json:"application-reference"`
	SubmittedAt          string `json:"submitted-at"`
}

// Submit returns the HTTP handler for POST /loan-applications/v1.  It records the
// exact body and acknowledges it; request bodies are never logged.
//
// HTTP status codes:
//   - 202 — application recorded; it has not been validated yet
//   - 400 — body is not a JSON object
//   - 401 — missing or wrong api-key (enforced by upstream middleware)
//   - 403 — partner agreement is not active (enforced by upstream middleware)
//   - 413 — body is over 16 MiB
//   - 500 — any unhandled storage error
func Submit(logger zerolog.Logger, store *pgstore.LoanApplicationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partnerID, ok := middleware.PartnerIDFrom(r.Context())
		if !ok {
			writeInternalError(w)
			return
		}

		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
		if err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				response.WriteErrorWithEnvelope(w, http.StatusRequestEntityTooLarge,
					"REQUEST_BODY_TOO_LARGE", "request body is over 16 MiB")
				return
			}
			writeInvalidBody(w)
			return
		}

		// Unmarshal accepts JSON null for a struct, so check for an object first.
		var req SubmitRequest
		if !isJSONObject(body) || json.Unmarshal(body, &req) != nil {
			writeInvalidBody(w)
			return
		}

		id, err := uuid.NewV7()
		if err != nil {
			logger.Error().Err(err).Msg("failed to generate application id")
			writeInternalError(w)
			return
		}
		reference := "LA-" + id.String()
		submittedAt := time.Now().UTC()

		if err := store.Create(r.Context(), pgstore.CreateApplicationParams{
			ID:          id,
			Reference:   reference,
			PartnerID:   partnerID,
			SubmittedAt: submittedAt,
			Body:        body,
		}); err != nil {
			logger.Error().Err(err).Msg("failed to record loan application")
			writeInternalError(w)
			return
		}

		response.WriteJSONWithEnvelope(w, http.StatusAccepted, SubmitResponse{
			ApplicationReference: reference,
			SubmittedAt:          submittedAt.Format(time.RFC3339),
		})
	}
}

func isJSONObject(body []byte) bool {
	return bytes.HasPrefix(bytes.TrimLeft(body, " \t\r\n"), []byte("{"))
}

func writeInvalidBody(w http.ResponseWriter) {
	response.WriteErrorWithEnvelope(w, http.StatusBadRequest,
		"INVALID_REQUEST_BODY", "request body must be a JSON object")
}

func writeInternalError(w http.ResponseWriter) {
	response.WriteErrorWithEnvelope(w, http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR", "internal server error")
}
