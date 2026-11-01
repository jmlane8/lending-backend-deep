import { check, group } from 'k6';

// Source  : test-scenarios/TS-001-submit-loan-application-customer-data.md - TC-001-048 Customer born 29 February reaches 18 on 1 March of a non-leap year
// OpenAPI : http://localhost:8080/openapi.json (3.0.3)
// Manual preconditions: the service is running, with the sandbox partners seeded and BUSINESS_TIMEZONE=UTC.
// Forced mismatch: this case calls validation.Validate with a fixed date, which no API operation allows.

export const options = {
  vus: 1,
  iterations: 1,
  thresholds: {
    checks: ['rate==1.0'], // any failed check fails the run
  },
};

export default function () {
  group("TC-001-048 — Customer born 29 February reaches 18 on 1 March of a non-leap year", () => {
    // TODO(gap): the API always uses the current date, so date-of-birth `2008-02-29`, Submission date `2026-03-01` can't be set from k6.
    // This fails on purpose until the case runs elsewhere. Expected: No reasons are returned.
    check(null, {
      "TC-001-048: not expressible in k6, run against validation.Validate": () => false,
    }, { tc: 'TC-001-048', type: 'positive' });
  });
}
