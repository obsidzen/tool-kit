# tool-kit

[한국어](README.md) · [English](README.en.md)

> Source: README.md (ko).

Shared Go modules for building small operational CLI/TUI tools.

## Modules

- `cli-kit`: Cobra-based CLI helpers for root/subcommand construction, `menu` entrypoints, shell completion, flag completion, and args/flag validation.
- `run-kit`: process execution helpers, including `Runner`, `CommandSpec`, task execution, and streaming boundaries.
- `tui-kit`: Bubble Tea/Bubbles-based TUI helpers, including type aliases, runners, theme, screen frame, component adapters, key helpers, and reusable models.

## Import Paths

```go
import (
 clikit "github.com/obsidzen/tool-kit/cli-kit"
 runkit "github.com/obsidzen/tool-kit/run-kit"
 tuikit "github.com/obsidzen/tool-kit/tui-kit"
)
```

Depend only on the modules your tool needs.

```sh
go get github.com/obsidzen/tool-kit/cli-kit@v0.1.0
go get github.com/obsidzen/tool-kit/run-kit@v0.1.0
go get github.com/obsidzen/tool-kit/tui-kit@v0.1.0
```

## Development

```sh
./scripts/ci.sh
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution rules and [SECURITY.md](SECURITY.md) for vulnerability reporting.

## Structure

```text
cli-kit/      # Cobra-based CLI helpers
run-kit/      # command runner, task, and streaming helpers
tui-kit/      # Bubble Tea/Bubbles adapters and reusable TUI models
scripts/      # repo-local verification scripts
go.work       # local multi-module development workspace
.gitattributes # line-ending and binary normalization
CONTRIBUTING.md / CONTRIBUTING.en.md # contribution guide
SECURITY.md / SECURITY.en.md         # security reporting policy
```

## Release

Each child module is versioned independently. Release tags use the Go multi-module format `<module>/vX.Y.Z`.

```text
cli-kit/v0.1.0
run-kit/v0.1.0
tui-kit/v0.1.0
```

Use rc tags for staging candidates, such as `tui-kit/v0.2.0-rc.1`, then tag the same commit as `tui-kit/v0.2.0` for the final release.

## Requirements

- Build: Go 1.26, mise
- Runtime: none. Modules are statically linked into consuming tools.

## License

MIT. See [LICENSE](LICENSE).
