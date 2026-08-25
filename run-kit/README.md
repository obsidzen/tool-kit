# run-kit

obsidzen Go 도구가 외부 프로세스를 실행할 때 쓰는 **주입형 Runner** 모듈. 실제 실행은 `ExecRunner`, 테스트는 fake runner로 바꿔 명령 시퀀스를 검증한다.

역할 경계: [tool-kit README](../README.md).

## API

- `runkit.CommandSpec` — 실행 디렉터리, 명령, 인자, 선택적 환경변수
- `runkit.Task` — 하나의 사용자 action에 대응하는 단일 command, command sequence, 또는 Go stream function. `Description`은 메뉴 한 줄 설명, `Detail`은 TUI/CLI 상세 설명에 쓴다.
- `runkit.Runner` — `Run`, `Stream`
- `runkit.ExecRunner` — `os/exec` 기반 기본 구현
- `runkit.Event`, `EventStatus` — CLI/TUI/JSON이 공유하는 renderer-independent lifecycle event
- `runkit.EventSchemaVersion`, `EventJSONSchema` — 필수 schema version과 내장 JSON Schema
- `runkit.EventLine`, `DiagnosticLine`, `Line.DisplayText` — typed event와 source/phase가 있는 child diagnostic을
  기존 task stream에서 일관되게 표시
- `runkit.WriteEventJSON` — 검증된 event를 newline-delimited JSON으로 출력
- `runkit.MergeEnv(map[string]string)` — 현재 환경에 override를 병합
- `runkit.StreamTo(ctx, runner, spec, writer)` — stream command output to a writer
- `runkit.StreamLines(ctx, runner, spec)` — command output을 line channel로 변환해 TUI tail 화면에 연결
- `runkit.TaskSpecs`, `TaskWithArgs`, `RunTask`, `StreamTaskTo`, `StreamTaskLines` — task sequence 실행과 streaming
- `runkit.StreamTaskLinesWithFormatter(ctx, runner, task, formatter)` — TUI tail에 표시되는 `$ command...` 라인을 tool별 formatter/redactor로 바꿔 streaming
- `runkit.CommandLine(spec)` — 표시용 command line 렌더링
- `runkit.RedactCommandLine(spec, sensitiveFlags...)` — `--password value` 같은 flag 다음 인자를 표시용 command line에서 `********`로 마스킹

## 사용

```go
runner := runkit.ExecRunner{}
out, err := runner.Run(ctx, runkit.CommandSpec{
    Dir:  projectDir,
    Name: "adb",
    Args: []string{"devices"},
    Env:  os.Environ(),
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

Go 함수로 구현된 task도 같은 TUI/CLI stream 경로를 쓸 수 있다.

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

CLI와 TUI에 같은 상태를 표시하고 JSON 증거도 남겨야 하면 typed event를 사용한다.

```go
event := runkit.Event{
    SchemaVersion: runkit.EventSchemaVersion,
    EventID: "database-check-running",
    Tool:    "example-check",
    Command: "database",
    PhaseID: "database",
    Status:  runkit.StatusRunning,
    Message: "Check database schema",
    Attempt: 1,
    Progress: &runkit.Progress{Current: 1, Total: 3, Unit: "phases"},
}
line := runkit.EventLine(event)
fmt.Fprintln(os.Stdout, line.DisplayText())
```

status는 `planned`, `running`, `passed`, `failed`, `skipped`, `not-applicable`로 고정된다. 실행 event는
1부터 시작하는 typed `Attempt`가 필요하고 failed event는 stable `ErrorCode`가 필수다. 진행률은
`Progress`의 current, total, unit을 함께 사용해 분모의 의미를 명시한다. `SchemaVersion`은 현재
`EventSchemaVersion`으로 명시해야 하며 구형 또는 누락된 version은 실패한다. 상태 icon·색·spinner는
renderer가 소유하며 `Message`에 넣지 않는다.

TUI에서 실행 command를 보여줄 때 secret 인자가 있으면 formatter를 넘긴다.

```go
lines, err := runkit.StreamTaskLinesWithFormatter(ctx, runner, task, func(spec runkit.CommandSpec) string {
    return runkit.RedactCommandLine(spec, "--password", "--token", "--db-url")
})
```
