# Contract: Rejection Webhook

The platform calls this on the partner's system when a submission fails customer data validation (FR-013). Only rejections are sent. Use the [validation result lookup](loan-applications.openapi.yaml) to see passed applications, or if a webhook never arrives.

## Request

`POST <partner notification_url>`

| Header | Value |
|---|---|
| `Content-Type` | `application/json` |
| `x-event-id` | Same value as `event-id` in the body |
| `x-timestamp` | Unix time in seconds when this attempt was signed |
| `x-signature` | `sha256=` followed by the lowercase hex HMAC-SHA256 of `<x-timestamp>.<raw body>`, using the partner's notification secret |

```json
{
  "event-id": "019a3b1d-22aa-7c31-8f10-5d9e3a7b1c02",
  "event-type": "loan-application.rejected",
  "application-reference": "LA-019a3b1c-7d2e-7f10-9c44-2b6e1f0a8d11",
  "rejected-at": "2026-10-31T04:15:09Z",
  "reasons": [
    {
      "field": "customer.phone-number",
      "code": "INVALID",
      "message": "phone-number must be an optional + followed by 8 to 15 digits"
    }
  ]
}
```

`reasons` uses the `RejectionReason` schema and never contains submitted values.

## Partner responsibilities

- Verify the signature with a constant-time comparison before trusting the body.
- Reject requests with an old `x-timestamp`. Five minutes is a sensible limit.
- De-duplicate on `event-id`. Delivery is at-least-once, so the same event can arrive more than once.
- Return any `2xx` within 5 seconds.

## Delivery behaviour

Any other status, a timeout, or a connection failure counts as a failed attempt. Retries happen 1 minute, 5 minutes, 30 minutes, 2 hours, 6 hours and 12 hours after the previous attempt, for 7 attempts in total. After that, the platform stops, and the result stays available through the lookup.
