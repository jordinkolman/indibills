# v1 Functional Spec

Use **RFC-2119** keywords. Each feature includes purpose, states, edge cases, and acceptance criteria.

## 1. Authentication & Onboarding
**Purpose**
Provide secure signup/login and guide new users through initial setup.

**Requirements**
- Users **MUST** sign up with email and password.
- After first login, a 3-step checklist **MUST** prompt users to: link accounts (Plaid or CSV), add/identify expenses & income, and view runway.

**States**
1. `new` – no checklist items complete.
2. `partial` – at least one item complete.
3. `ready` – link/import and ≥1 expense/income identified.

**Edge Cases**
- User closes browser mid-checklist; progress **MUST** persist.
- Duplicate email signup **MUST** be rejected.

**Acceptance**
- **Given** a new user
- **When** they link a bank or import a CSV and create or identify at least one expense
- **Then** they **MUST** see a runway number.

## 2. Accounts & Balances
**Purpose**
Represent financial accounts and derive balances from ledger entries.

**Requirements**
- Accounts **MUST** have type `asset` or `liability`.
- Balances **MUST** derive from immutable ledger postings; no mutable running totals.

**States**
- `active` – synced or manually managed.
- `disabled` – user removed or deauthorized; historical data retained.

**Edge Cases**
- Negative balances on liabilities **MUST** display with sign.
- Unknown account types from Plaid **MUST** map to default `asset` until user edits.

**Acceptance**
- **Given** an account with postings totaling $200 debit and $50 credit
- **When** balance is computed
- **Then** resulting balance **MUST** be $150.

## 3. Transactions & Ledger
**Purpose**
Capture all money movement with double-entry journal entries.

**Requirements**
- Imported or manual transactions **MUST** map to journal entries with ≥2 postings summing to zero.
- Transactions **MUST NOT** be deleted; corrections are reversing entries.

**States**
- `pending` – imported but not reviewed.
- `cleared` – verified and reconciled.

**Edge Cases**
- Duplicate imports **MUST** be idempotent.
- Partial reversals **SHOULD** allow correcting amount while retaining audit history.

**Acceptance**
- **Given** a transaction mistakenly categorized
- **When** a reversing entry is posted
- **Then** the ledger **MUST** net to zero with an audit trail.

## 4. Schedules (Expenses & Expected Income)
**Purpose**
Track recurring and one-time cash flows.

**Requirements**
- Users **MUST** create schedules with due date or RRULE, amount, and category.
- The system **SHOULD** support flexing to nearest business day.

**States**
- `pending` – upcoming occurrence.
- `partial` – linked but not fully satisfied.
- `satisfied` – amount fully matched to transactions.

**Edge Cases**
- Editing amount after occurrences exist **MUST** only affect future ones.
- Skipped occurrences **MAY** be manually marked satisfied.

### 4.1 Recurrence Recognition & Linking
**Expenses**
1. **Label from Transaction** – From a synced transaction, user marks it recurring. The system **MUST** create a schedule and store a recognition signature (merchant, descriptor, amount band, account, cadence).
2. **Manual First, Link Later** – User adds a manual schedule and later links it when the transaction arrives. Linking **MUST** satisfy the occurrence and refine the signature.

**Income**
The same flows **MUST** exist for income items. Labeling income **MUST** create a schedule with income kind and recognition signature.

**Auto-Detection**
On each sync, the system **SHOULD** match incoming transactions to known signatures and **SHOULD** propose confirmations. False positives **MUST** be dismissible and down-weight the signature.

**Acceptance Examples**
- **Given** a DoorDash deposit with descriptor "DOORDASH*PAY" every week ±2 days and amount within ±10%
- **When** labeled recurring
- **Then** future matches **SHOULD** auto-link.
- **Given** a manual "Rent $1200 due on 1st (flex to business day)"
- **When** a $1200 ±$5 transaction posts within 3 days of due date
- **Then** it **SHOULD** be matched and marked satisfied.

## 5. Runway Engine
**Purpose**
Forecast days of cash runway.

**Requirements**
- Engine **MUST** simulate daily cash flows until balance < 0.
- Inputs include current balances, schedules, expected income, known transfers, and unexpected-spend buffer.
- Outputs **MUST** include days of runway, projected zero date, and assumptions.
- A mode toggle **MAY** adjust expected-income inclusion (0%/50%/100%) and buffer intensity.

**Unexpected-Spend Buffer**
- System **MUST** estimate non-scheduled spend using a rolling average of prior unlinked debit spend.
- Supported settings: `buffer_mode` ∈ {off, avg30, avg90, manual}; `buffer_floor_cents`; `exclude_categories`.
- Users **MUST** be able to mark transactions as excluded from buffer calculations.

**Examples**
- $1,000 balance with $50/day scheduled burn and $15/day buffer → runway **MUST** be 15 days.
- $600 income expected day 10 with neutral mode (50%) **SHOULD** extend runway proportionally.

**Edge Cases**
- Negative starting balance **MUST** yield zero days.
- Buffer mode `off` **MUST** omit buffer from simulation.

**Acceptance**
- **Given** a user with configured schedules and buffer
- **When** they request a runway summary
- **Then** API **MUST** return days, zero date, and assumptions.

## 6. Reconciliation
**Purpose**
Match transactions to schedule occurrences.

**Requirements**
- System **SHOULD** auto-match transactions via signature/date/amount.
- Users **MUST** be able to relink or unlink matches.
- Occurrences **MUST** be marked satisfied once matched; partial matches allowed.

**Edge Cases**
- Multiple transactions satisfying one occurrence **MUST** mark partial then satisfied.
- Over-matching **MUST** be reversible.

**Acceptance**
- **Given** a schedule occurrence and matching transaction
- **When** auto-matched
- **Then** status **MUST** change to satisfied.

## 7. Audit & Safety
**Purpose**
Ensure traceability and safe retries.

**Requirements**
- All mutation endpoints **MUST** accept an `Idempotency-Key` header.
- Every mutation **MUST** emit a redacted audit log event.

**Edge Cases**
- Replayed requests with same key **MUST** return prior result.

**Acceptance**
- **Given** a valid mutation request with `Idempotency-Key`
- **When** the request is repeated
- **Then** the system **MUST** return the same response and log once.

