# tool-kit

> Zero-context agent entrypoint. Keep this file self-contained and aligned with the repo.
> `CLAUDE.md` and `GEMINI.md` are symlinks to this file; edit `AGENTS.md` only.

## Overview

Go 1.26 multi-module repo for shared CLI/TUI building blocks:

- `cli-kit`: Cobra command helpers.
- `run-kit`: injectable command runner and task execution helpers.
- `tui-kit`: Bubble Tea/Bubbles adapters and reusable TUI models.

## Commands

```bash
./scripts/ci.sh
```

**Definition of Done — before finishing**

- [ ] `./scripts/ci.sh` passes.
- [ ] Public API changes are reflected in the relevant module `README.md` and `README.en.md`.
- [ ] Module paths and tag examples remain aligned with `github.com/obsidzen/tool-kit/<kit>`.
- [ ] `AGENTS.md` still matches the actual structure, commands, and testing approach.

## Project Structure

```text
cli-kit/      # module github.com/obsidzen/tool-kit/cli-kit
run-kit/      # module github.com/obsidzen/tool-kit/run-kit
tui-kit/      # module github.com/obsidzen/tool-kit/tui-kit
scripts/      # repo-local verification scripts
.github/      # PR, issue, and CI configuration
```

Root files:

- `README.md`: Korean source README.
- `README.en.md`: English translation; keep it aligned with `README.md`.
- `CONTRIBUTING.md`: Korean source contribution guide.
- `CONTRIBUTING.en.md`: English contribution guide translation.
- `SECURITY.md`: Korean source security reporting policy.
- `SECURITY.en.md`: English security reporting policy translation.
- `LICENSE`: MIT license.
- `go.work`: local multi-module development workspace.
- `mise.toml`: tool versions.
- `.gitattributes`: line-ending and binary normalization.

## Code Style

Keep APIs small, explicit, and tool-agnostic.

```go
type Runner interface {
	Run(ctx context.Context, spec CommandSpec) ([]byte, error)
	Stream(ctx context.Context, spec CommandSpec) (<-chan Line, error)
}
```

- Add shared Bubble Tea/Bubbles adapters or reusable models to `tui-kit` before consumers import those libraries directly.
- Keep command execution behind `run-kit.Runner` when behavior needs tests.
- Use English for code identifiers, package docs, CLI-facing API names, and commit messages.
- Keep README source in Korean and maintain `README.en.md` for public surface docs.
- Code comments should explain invariants or non-obvious behavior, not restate code.
- Do not reference external workspace paths or issue numbers in code comments.

## Testing

- Test command: `./scripts/ci.sh`.
- Each module keeps unit tests beside its package.
- Use fake runners for process execution behavior.
- Tests must not require device, emulator, network service, or host-specific state.
- The GitHub Actions workflow is a thin trigger for `./scripts/ci.sh`; keep CI logic in the script.

## Git Workflow

Use GitHub-flow, Conventional Commits, squash merge by default, and SemVer module tags.

- Default branch: `main`.
- Work branches: `<type>/<issue#>-<kebab-summary>`; issue number may be omitted.
- Commit format: `<type>(<scope>): <subject>`.
- Types: `feat fix docs refactor test chore build ci perf style revert`.
- Breaking changes: `feat!:` or `BREAKING CHANGE:`.
- PR merge default: squash merge.
- Rebase personal work branches on latest `main`; do not rewrite shared branch history.
- GitHub repo settings: default branch `main`, squash merge enabled by default, delete head branches after merge.
- Branch protection: no direct push or force push to `main`; required check is `./scripts/ci.sh`.
- Module tags: `cli-kit/vX.Y.Z`, `run-kit/vX.Y.Z`, `tui-kit/vX.Y.Z`.
- Release candidates: `<module>/vX.Y.Z-rc.N`; final release tags the same commit as `<module>/vX.Y.Z`.
- Backports, if ever needed, use `cherry-pick -x` onto a maintenance branch.
- PR titles use Conventional Commit format. Public API or usage changes must update `README.md` and `README.en.md`.

## Boundaries

- **Always:** keep public API docs and tests aligned with code changes.
- **Always:** keep root and module README translations aligned when public behavior changes.
- **Always:** keep `scripts/ci.sh` as the single local and CI verification gate.
- **Ask first:** module path changes, dependency upgrades, toolchain pin changes, file or directory moves.
- **Ask first:** repo initialization, commit, push, tag publishing, or release publishing.
- **Never:** commit secrets, `.env` values, generated build caches, or machine-local state.
- **Never:** add tests that require real devices, emulators, external services, or private host configuration.
- **Never:** run public PR jobs on self-hosted runners.
