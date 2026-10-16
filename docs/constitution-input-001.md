# BUSINESS PROBLEMS
- Our company wants to expand its loan submission capability to external partners. The current system only supports internal submissions, blocking any collaboration with outside organizations.
- Our market growth has started to slow. The company wants to increase revenue and customer reach through external distribution channels without significantly expanding the marketing team or increasing advertising spend.
- External partners (digital-native companies or startups with 40,000+ users) will serve as loan submission channels for car-based collateral financing. Partners earn a commission on approved loans; the company handles the financing.
- We want scalable business growth through independent external partnerships.

---

# BUSINESS OBJECTIVES
- O1: Revenue growth: 30,000 new loan applications must be sourced through external partner channels within 12 months of launch.
- O2: Partner self-sufficiency: A new partner must be able to read documentation, test in the sandbox, and go live within 5 business days with minimal hand-holding from our team.
- O3: API performance: The API must return a submission acknowledgment within 2 seconds.
- O4: Operational scale: The platform must support at least 20 active partner integrations within 12 months of launch without increasing operational team headcount.

---

# BOUNDARIES
- Users: the direct user of the platform is the external partner system (website, mobile app, or partner backend), not the end consumer. End consumers never call our platform directly. We identify partners, not individual users.
- Authentication: API key-based machine-to-machine authentication. Keys are revocable immediately if compromised, and each key is traceable to a specific partner.

## IN SCOPE
- REST API, chosen over gRPC and GraphQL to minimize integration friction for partners and support the headcount objective.
- Asynchronous loan processing. The API accepts the submission and returns an acknowledgment within 2 seconds. The loan decision is delivered to the partner separately.
- Vehicle identification using the existing internal loan system, which responds in 10 to 200 milliseconds.
- Credit rating check via an external credit bureau. Response time is unpredictable, so this is processed asynchronously after submission.
- Underwriting and credit assessment. This process may involve complex evaluations or manual steps and is also processed asynchronously after submission.
- Partner API credential provisioning. Our team issues the initial credential after the legal agreement is signed. Subsequent rotations are self-service for the partner.
- Developer documentation and sandbox environment to support the 5-day self-onboarding objective.
- Master data endpoints for accepted vehicle brands and valid enumerations, such as gender, to ensure partners submit consistent values.
- Design Decision: Async Submission Model. All validations, including vehicle identification, are processed after the submission is accepted. The API returns an acknowledgment immediately, and partners receive status updates separately. Rationale: If the existing loan system or credit bureau is unavailable, partners can still submit successfully. Failed validations are retried internally without causing API errors during submission. This design supports the 2-second response objective and improves resilience.

## OUT OF SCOPE
- Consumer-facing UI or partner-facing dashboard.
- Dedicated authentication service, which is unnecessary for a static and pre-vetted partner set.
- Synchronous credit checks or underwriting within the request-response cycle.
- OAuth or session-based authentication flows.
