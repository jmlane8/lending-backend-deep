import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { ENDPOINT, API_KEY, headers, validationResultUrl } from './config.js';

// Source  : test-scenarios/TS-001-submit-loan-application-customer-data.md - TC-001-033 Marital status DIVORCED is accepted
// OpenAPI : http://localhost:8080/openapi.json (3.0.3)
// Manual preconditions: the service is running, with the sandbox partners seeded and BUSINESS_TIMEZONE=UTC.

function validPayload() {
  return {
    "collateral": {},
    "customer": {
      "date-of-birth": "1990-02-11",
      "domicile-city": "Springfield",
      "domicile-district": "Central",
      "domicile-province": "North",
      "domicile-street": "12 Example Street",
      "email": "alex@example.com",
      "full-legal-name": "Alex Morgan",
      "gender": "FEMALE",
      "government-id-expiry-date": "2099-05-14",
      "government-id-photo": "/9j/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIAAQABAMBIgACEQEDEQH/xAGiAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+gEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoLEQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/APn+iiigD//Z",
      "marital-status": "SINGLE",
      "monthly-income": 4500.5,
      "national-id-number": "3201012345678901",
      "number-of-dependents": 1,
      "occupation": "EMPLOYEE",
      "phone-number": "+15551234567",
      "place-of-birth": "Springfield",
      "selfie-photo": "/9j/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIAAQABAMBIgACEQEDEQH/xAGiAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+gEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoLEQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/APn+iiigD//Z",
      "work-duration-years": 6
    },
    "loan": {}
  };
}

export const options = {
  vus: 1,
  iterations: 1,
  thresholds: {
    checks: ['rate==1.0'], // any failed check fails the run
  },
};

// Reads the lookup until the status is no longer PENDING, for at most 1 minute (TS-001 lookup timing).
function waitForResult(reference, tags) {
  const deadline = Date.now() + 60000;
  let res = http.get(validationResultUrl(reference), { headers: headers(API_KEY), tags });
  while (res.status === 200 && res.json('data.status') === 'PENDING' && Date.now() < deadline) {
    sleep(0.5);
    res = http.get(validationResultUrl(reference), { headers: headers(API_KEY), tags });
  }
  return res;
}

export default function () {
  group("TC-001-033 — Marital status DIVORCED is accepted", () => {
    const tags = { tc: 'TC-001-033', type: 'positive' };
    const body = validPayload();
    body.customer["marital-status"] = "DIVORCED";
    const res = http.post(ENDPOINT, JSON.stringify(body), { headers: headers(API_KEY), tags });
    // A failed POST has no reference; the lookup then returns an error and the check below fails.
    const reference = res.status === 202 ? res.json('data.application-reference') : 'not-acknowledged';
    const result = waitForResult(reference, tags);
    check(result, {
      "TC-001-033: lookup returns PASSED with no reasons": (r) => { try { return r.status === 200 && r.json('data.status') === 'PASSED' && (r.json('data.reasons') === null || r.json('data.reasons') === undefined); } catch (_) { return false; } },
    }, tags);
  });
}
