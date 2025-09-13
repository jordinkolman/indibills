# Security Policy

## Reporting
Report vulnerabilities to security@indibills.dev. We acknowledge reports within 2 business days and aim to resolve critical issues within 7 days.

## Dependency Policy
- Dependencies **MUST** come from the approved allowlist.
- Automated tooling **SHOULD** check for known CVEs weekly.

## Code Scanning
- CI **MUST** run static analysis and dependency audits.
- Secrets scanners **MUST** block commits containing credentials.

## Secrets Handling
- Secrets **MUST** be stored in the secret manager; never committed to the repo.
- Access tokens (e.g., Plaid) **MUST** be encrypted at rest and rotated periodically.
- Logs **MUST NOT** contain PII or secret values.

