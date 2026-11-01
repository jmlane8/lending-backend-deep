package middleware

import (
	"context"
	"crypto/sha256"
	"errors"
	"net/http"

	"lending-backend-deep/internal/rest/response"
	pgstore "lending-backend-deep/internal/store/pg"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const apiKeyHeader = "api-key"

type partnerIDKey struct{}

// PartnerAuth returns a chi-compatible middleware that authenticates the
// api-key header against the stored key hashes.
//
// A missing or unknown key short-circuits with HTTP 401, and a valid key whose
// partner agreement is not active with HTTP 403.  Otherwise the partner ID is
// placed in the request context (see PartnerIDFrom) and the request is logged
// against that partner, so every authenticated request is attributable.
//
// Usage:
//
//	router.Use(middleware.PartnerAuth(logger, partnerStore))
func PartnerAuth(logger zerolog.Logger, store *pgstore.PartnerStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get(apiKeyHeader)
			if key == "" {
				writeUnauthorized(w)
				return
			}

			hash := sha256.Sum256([]byte(key))
			credential, err := store.FindByKeyHash(r.Context(), hash[:])
			if errors.Is(err, pgstore.ErrNotFound) {
				writeUnauthorized(w)
				return
			}
			if err != nil {
				logger.Error().Err(err).Msg("failed to look up partner credential")
				response.WriteErrorWithEnvelope(w, http.StatusInternalServerError,
					"INTERNAL_SERVER_ERROR", "internal server error")
				return
			}

			logger.Info().
				Str("partner_id", credential.PartnerID.String()).
				Str("agreement_status", credential.AgreementStatus).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("partner request")

			if credential.AgreementStatus != "ACTIVE" {
				response.WriteErrorWithEnvelope(w, http.StatusForbidden,
					"PARTNER_AGREEMENT_INACTIVE", "forbidden: partner agreement is not active")
				return
			}

			ctx := context.WithValue(r.Context(), partnerIDKey{}, credential.PartnerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// PartnerIDFrom returns the authenticated partner ID placed in ctx by PartnerAuth.
func PartnerIDFrom(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(partnerIDKey{}).(uuid.UUID)
	return id, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	response.WriteErrorWithEnvelope(w, http.StatusUnauthorized,
		"UNAUTHORIZED", "unauthorized: invalid or missing api-key")
}
