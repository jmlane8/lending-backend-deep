# OB-001 Business Objectives & Boundaries for Loan Submission 

## O1: Revenue Growth
### Objective
30,000 loan applications sourced through external partner channels within 12 months of launch (the threshold that makes the investment worthwhile). 

### Boundaries
- In Scope: Integrations from multiple parties outside our organization.
- Out of Scope: Direct consumer UX. No login, no end-user account model.
- System edges are partner systems, not end users.


## O2: Partner self-sufficiency
### Objective
Read docs, test sandbox, understand how-to-integrate in 5 days with minimal hand-holding. Drives the items beyond the API itself: docs, sandbox, …

### Boundaries
- Do not rely on staff-assisted onboarding for system usage. Provide comprehensive self-service documentation instead.
- The trust gate: initial credential provisioning is performed after the legal agreement is signed. 
- Partners can rotate and regenerate credentials whenever needed.

## O3: Fast API performance
### Objective
A loan submission must be acknowledged within 2 seconds. Any delay or redirect disrupts the partner's user experience. 

### Boundaries
- The platform must adopt an asynchronous pattern for long-running operations (e.g., third party API calls)

## O4: Operational Scale
### Objective
Within 12 months, acquire 20 active partners without increasing operational headcount, as the current revenue may not yet support additional staffing costs.
