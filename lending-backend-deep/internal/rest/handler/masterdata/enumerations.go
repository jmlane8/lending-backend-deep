package masterdata

import (
	"net/http"

	"lending-backend-deep/internal/rest/response"
	pgstore "lending-backend-deep/internal/store/pg"

	"github.com/rs/zerolog"
)

// EnumerationsResponse is the JSON response body for GET /master-data/v1/enumerations.
type EnumerationsResponse struct {
	Gender        []string `json:"gender"`
	MaritalStatus []string `json:"marital-status"`
	Occupation    []string `json:"occupation"`
}

// Enumerations returns the HTTP handler for GET /master-data/v1/enumerations.
//
// HTTP status codes:
//   - 200 — accepted values by category
//   - 401 — missing or wrong api-key (enforced by upstream middleware)
//   - 403 — partner agreement is not active (enforced by upstream middleware)
//   - 500 — any unhandled storage error
func Enumerations(logger zerolog.Logger, store *pgstore.MasterDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, err := store.ListActive(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("failed to list master data")
			response.WriteErrorWithEnvelope(w, http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR", "internal server error")
			return
		}

		response.WriteJSONWithEnvelope(w, http.StatusOK, EnumerationsResponse{
			Gender:        values["GENDER"],
			MaritalStatus: values["MARITAL_STATUS"],
			Occupation:    values["OCCUPATION"],
		})
	}
}
