# Ops Runbooks

## Plaid Outage
1. Confirm outage via Plaid status page.
2. Pause sync jobs; queue retries.
3. Communicate status to users via banner.
4. Resume jobs once Plaid reports recovery; monitor for backlog clearing.

## Webhook Replay
1. Verify webhook signature.
2. If replay, ensure `Idempotency-Key` prevents duplicate processing.
3. Log replay event with correlation ID.
4. Reprocess if previous attempt failed.

## CSV Import Errors
1. Validate CSV schema; reject rows with missing required fields.
2. Return aggregated error report to user.
3. Allow user to fix and re-upload; prior valid entries **MUST NOT** duplicate.

## Data Redaction Request
1. Authenticate requester.
2. Queue redaction job referencing user ID.
3. Remove personal data from logs, storage, and database per policy.
4. Confirm completion to requester within SLA.

