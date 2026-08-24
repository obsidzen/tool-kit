# run-kit

> Source: README.md (ko).

Injectable process runner module for obsidzen Go tools. Production code uses `ExecRunner`; tests can swap in fake runners to verify command sequences.

## API

- `runkit.CommandSpec` — working directory, command, arguments, and optional environment
- `runkit.Task` — one user action backed by a single command, a command sequence, or a Go stream function. `Description` is the one-line menu text; `Detail` is for TUI/CLI detail help.
- `runkit.Runner` — `Run`, `Stream`
- `runkit.ExecRunner` — default `os/exec` implementation
- `runkit.Event`, `EventStatus` — renderer-independent lifecycle events shared by CLI, TUI, and JSON output
- `runkit.EventLine`, `DiagnosticLine`, `Line.DisplayText` — consistent display for typed events and child
  diagnostics that retain source and phase
- `runkit.WriteEventJSON` — validate and write a newline-delimited JSON event
- `runkit.MergeEnv(map[string]string)` — merge overrides into the current environment
- `runkit.StreamTo(ctx, runner, spec, writer)` — stream command output to a writer
- `runkit.StreamLines(ctx, runner, spec)` — convert command output into a line channel for TUI tail screens
- `runkit.TaskSpecs`, `TaskWithArgs`, `RunTask`, `StreamTaskTo`, `StreamTaskLines` — task sequence execution and streaming
- `runkit.StreamTaskLinesWithFormatter(ctx, runner, task, formatter)` — stream task output while formatting/redacting the `$ command...` line shown in TUI tails
- `runkit.CommandLine(spec)` — display command line rendering
- `runkit.RedactCommandLine(spec, sensitiveFlags...)` — mask the argument after flags such as `--password value` in display command lines

## Usage

```go
runner := runkit.ExecRunner{}
out, err := runner.Run(ctx, runkit.CommandSpec{
    Dir:  projectDir,
    Name: "adb",
    Args: []string{"devices"},
})
```

```go
task := runkit.Task{
    Key: "deploy",
    Description: "deploy workers",
    Detail: "Runs env checks and deploys the Workers runtime.",
    Specs: []runkit.CommandSpec{
        {Dir: projectDir, Name: "npm", Args: []string{"run", "env:schema:check"}},
        {Dir: projectDir, Name: "npm", Args: []string{"run", "workers:deploy"}},
    },
}
err := runkit.StreamTaskTo(ctx, runner, task, os.Stdout)
```

Go-function tasks can use the same TUI/CLI streaming path.

```go
task := runkit.Task{
    Key: "verify",
    Description: "verify generated files",
    Stream: func(ctx context.Context) (<-chan runkit.Line, error) {
        lines := make(chan runkit.Line, 1)
        go func() {
            defer close(lines)
            lines <- runkit.Line{Text: "ok"}
            lines <- runkit.Line{Done: true}
        }()
        return lines, nil
    },
}
```

Use a typed event when CLI and TUI must show the same status and JSON evidence is
also required.

```go
event := runkit.Event{
    EventID: "database-check-running",
    Tool:    "example-check",
    Command: "database",
    PhaseID: "database",
    Status:  runkit.StatusRunning,
    Message: "Check database schema",
}
line := runkit.EventLine(event)
fmt.Fprintln(os.Stdout, line.DisplayText())
```

Statuses are limited to `planned`, `running`, `passed`, `failed`, `skipped`, and
`not-applicable`. Failed events require a stable `ErrorCode`. Renderers own status
icons, colors, and spinners; they do not belong in `Message`.

Use a formatter when the displayed command line may include secret arguments.

```go
lines, err := runkit.StreamTaskLinesWithFormatter(ctx, runner, task, func(spec runkit.CommandSpec) string {
    return runkit.RedactCommandLine(spec, "--password", "--token", "--db-url")
})
```
