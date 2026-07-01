# Contributing

[한국어](CONTRIBUTING.md) · [English](CONTRIBUTING.en.md)

> Source: CONTRIBUTING.md (ko).

This document defines the human-facing contribution rules for `tool-kit`. AI agent work follows `AGENTS.md`.

## Development Flow

- Use GitHub-flow. Do not push directly to `main`; create PRs from short-lived work branches.
- PR titles and commit messages follow Conventional Commits: `<type>(<scope>): <subject>`.
- Squash merge is the default merge strategy.

## Getting Started

```bash
./scripts/ci.sh
```

## Definition of Done

Before opening a PR, check:

- [ ] `./scripts/ci.sh` passes.
- [ ] Public API, import path, or CLI/TUI usage changes update the relevant module `README.md` and `README.en.md`.
- [ ] Repo operation or security reporting changes update `CONTRIBUTING.md`/`CONTRIBUTING.en.md` or `SECURITY.md`/`SECURITY.en.md`.
- [ ] New behavior has appropriate tests or verification steps.

## Code Rules

- Keep each kit as a small independent Go module.
- When a consuming tool needs direct Bubble Tea/Bubbles usage, first consider whether it belongs in `tui-kit` as an adapter or reusable model.
- Put external process execution behind the testable `run-kit.Runner` boundary.
- Keep public APIs small and explicit.
- Comments should explain why, invariants, or hazards that the code cannot express by itself.

## Issues And PRs

- Use the forms under `.github/ISSUE_TEMPLATE/` for bugs, feature requests, tasks, and questions.
- Follow `.github/pull_request_template.md` for PR descriptions.
- Do not report security vulnerabilities in public issues or PRs. Follow `SECURITY.md`.

## Releases

- Each module is versioned independently.
- Release tags use the Go multi-module format `<module>/vX.Y.Z`.
- Tags are created only on commits that pass `./scripts/ci.sh`.
