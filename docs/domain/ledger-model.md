# Ledger Model

## Entities
- **Account**: `id`, `user_id`, `name`, `type {asset|liability}`, `institution`, `last_sync_at`.
- **JournalEntry**: `id`, `user_id`, `occurred_at`, `source {csv|plaid|manual}`, `external_id`.
- **Posting**: `id`, `entry_id`, `account_id`, `amount_cents` (signed), `currency`, `memo`.
- **Transaction (view)**: joins entry with normalized merchant/category and plaid ids.
- **Schedule**: recurring or one-time expense/income definitions.
- **ScheduleOccurrence**: specific due instances with `status {pending|partial|satisfied}` and optional link to ledger entry.
- **Signature**: recognition hints for matching transactions to schedules.
- **Match**: links occurrences to entries and tracks partial/full amounts.

## Invariants
- Sum of `amount_cents` in postings for a journal entry **MUST** equal zero.
- Money is stored as integer cents in USD for v1.
- Deletions are soft or handled via reversing entries.

## Posting Examples
```
Entry: Paycheck
  Debit  Asset:Checking  +100000
  Credit Income:Wages   -100000
```
```
Entry: Rent
  Debit  Expense:Rent    +120000
  Credit Asset:Checking -120000
```

## Reconciliation Rules
- Schedule occurrences **MUST** be marked satisfied when linked to postings totaling the expected amount.
- Partial matches **MAY** leave the occurrence in `partial` state until fully satisfied.

## Reversing Entries
- To correct mistakes, create a new journal entry that negates prior postings.
- Audit logs **MUST** record original and reversing entry references.

## Signatures
- Stored per schedule to aid auto-detection: merchant, ACH descriptor, MCC, typical amount band, cadence window, account filter.
- False positives **MUST** be down-weighted when dismissed by users.

