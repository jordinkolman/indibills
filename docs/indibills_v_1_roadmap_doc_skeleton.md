# Indibills — v1 Roadmap & Documentation Skeleton (Revised)

A concise, agent‑friendly plan and doc set tailored for a gig/freelance personal finance app with “runway‑style” budgeting and Plaid‑based account aggregation.

---

## 0) One‑Page Product Brief
**Problem**: Gig and freelance workers struggle to see how long they can survive on current cash given irregular income and upcoming expenses (including surprise spending).

**Audience**: US‑based gig drivers (Uber, Lyft, DoorDash) and freelancers (contract or project‑based), early users comfortable linking bank accounts.

**Core Value**: A clear **“days of runway”** number + simple, deterministic cash‑flow simulation that updates as money moves.

**Primary Job Stories**
- When I get paid unpredictably, I want to know how many days I can last, so I can decide whether to accept more work or cut spending.
- When I log/identify upcoming or recurring income/expenses, I want the runway to update immediately, so I can see the impact.
- When I have unexpected one‑off spending, I want the app to account for it in my runway forecasts without overreacting.

**Key Constraints**
- Security/PII first: encrypted tokens, minimal data retention, strict logging.
- v1 includes: Plaid for banks; manual CSV import fallback; **no platform OAuth connectors or mileage** yet.

**Out‑of‑Scope for v1 (but potential future features)**
- Automated **tax estimation**, **invoice generation**, **ML forecasting**. (Not implemented in v1; considered in roadmap after beta feedback.)

**Success Metrics (first 60 days)**
- TTFV (time to first value): <10 minutes from signup to seeing a runway.
- Daily/Weekly active of ≥40% of signups in week 1.
- ≥90% correctness on reconciliation tests; zero PII in logs.

---

## 1) Rollout Timeline (America/Chicago)
**Week 0: Sep 12–14, 2025 — Foundations**
- Repo/init, CI, secrets management, environments (dev/stage).
- ADRs: double‑entry ledger, money type (integer cents), date/time conventions, **FastAPI** backend, Plaid as aggregator, storage/hosting, frontend framework.
- SECURITY.md, threat model draft, lint/format/test baselines.

**Week 1: Sep 15–21, 2025 — Domain & Data**
- Implement Money value object; define Account/Entry/Posting invariants.
- Schema migration v1; seed data; property tests for ledger (debits==credits).
- Auth + session; minimal UI shell.

**Week 2: Sep 22–28, 2025 — Ingestion & Manual Input**
- Manual CSV import for transactions (bank exports); categorization rules MVP.
- Expense & income schedules (recurring/one‑time), CRUD screens.

**Week 3: Sep 29–Oct 5, 2025 — Runway Engine**
- Deterministic forward cash‑flow simulation (day‑by‑day) using current balance + schedules + expected income + **unexpected‑spend buffer**.
- RunwaySummary API + UI: number of days, projected zero date, key assumptions.
- Acceptance tests on known scenarios.

**Week 4: Oct 6–12, 2025 — Plaid Integration (MVP)**
- Link flow, token exchange, account selection, initial transaction sync.
- Webhook ingestion; idempotent upserts; reconciliation to ledger.

**Week 5: Oct 13–19, 2025 — Recurrence Recognition & Reconciliation**
- Label‑from‑transaction flow; pattern extraction; auto‑detection on new syncs.
- Manual schedule linking to transactions; rule‑based matching.
- Empty/error states, optimistic updates, audit trail, redaction in logs.

**Week 6: Oct 20–26, 2025 — Beta Hardening**
- Performance budget, loading states, rate limiting, abuse protection.
- Security pass vs checklist; docs freeze; invite‑only beta.

**Post‑v1 Backlog**
- Mileage tracker; platform connectors (Uber/Lyft/DoorDash/Upwork); notifications; tax‑time exports; automatic tax estimation; invoicing; ML‑assisted forecasts.

---

## 2) v1 Functional Spec (agent‑oriented)
Use **RFC‑2119** keywords. Include UX states & edge cases. Each story has acceptance criteria (Gherkin‑style abbreviated).

### 2.1 Authentication & Onboarding
- Users **MUST** sign up with email + password (magic link optional later).
- On first login, users **MUST** see a 3‑step checklist: link accounts (or CSV), add/identify expenses & income, view runway.

**Acceptance**
- Given a new user, when they complete at least one of: link bank OR import CSV AND create/identify ≥1 expense, then they **MUST** see a runway number.

### 2.2 Accounts & Balances
- The system **MUST** represent bank/credit accounts and cash wallets as **Accounts** with type (asset/liability).
- Balance **MUST** derive from immutable ledger entries (no mutable running totals).

### 2.3 Transactions & Ledger
- Ingested transactions (Plaid/CSV/manual) **MUST** map to double‑entry journal entries: ≥2 postings that net to zero.
- Transactions **MUST NOT** be deleted; corrections are reversing entries.

### 2.4 Schedules (Expenses & Expected Income)
- Users **MUST** create recurring (RRULE‑like) and one‑time items with due date, amount, and category.
- System **SHOULD** support “flexing” to nearest business day (toggle).

#### 2.4.1 Recurrence Recognition & Linking (v1)
**Expenses**
1) **Label from Transaction**: A user **MUST** be able to mark an existing synced transaction as a recurring expense. The system **MUST** create a Schedule from it and **MUST** store a **recognition signature** (e.g., merchant/ACH descriptor, MCC/category, normalized name, typical amount range, account, cadence).
2) **Manual First, Link Later**: A user **MUST** be able to manually add recurring or single expenses and later **link** them to a specific transaction once it appears. Linking **MUST** mark the schedule occurrence as satisfied (full or partial) and refine the recognition signature.

**Income**
- The same two flows **MUST** exist for income (e.g., weekly contract payment, DoorDash deposits). Labeling income **MUST** create a Schedule with an income kind and its own recognition signature.

**Auto‑Detection**
- On each sync, the system **SHOULD** attempt to match incoming transactions to known schedule signatures and **SHOULD** propose confirmations. False positives **MUST** be dismissible and used to down‑weight that signature.

**Acceptance examples**
- Given a DoorDash deposit with descriptor "DOORDASH*PAY" occurring weekly ±2 days and amount in a ±10% band, when labeled recurring, future transactions matching that signature **SHOULD** auto‑link.
- Given a manually created “Rent $1200 due on 1st (flex to business day)”, when a Plaid transaction for $1200 ±$5 posts within 3 days of the due date, it **SHOULD** be matched and marked satisfied.

### 2.5 Runway Engine
- The engine **MUST** simulate daily cash flows from today until balance < 0.
- Inputs: current cleared balance(s), scheduled expenses, expected incomes, known transfers, **unexpected‑spend buffer**.
- Outputs: days of runway, projected zero date, top drivers & assumptions.
- A “conservative/neutral/optimistic” toggle **MAY** adjust expected‑income inclusion (0%/50%/100%) and buffer intensity.

**Unexpected‑Spend Buffer (v1)**
- The system **MUST** estimate a budget for non‑scheduled, one‑off spending using a **rolling average** of prior unlinked debit spend (e.g., last 30 days) with user‑tunable settings:
  - `buffer_mode` ∈ {off, avg30, avg90, manual}
  - `buffer_floor_cents` (minimum per‑day buffer)
  - `exclude_categories` (e.g., transfers)
- Users **MUST** be able to mark any transaction as “exclude from buffer calculations”.

**Acceptance examples**
- With $1,000 balance and $50/day scheduled burn plus $15/day buffer, runway **MUST** be 1,000 ÷ 65 ≈ 15 days (floor to 15).
- If a $600 income is expected on day 10 and mode=neutral (50%), runway **SHOULD** extend proportionally.

### 2.6 Reconciliation
- System **SHOULD** auto‑match ingested transactions to schedules by signature/date/amount; user can relink.
- Each schedule occurrence **MUST** be marked satisfied once matched; partial matches allowed.

### 2.7 Audit & Safety
- All mutation endpoints **MUST** accept `Idempotency-Key` and write an audit log event (redacted).

---

## 3) System Overview (diagram notes)
**Stack**: **FastAPI** (Python 3.12) + Postgres + Redis (jobs/cache) + Worker (Dramatiq/Celery) + Web client (see §9 recommendation) + Plaid + Object storage (CSV) + Observability (OTel logs/metrics/trace).

**Trust Boundaries**: Browser ↔ API; API ↔ Plaid; API ↔ DB. Secrets server‑side only. No PII in logs.

**Runtime Constraints**: 30s request timeout; job retries w/ backoff; Plaid rate limits respected.

---

## 4) Domain & Data Model

### 4.1 Entities (v1)
- **User**: id, email, password_hash, created_at.
- **Account**: id, user_id, name, type {asset|liability}, institution, last_sync_at.
- **JournalEntry**: id, user_id, occurred_at, source {csv|plaid|manual}, external_id.
- **Posting**: id, entry_id, account_id, amount_cents (signed), currency, memo.
- **Transaction** (view): entry_id + normalized merchant/category + plaid_ids.
- **Schedule**: id, user_id, kind {expense|income}, name, amount_cents, currency, due_date or rrule, category, flexibility {exact|nearest_business_day}.
- **ScheduleOccurrence**: id, schedule_id, due_on, status {pending|partial|satisfied}, linked_entry_id? (nullable), amount_expected_cents, amount_actual_cents.
- **Signature**: id, schedule_id, fields {merchant, ach_descriptor, mcc, typical_amount_band, cadence_window_days, account_filter}.
- **Match**: id, schedule_occurrence_id, entry_id, amount_cents, status {partial|full}.
- **PlaidItem**: user_id, item_id (enc), access_token (enc), institution_id.
- **AuditEvent**: id, user_id, action, entity_ref, redacted_payload, created_at.
- **BufferSettings**: user_id, mode, floor_cents, exclude_categories.

### 4.2 Invariants
- Sum(postings.amount_cents) per entry == 0.
- Money stored as integer cents; currency per user (USD v1).
- All deletions are soft or reversing entries.

---

## 5) API Contract (OpenAPI 3.1 — updated sketch)
```yaml
openapi: 3.1.0
info:
  title: Indibills API
  version: 0.1.0
servers:
  - url: https://api.indibills.dev
paths:
  /auth/signup:
    post: { summary: Create user }
  /accounts:
    get: { summary: List accounts }
  /transactions:
    get: { summary: List transactions }
    post:
      summary: Create manual transaction
      parameters:
        - name: Idempotency-Key
          in: header
          schema: { type: string }
  /schedules:
    get: { summary: List schedules }
    post: { summary: Create schedule (manual) }
  /schedules/{id}/occurrences:
    get: { summary: List occurrences }
  /schedules/from-transaction:
    post: { summary: Create schedule by labeling an existing transaction }
  /transactions/{id}/link-schedule:
    post: { summary: Link a transaction to a schedule occurrence }
  /runway:
    get:
      summary: Get runway summary
      parameters:
        - in: query
          name: optimism
          schema: { type: string, enum: [conservative, neutral, optimistic] }
  /buffer/settings:
    get: { summary: Get unexpected‑spend buffer settings }
    put: { summary: Update unexpected‑spend buffer settings }
  /plaid/link/token:
    post: { summary: Create link token }
  /plaid/item/exchange:
    post: { summary: Exchange public_token for access_token }
  /webhooks/plaid:
    post: { summary: Plaid webhook receiver }
components:
  schemas:
    Money:
      type: object
      properties:
        amount_cents: { type: integer }
        currency: { type: string, enum: [USD] }
    Account:
      type: object
      properties:
        id: { type: string }
        name: { type: string }
        type: { type: string, enum: [asset, liability] }
        balance: { $ref: '#/components/schemas/Money' }
    RunwaySummary:
      type: object
      properties:
        days: { type: integer }
        zero_date: { type: string, format: date }
        assumptions: { type: array, items: { type: string } }
```

---

## 6) Security & Privacy Pack (v1 scope)
- **PII Minimization**: store only what’s needed; never store bank creds; Plaid access tokens encrypted at rest; rotate keys.
- **Logging**: structured; redact names/emails/tokens; disable request body logging for auth/plaid routes.
- **Threat Model (STRIDE)**
  - *Spoofing*: Strong auth, session hardening, webhook secret validation.
  - *Tampering*: Idempotency keys, immutable ledger, signed migrations.
  - *Repudiation*: Audit events for all mutations.
  - *Information Disclosure*: TLS everywhere, field‑level encryption for sensitive columns.
  - *Denial of Service*: Rate limits per IP/user, backpressure on webhooks.
  - *Elevation of Privilege*: RBAC (future multi‑user); least‑privilege tokens.
- **Compliance Posture**: Not a custodial service; no payment initiation; aim for OWASP ASVS L1/L2 coverage check.

---

## 7) Testing Strategy
- **Property tests** (Hypothesis) for ledger invariants and signature matching.
- **Scenario tests** for runway engine (fixtures for varying income/expense cadence + buffer settings).
- **Contract tests** from OpenAPI (request/response validation; negative cases).
- **E2E happy paths**: signup → link/import → create/label schedules → view runway → reconcile.

---

## 8) Ops & Quality Conventions
- **Config** via env; no secrets in code; secret store for production.
- **Logs** as streams; correlation IDs per request; sampling for verbose routes.
- **Migrations** versioned and reversible; seed scripts for demos; feature flags.
- **Runbooks**: Plaid outages, webhook replay, CSV import errors, data redaction requests.

---

## 9) Stack Recommendation (frontend) & Agent‑Safety Rules
**Backend (chosen)**: **FastAPI** (+ Pydantic v2, SQLAlchemy, Alembic, Uvicorn, Redis, Dramatiq/Celery for jobs).

**Frontend (recommended)**: **Next.js (React + TypeScript)** with TanStack Query, Tailwind, and OpenAPI‑generated client types. Rationale: huge ecosystem, excellent docs, easy auth/session patterns, and strong community for finance‑app UI patterns; works well even if backend is your strength. Alternative if you want simpler mental model: **SvelteKit + TypeScript** (leaner, very ergonomic forms), but fewer ready‑made components.

**Agent‑Safety Rules**
- Allowed to: modify code in `/api`, `/web`, `/jobs`, `/infra`; add migrations under `/db/migrations`.
- Not allowed to: alter `/security/*`, change ADRs without creating a new ADR, write secrets, disable lint/tests.
- All schema or API contract changes **MUST** include: migration, OpenAPI update, and an ADR entry.
- External libs **MUST** go through a dependency allowlist.

---

## 10) ADR Index (start here)
1. **Ledger model**: double‑entry with integer cents; immutable entries; reversing corrections.
2. **Money type**: integer cents; USD‑only v1; scale via currency table later.
3. **Time**: server times UTC; display per user tz.
4. **Aggregator**: Plaid for v1; add others only if needed later.
5. **Forecasting**: deterministic simulation + buffer; no ML in v1.
6. **Frontend framework**: Next.js + TS (alt: SvelteKit).
7. **Job runner**: Dramatiq/Celery (pick one).

---

## 11) Project Layout (suggested)
```
/docs/00-brief.md
/docs/specs/v1-functional-spec.md
/docs/architecture/overview.md
/docs/architecture/adr/0001-ledger-model.md
/docs/domain/ledger-model.md
/docs/domain/money.md
/api/openapi.yaml
/security/SECURITY.md
/security/threat-model.md
/tests/
/ops/runbooks.md
```

---

## 12) Skeleton Prompts for Each Doc (to fill fast)

**/docs/00-brief.md**
- Problem, Audience, Value, Constraints, v1 Out‑of‑Scope (future), Success Metrics, Open Questions.

**/docs/specs/v1-functional-spec.md**
- For each feature: Purpose → MUST/SHOULD/MAY → States → Edge cases → Acceptance (Given/When/Then).
- Include Recurrence Recognition flows (label‑from‑transaction; manual then link later) for both expenses and income.
- Include Unexpected‑Spend Buffer parameters and examples.

**/docs/architecture/overview.md**
- Context diagram; components; trust boundaries; data stores; runtime limits.

**/docs/architecture/adr/0001-ledger-model.md**
- Context → Decision → Consequences → Alternatives → Status.

**/docs/domain/ledger-model.md**
- Entities; invariants; posting examples; reconciliation rules; reversing entries; signatures.

**/security/SECURITY.md**
- Reporting policy; dependency policy; code scanning; secrets handling.

**/security/threat-model.md**
- STRIDE table; mitigations; residual risks; test hooks.

**/api/openapi.yaml**
- Keep in sync with code; generate clients/mocks; validate in CI.

---

This revision bakes in: recurrence labeling/recognition for **both expenses and income**, a tunable **unexpected‑spend buffer** that feeds the runway, and your **FastAPI** backend choice with a pragmatic frontend recommendation for your background.

