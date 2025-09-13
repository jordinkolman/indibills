# Threat Model

| STRIDE Category | Threat Example | Mitigation | Residual Risk | Test Hook |
| --- | --- | --- | --- | --- |
| Spoofing | Impersonated user sessions | Strong auth, secure cookies | Low if passwords reused elsewhere | Session hijack test cases |
| Tampering | Ledger modification | Idempotency keys, immutable postings | None | Property tests on ledger |
| Repudiation | User denies action | Audit events with redaction | Low | Verify audit log entries |
| Information Disclosure | PII leakage in logs | Structured logging with redaction, TLS | Medium via 3rd-party breach | Log scanning in CI |
| Denial of Service | Plaid webhook flood | Rate limiting, backpressure | Medium | Load tests for webhook endpoints |
| Elevation of Privilege | Access to other users' data | RBAC (future), strict token scoping | Medium until multi-user | Integration tests on auth checks |

