// Shared connection settings for TS-001-submit-loan-application-customer-data.
// Values come from the API's OpenAPI document (http://localhost:8080/openapi.json).
// Connection only: per-case values live in each case file.

// BASE_URL: from the OpenAPI `servers`.
export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Endpoints: from the OpenAPI `paths`.
export const ENDPOINT = `${BASE_URL}/loan-applications/v1`;
export function validationResultUrl(reference) {
  return `${ENDPOINT}/${reference}/validation-result`;
}

// Credential: the active sandbox partner named in TS-001.
export const API_KEY = __ENV.API_KEY || 'sandbox-partner-key';

// Auth: `ApiKeyAuth` in the OpenAPI `securitySchemes` is an apiKey in the `api-key` header.
export function headers(apiKey) {
  return { 'Content-Type': 'application/json', 'api-key': apiKey };
}
