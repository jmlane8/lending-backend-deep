// =============================================================================
// k6 test script — ONE test case (TC-NNN-NNN) from a Test Scenario (TS-NNN).
// This is the canonical reference structure. Every generated case file follows
// THIS section order so all scripts look the same:
//   1. Imports (incl. ./config.js)   2. File header   3. Payload helper
//   4. Custom metrics (only if needed)   5. options   6. default function
//
// One file per test case. File name: k6-<TC-id>-<short-kebab-title>.js
//   e.g. k6-TC-001-002-non-contracted-party-submission-is-rejected.js
// Shared connection settings come from ./config.js (the one allowed local import).
// A case file must NOT import another case file.
// Run from inside the scenario folder: k6 run k6-TC-NNN-NNN-<short-kebab-title>.js
//
// Keep payloads/values generic to the scenario you are converting —
// do NOT copy the placeholders here verbatim.
// =============================================================================

// 1. Imports — k6 built-ins, public jslib (optional), and the shared config.
import http from 'k6/http';
import { check, group } from 'k6';
// import { Counter } from 'k6/metrics'; // only for cross-request (concurrency) cases
import { ENDPOINT, API_KEY, headers } from './config.js';

// 2. File header — fill in for each generated file:
//   Source : test-scenarios/<TS-file>.md  —  TC-NNN-NNN <title>
//   Manual preconditions this script cannot create (seed data, pre-existing
//     state, services forced down): …
//   Open gaps (see TODO(gap) markers): …

// 3. Payload — the canonical valid body this case starts from (test data stays per file).
// Shape (field names/types/required) comes from the OpenAPI request schema; concrete values
// come from the scenario's reference positive case.
// TODO(gap) example — when NEITHER the scenario nor the OpenAPI defines a wire format, use a
// clearly named placeholder and flag it; do NOT fabricate a real value.
// const VALID_PHOTO = __ENV.VALID_PHOTO || 'PLACEHOLDER_VALID_PHOTO';
function validPayload() {
  return {
    // Canonical valid request body = OpenAPI request schema, filled with the scenario's values.
    // A negative case mutates one field of this; a positive case sends it as-is.
  };
}

// 4. Custom metrics — only when this case asserts a cross-request rule (concurrency).
// const concurrentAccepted = new Counter('concurrent_accepted');

// 5. options — this single case runs once. checks==1.0 makes the run a pass/fail gate.
//    (vus/iterations/thresholds are per-case and live HERE, never in config.js.)
export const options = {
  vus: 1,
  iterations: 1,
  thresholds: {
    checks: ['rate==1.0'], // any failed check fails the whole run (CI gate)
    // Time-SLA case — bound the request duration (tag keeps it case-scoped):
    // 'http_req_duration{tc:TC-NNN-NNN}': ['p(100)<2000'],
    // Concurrency case — set vus:2, iterations:1 above and assert the aggregate rule:
    // concurrent_accepted: ['count==1'],
  },
};

// 6. default function — perform the request(s) and the named, tagged checks.
export default function () {
  group('TC-NNN-NNN — <title>', () => {
    const body = validPayload();
    // For a negative case, mutate one field here, e.g. delete body.<field>;
    // http.<method> + endpoint come from the OpenAPI operation (POST shown as the common case).
    const res = http.post(ENDPOINT, JSON.stringify(body), {
      headers: headers(API_KEY),
      tags: { tc: 'TC-NNN-NNN', type: 'positive' }, // or 'negative'
    });
    check(res, {
      'TC-NNN-NNN: <expected result>': (r) => r.status === 200 || r.status === 202,
      // Safe JSON-field read:
      // 'TC-NNN-NNN: body returns an id': (r) => { try { return r.json('id') !== undefined; } catch (_) { return false; } },
    }, { tc: 'TC-NNN-NNN' });

    // Concurrency variant — with options.vus = 2, each VU runs this once; a VU sees only
    // its own response, so count the wins and assert count==1 via the threshold above:
    // const accepted = res.status === 200 || res.status === 202;
    // if (accepted) concurrentAccepted.add(1);
  });
}
