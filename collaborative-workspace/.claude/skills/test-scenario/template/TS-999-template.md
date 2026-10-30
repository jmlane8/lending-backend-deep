# Test Scenario - [Feature Title]

## Changelog
<!--
Use semantic versioning. Start at 0.0.1 and increment only the patch version. Major/minor versions are user-managed. Sort versions newest first.
Document changes as concise bullet points. Each bullet should ideally be no more than two sentences.
-->

| Version | Date | Changes |
|---------|------|---------|
| `x.y.z` | yyyy-mm-dd | bullet points showing changes |

## Source Requirement
<!--
Mandatory, one or more rows. The use case(s) / user story(ies) this test scenario is derived from.
Use the artifact's official ID as a relative hyperlink where possible.
-->

| ID | Type | Description |
|----|------|-------------|
| [UC-002](../use-cases/UC-002-create-support-ticket.md) | Use Case | Create Customer Support Ticket |
| [US-007](../user-stories/US-007-submit-support-request.md) | User Story | Submit support request |

---

## Test Basis & Assumptions
<!--
Optional, can be zero or more.
Preconditions common to every test case: environment, test accounts, baseline data state.
Do NOT invent requirements the source does not state — flag the gap to the user instead.
-->

- Assumption / shared precondition 1
- Assumption / shared precondition N

## Coverage Matrix
<!--
Map every source element to the test case(s) that cover it, so gaps are visible at a glance.
Source elements come from the source requirement: each Normal Flow / Alternative Flow / Exception
(use case) or each Acceptance Criteria (user story).
Both positive (success) and negative (failure/exception) paths MUST appear.
Every source element must have at least one covering test case — an empty "Covered By" cell is a gap.
-->

| Source Element | Path | Covered By |
|---|---|---|
| Normal Flow / AC-007-001 | Positive | TC-001-001 |
| Alternative Flow AF-001 | Positive | TC-001-002 |
| Exception E-003 / AC-007-005 | Negative | TC-001-005 |

---

## Test Scenarios
<!--
Group related test cases under a named scenario (a high-level behavior to verify).
Each scenario contains one or more test cases.
Across the whole document there MUST be both positive (happy path) and negative (not-happy path) test cases.
-->

### Scenario 1: [Scenario Name]

#### TC-999-001 - [Test Case Title]
<!--
Test case ID format: TC-[TS digits]-[3 digits], e.g. for TS-007 -> TC-007-001.
Each test case has exactly one Expected Result so the outcome is a binary pass/fail.
-->
**Type**: Positive | Negative <br>
**Priority**: High | Medium | Low <br>
**Traces To**: source element (e.g., Normal Flow step 4, AF-001, E-003, AC-007-002) <br>

**Preconditions**
<!-- The state that must be true before this test runs (the "Given"). Set to "None" if not applicable. -->
- Precondition 1

**Test Data**
<!-- Concrete, representative input values. For negative cases use boundary or invalid values. Set to "None" if not applicable. -->
| Field | Value |
|---|---|
| field name | concrete value |

**Steps**
<!-- The actions the tester performs, in order (the "When"). Keep each step atomic. -->
| Step | Action |
|---|---|
| 1 | First action. |
| N | Continue until the behavior is exercised. |

**Expected Result**
<!-- The single observable outcome that determines pass/fail (the "Then"). Specific and binary. -->
- Expected observable outcome.

#### TC-999-00N - [Test Case Title]
**Type**: Positive | Negative <br>
**Priority**: High | Medium | Low <br>
**Traces To**: source element <br>

**Preconditions**
- Precondition 1

**Test Data**
| Field | Value |
|---|---|
| field name | concrete value |

**Steps**
| Step | Action |
|---|---|
| 1 | First action. |
| N | ... |

**Expected Result**
- Expected observable outcome.

### Scenario N: [Scenario Name]
<!-- Repeat the test case structure above for each scenario. -->
