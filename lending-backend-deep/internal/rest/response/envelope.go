package response

import (
	"encoding/json"
	"net/http"
)

// ResponseEnvelope is the standard JSON envelope returned by every endpoint.
//
// Example success:
//
//	{"data":{...}}
//
// Example error:
//
//	{"error-code":"UNAUTHORIZED","error-message":"unauthorized: invalid or missing api-key"}
type ResponseEnvelope struct {
	Data         any    `json:"data,omitempty"`
	ErrorCode    string `json:"error-code,omitempty"`
	ErrorMessage string `json:"error-message,omitempty"`
}

// WriteJSONWithEnvelope serialises data inside the standard envelope and writes it with
// the given HTTP status code.  The Content-Type header is set to
// "application/json".
func WriteJSONWithEnvelope(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ResponseEnvelope{
		Data: data,
	})
}

// WriteErrorWithEnvelope writes a JSON error envelope with a machine-readable
// error code and a human-readable message. Use this for 4xx and 5xx responses.
func WriteErrorWithEnvelope(w http.ResponseWriter, statusCode int, errorCode, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ResponseEnvelope{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	})
}
