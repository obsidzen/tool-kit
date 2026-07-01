# cli-kit

obsidzen CLI 도구(Go + Cobra)가 공유하는 **CLI 부품** 모듈. 공통 root command 설정, action command 등록, shell completion command, enum completion helper를 제공한다.

역할 경계: [tool-kit README](../README.md).

## API
- `clikit.NewRoot(opts clikit.RootOptions) *clikit.Command` — 공통 root command 생성
- `clikit.NewCommand(opts clikit.CommandOptions) *clikit.Command` — 공통 subcommand 생성
- `clikit.Action` — `key`, `description`, `detail`, `aliases`, `RunE`/`RunArgsE` 기반 action command 정의
- `clikit.VersionInfo` — `--version` flag와 `version` command에 쓰는 tool/version/commit/build date 정보
- `clikit.CompletionCommand()` — `completion [bash|zsh|fish|powershell]` command
- `clikit.EnumCompletion(...)` — enum flag completion helper
- `clikit.NoFileCompletion(fn)` — 파일 완성 없이 동적 후보만 제공
- `clikit.MustRegisterFlagCompletion(cmd, name, fn)` — completion 등록 실패를 즉시 드러냄
- `clikit.RequireExactlyOne(map[string]string)` — mutually exclusive flag 값 검증
- `clikit.NoArgs`, `ExactArgs`, `RangeArgs`, `ArbitraryArgs` — Cobra args helper 재노출
- `Action.Args`, `Action.DisableFlagParsing` — action command가 하위 runner나 domain command로 flag/args를 그대로 넘겨야 할 때 사용
- `clikit.ChangedFlags(cmd)` — config 병합 등에 쓰는 changed flag set
- `clikit.Execute(root)` — 공통 에러 출력 후 종료

## 사용
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

root 기본 실행을 직접 제어해야 하는 도구는 `RunE`를 쓰고, `Menu`를 함께 넘기면 `menu` subcommand도 자동으로 생긴다.
`Version`을 넘기면 `--version` flag와 `version` subcommand가 함께 생긴다.
