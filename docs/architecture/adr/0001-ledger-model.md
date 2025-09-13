# ADR 0001: Ledger Model

## Context
Transactions must be recorded with accuracy and traceability. A traditional double-entry ledger ensures every movement of money balances and provides an audit trail.

## Decision
Adopt a double-entry ledger using integer cents. Each journal entry **MUST** contain at least two postings that sum to zero. Corrections are done with reversing entries; deletions are disallowed.

## Consequences
- Guarantees accounting integrity and supports reconciliation.
- Requires developers to model transactions explicitly and handle reversing entries.
- Enables future reporting and audit features without schema changes.

## Alternatives
- Single-entry ledger with mutable balances (rejected due to reconciliation difficulty).
- Floating-point money representation (rejected for precision issues).

## Status
Accepted – Sept 2025.

