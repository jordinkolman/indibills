# Indibills – One-Page Product Brief

## Problem
Gig and freelance workers struggle to see how long they can survive on current cash given irregular income and upcoming expenses (including surprise spending).

## Audience
US-based gig drivers (Uber, Lyft, DoorDash) and freelancers (contract or project-based) who are comfortable linking bank accounts.

## Core Value
A clear **days-of-runway** number paired with a deterministic cash-flow simulation that updates as money moves.

## Key Constraints
- Security and PII protection come first: encrypted tokens, minimal data retention, strict logging.
- v1 integrates Plaid for banks and offers manual CSV import as a fallback.
- Platform OAuth connectors and mileage tracking are out of scope for v1.

## v1 Out-of-Scope (Future)
- Automated tax estimation.
- Invoice generation.
- ML forecasting.

## Success Metrics (first 60 days)
- Time to first value (TTFV) is **less than 10 minutes** from signup to viewing a runway.
- Daily/weekly active users reach **≥40%** of signups in week 1.
- Reconciliation tests achieve **≥90%** correctness with zero PII in logs.

## Open Questions
- Which additional gig platforms should be prioritized for future connectors?
- Are users satisfied with the conservative/neutral/optimistic runway modes?
- What manual-import enhancements are needed for non-Plaid banks?

