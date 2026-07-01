# Security Policy

[한국어](SECURITY.md) · [English](SECURITY.en.md)

> Source: SECURITY.md (ko).

## Reporting

Do not report security vulnerabilities in public issues, discussions, or PRs.

Use GitHub private vulnerability reporting when available. If it is not available, contact the maintainer through a private channel.

Include as much of the following as possible:

- Affected module: one of `cli-kit`, `run-kit`, or `tui-kit`
- Affected version or commit
- Reproduction steps
- Expected impact
- Relevant logs, configuration, or inputs
- Minimal proof of concept that is safe to share

## Scope

Supported:

- Current `main`
- Latest stable module release
- Recent release candidates, if any

Not explicitly supported:

- Changes that exist only in forks
- End-of-life release lines
- Consumer tool code or operating environment issues
- Local misconfiguration, exposed personal secrets, or issues in the user's own operating environment

## Response Expectations

- The maintainer will acknowledge reports as soon as practical.
- Reproducibility and impact are checked before choosing a fix path.
- When needed, a security fix release and advisory are prepared together.

Response time is not guaranteed, but avoid public disclosure before coordination.

## Secret Handling

- Do not paste tokens, keys, real `.env` values, or private endpoints directly into the report.
- If a secret has already been exposed, revoke it immediately and replace it with a new value.
- If secret exposure is suspected, check whether rotation is required separately from any code fix.

## Disclosure

Vulnerabilities are disclosed after a fix, release, and user mitigation guidance are ready. High-risk exploit details are limited to the necessary scope.
