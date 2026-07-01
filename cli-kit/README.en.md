# cli-kit

> Source: README.md (ko).

Shared **CLI components** for obsidzen Go + Cobra tools. It provides common root command setup, action command registration, shell completion commands, and enum completion helpers.

## API
- `clikit.NewRoot(opts clikit.RootOptions) *clikit.Command` — creates a common root command
- `clikit.NewCommand(opts clikit.CommandOptions) *clikit.Command` — creates a common subcommand
- `clikit.Action` — action command definition with `key`, `description`, `detail`, `aliases`, and `RunE`/`RunArgsE`
- `clikit.VersionInfo` — tool/version/commit/build date metadata for `--version` and `version`
- `clikit.CompletionCommand()` — `completion [bash|zsh|fish|powershell]` command
- `clikit.EnumCompletion(...)` — enum flag completion helper
- `clikit.NoFileCompletion(fn)` — dynamic candidates with file completion disabled
- `clikit.MustRegisterFlagCompletion(cmd, name, fn)` — fail fast on completion registration errors
- `clikit.RequireExactlyOne(map[string]string)` — mutually exclusive flag value validation
- `clikit.NoArgs`, `ExactArgs`, `RangeArgs`, `ArbitraryArgs` — re-exported Cobra args helpers
- `Action.Args`, `Action.DisableFlagParsing` — pass-through commands for forwarding flags/args to downstream runners or domain commands
- `clikit.ChangedFlags(cmd)` — changed flag set for config merging
- `clikit.Execute(root)` — common error output and exit handling

## Usage
```go
root := clikit.NewRoot(clikit.RootOptions{
    Use:   "my-tool",
    Short: "My tool",
    Version: clikit.VersionInfo{Name: "my-tool", Version: version, Commit: commit, Date: date},
    Menu:  func(ctx context.Context) error { return runTUI() },
    Actions: []clikit.Action{
        {Key: "status", Description: "show status", Detail: "Show current service status and important paths.", RunE: runStatus},
    },
})
clikit.Execute(root)
```

Tools that need custom root behavior can provide `RunE`; passing `Menu` as well adds the standard `menu` subcommand.
Passing `Version` adds both `--version` and the `version` subcommand.
