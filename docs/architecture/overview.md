# System Architecture Overview

## Context
Indibills is a personal finance web app for gig and freelance workers. It exposes a FastAPI backend with a Next.js (React + TypeScript) frontend. The backend integrates Plaid, Postgres, Redis, a job runner, and object storage for CSV imports.

## Components
- **Web Client**: Next.js app communicating via JSON over HTTPS.
- **API**: FastAPI application providing REST endpoints and Webhooks.
- **Worker**: Background processor (Dramatiq or Celery) for sync and reconciliation jobs.
- **Databases**:
  - Postgres for relational data and ledger entries.
  - Redis for caching and job queues.
  - Object storage for uploaded CSV files.
- **Plaid**: Third-party bank aggregator.

## Trust Boundaries
```
[Browser] <--HTTPS--> [API] <--TLS--> [Plaid]
                      |
                      +-- [Postgres]
                      +-- [Redis]
                      +-- [Object Storage]
```
Secrets stay server-side. No PII in logs.

## Data Stores
- Postgres: user accounts, ledger, schedules, audit logs.
- Redis: session cache, job queue.
- Object storage: raw CSV imports.

## Runtime Limits
- 30s maximum request time.
- Job retries with exponential backoff.
- Plaid rate limits respected; API throttles per IP/user.

