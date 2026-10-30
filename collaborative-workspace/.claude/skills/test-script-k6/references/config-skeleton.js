// =============================================================================
// config.js — shared connection settings for ONE test scenario folder.
// Imported by every k6-TC-*.js case file in the same folder. Edit these for your
// environment, or override via env vars at run time. Uniform name: always "config.js".
//
// Keep this MINIMAL — connection only. Per-case values (vus, iterations, thresholds,
// payload, checks) belong in each case file, NEVER here.
// Values below come from the API's OpenAPI document, not from guessing.
// =============================================================================

// BASE_URL — from the OpenAPI `servers`.
export const BASE_URL = __ENV.BASE_URL || 'https://api.test.example.internal/v1';

// Endpoint path(s) — from the OpenAPI `paths`. Keep here only when every case hits the same one.
export const ENDPOINT = `${BASE_URL}/<resource>`;

// Credentials — one export per distinct key the scenario uses.
export const API_KEY = __ENV.API_KEY || '<api-key>';
// e.g. when a scenario distinguishes contracted vs non-contracted callers:
// export const CONTRACTED_PARTNER_KEY = __ENV.CONTRACTED_PARTNER_KEY || '<contracted-partner-key>';
// export const NON_CONTRACTED_KEY = __ENV.NON_CONTRACTED_KEY || '<non-contracted-key>';

// Small shared helper — auth shape from the OpenAPI `securitySchemes`. Match your contract:
// bearer -> `Bearer ${apiKey}`; apiKey-in-header -> the named header instead of Authorization.
export function headers(apiKey) {
  return { 'Content-Type': 'application/json', Authorization: `Token ${apiKey}` };
}
