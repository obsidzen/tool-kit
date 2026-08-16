# tui-kit

obsidzen CLI 도구(Go + Bubble Tea)가 공유하는 **TUI 부품** 모듈. Bubble Tea 타입/실행기를 표준 진입점으로 감싸고, 테마(색·스타일), 고정 상태 글리프(`▶ ✓ ⚠ ✗`), 공통 화면 프레임 `Screen()`을 제공해 도구 간 룩과 실행 방식을 통일한다.

버전은 `VERSION` 파일과 module tag로 관리한다.

## API

- `tuikit.Model/Msg/Cmd/KeyMsg/WindowSizeMsg/QuitMsg` — Bubble Tea 타입 alias
- `tuikit.RunAltScreen(model)` — 표준 alt-screen 실행
- `tuikit.ExecProcess(cmd, fn)`, `tuikit.QuitCmd()` — Bubble Tea command helper
- `tuikit.KeyString`, `KeyRunes`, `KeyEnter`, `AssertQuit` — key/test helper
- `tuikit.TextInput`, `NewTextInput`, `TextInputBlink` — Bubbles text input wrapper
- `tuikit.KeyBinding`, `NewKeyBinding`, `MatchesKey` — Bubbles key binding wrapper
- `tuikit.KeyHelpModel`, `NewKeyHelpModel` — Bubbles help wrapper
- `tuikit.FilterListItem`, `FilterListModel`, `NewFilterList` — Bubbles filterable list wrapper
- `tuikit.Language`, `MessageCatalog`, `Translator`, `CommonMessages`, `MergeCatalogs` — TUI 언어 선택/i18n catalog helper
- `tuikit.KeyLanguage`, `tuikit.KeyDetail` — language selector 기본 키(`l`), 상세 설명 기본 키(`?`)
- `tuikit.Screen(width, title, body, footer string) string` — 둥근 테두리 + 제목 + 구분선 + 푸터 프레임
- `tuikit.CursorList(items []string, selected int) string` — 공통 커서 글리프를 쓰는 선택 목록
- `tuikit.Footer(status, help string) string` — 상태/도움말 푸터
- `tuikit.KeyValues(rows []tuikit.KV) string` — 설정/결과 요약용 key-value 표기
- `tuikit.ActionList(actions []tuikit.Action, selected int) string` — `key - description` action 메뉴 정렬 렌더링
- `tuikit.SelectItemsFromStrings`, `SelectFromStrings`, `SelectItemsFromActions` — 문자열/action 목록을 선택 모델로 변환
- `tuikit.TextScreen(...)`, `tuikit.TailScreen(...)`, `tuikit.TrimLines(...)` — 읽기 전용/로그 tail 화면 헬퍼
- `tuikit.RunTaskMenu(title, tasks)`, `RunTaskMenuWithVersion(title, version, tasks)`, `RunTaskMenuWithVersionAndCatalog(title, version, tasks, catalog)`, `RunTaskMenuWithVersionCatalogLanguage(title, version, tasks, catalog, language)`, `NewRunMenuModel(...)` — `run-kit` task 목록을 선택하고 실행 로그를 tail 화면으로 표시하는 공통 Bubble Tea model. version은 home footer에 표시하고, catalog가 있으면 `l` 키로 언어를 선택한다. `?`는 선택한 task의 상세 설명을 연다.
- `tuikit.IsQuitKey/IsBackKey/IsUpKey/IsDownKey` — 공통 키 판정 helper
- `tuikit.HelpSelectRunQuit`, `HelpEnterEscQuit`, `HelpEnterOpenEscQuit`, `HelpAnyHomeQuit` — 공통 help 문구
- reusable models:
  - `SelectModel` — 선택 목록 상태/이동/렌더링
  - `TextModel` — 읽기 전용 텍스트 화면
  - `ResultModel` — key-value 결과 + artifact 목록
  - `TailModel` — 로그/tail 라인 누적 및 렌더링
  - `RunMenuModel` — task 선택, 실행, 취소, 완료 후 home 복귀를 처리하는 action runner
  - `ConfirmModel` — 단축키 기반 확인 화면. 화면별 help 문구는 `Help`로 지정
  - `FormModel` — 가벼운 문자열 필드 폼. `FocusNext/FocusPrev`, `SetValueByKey`, `Value`, `Values`로 Bubbles `textinput` 같은 입력 컴포넌트와 결합. `FormField.Secret`은 화면 표시를 마스킹하고, `FormField.Options` + `CycleValue`는 `true/false` 같은 enum/toggle 필드를 구현할 때 사용한다. `FieldHasOptions`로 현재 필드가 토글 가능한지 판정한다.
- 스타일: `Title Rule Prompt Cursor Status Help Err OK Frame`, 색 `ColorAccent/OK/Err/Muted`
- 글리프: `GlyphStep/OK/Warn/Err`, `ChromeHeight`

## 사용 규칙

- tool 코드는 Bubble Tea/Bubbles를 직접 import하지 않고 `tui-kit`의 공개 API를 통해 사용한다. 새 Bubble Tea/Bubbles 기능이 필요하면 먼저 `tui-kit`에 adapter나 reusable model로 올린 뒤 tool에서 소비한다.
- 선택/읽기전용/결과/tail/확인/form 상태는 먼저 reusable model로 표현한다. 도메인별 label 변환만 tool 내부에 둔다.
- `internal/tui/theme.go`는 `tui-kit` symbol을 로컬 별칭으로 재노출하는 shim만 둔다.
- 높은 수준의 입력/목록 컴포넌트가 필요한 경우에도 navigation/value 상태는 가능한 reusable model에 맡기고, 컴포넌트는 렌더링·입력 처리 adapter로 감싼다.

## 사용

```go
import tuikit "github.com/obsidzen/tool-kit/tui-kit"
// go.mod: require github.com/obsidzen/tool-kit/tui-kit v0.1.0
view := tuikit.Screen(width, "Title", body, footer)
```

Form toggle 예:

```go
form := tuikit.FormModel{Fields: []tuikit.FormField{
    {Key: "enabled", Value: "false", Options: []string{"false", "true"}},
}}
if form.FieldHasOptions() {
    form = form.CycleValue()
}
```

Language catalog 예:

```go
catalog := tuikit.MessageCatalog{
    tuikit.LangKorean: {
        "actions.verify.description": "번들된 assets 검증",
        "actions.verify.detail": "manifest, PMTiles, glyphs, sprites를 검사합니다.",
    },
}
return tuikit.RunTaskMenuWithVersionCatalogLanguage("tool", version, tasks, catalog, tuikit.LangKorean)
```

소비자: `termux-ssh-launcher`, `aictl`, `adb-manager`, `android-app-dev`. (각 도구 `internal/tui/theme.go` 는 이 모듈을 로컬 별칭으로 재노출.)
