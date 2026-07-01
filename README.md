# tool-kit

[한국어](README.md) · [English](README.en.md)

작은 운영용 CLI/TUI 도구를 만들기 위한 공용 Go module 묶음.

## Modules

- `cli-kit`: Cobra 기반 CLI 공통부. root/subcommand 생성, `menu` 진입점, shell completion, flag completion, args/flag validation을 담당한다.
- `run-kit`: 외부 프로세스 실행 공통부. `Runner`/`CommandSpec` 주입 경계와 `ExecRunner` 기본 구현을 담당한다.
- `tui-kit`: Bubble Tea/Bubbles 기반 TUI 공통부. Bubble Tea 타입 alias, 실행기, 테마, 화면 프레임, component adapter, key helper, reusable model을 담당한다.

## Import Paths

```go
import (
	clikit "github.com/obsidzen/tool-kit/cli-kit"
	runkit "github.com/obsidzen/tool-kit/run-kit"
	tuikit "github.com/obsidzen/tool-kit/tui-kit"
)
```

개별 module을 필요한 만큼만 의존한다.

```sh
go get github.com/obsidzen/tool-kit/cli-kit@v0.1.0
go get github.com/obsidzen/tool-kit/run-kit@v0.1.0
go get github.com/obsidzen/tool-kit/tui-kit@v0.1.0
```

## Development

```sh
./scripts/ci.sh
```

기여 절차는 [CONTRIBUTING.md](CONTRIBUTING.md), 보안 신고는 [SECURITY.md](SECURITY.md)를 따른다.

## Structure

```text
cli-kit/      # Cobra 기반 CLI helper
run-kit/      # command runner, task, streaming helper
tui-kit/      # Bubble Tea/Bubbles adapter와 reusable TUI model
scripts/      # repo-local verification scripts
go.work       # local multi-module development workspace
.gitattributes # line-ending and binary normalization
CONTRIBUTING.md / CONTRIBUTING.en.md # 기여 가이드
SECURITY.md / SECURITY.en.md         # 보안 신고 정책
```

## Release

하위 module은 독립적으로 version pinning한다. release tag는 Go multi-module 규칙에 맞춰 `<module>/vX.Y.Z` 형식을 사용한다.

```text
cli-kit/v0.1.0
run-kit/v0.1.0
tui-kit/v0.1.0
```

staging 후보가 필요하면 `tui-kit/v0.2.0-rc.1`처럼 rc tag를 사용하고, 정식 release는 같은 commit에 `tui-kit/v0.2.0` tag를 단다.

## Requirements

- Build: Go 1.26, mise
- Runtime: 없음. 각 module은 소비 tool에 정적으로 링크된다.

## License

MIT. [LICENSE](LICENSE)를 참고한다.
