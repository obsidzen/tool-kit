# tui-kit

> Source: README.md (ko).

Shared **TUI components** for obsidzen Go + Bubble Tea tools. It wraps Bubble Tea types and execution helpers, provides a common theme, fixed status glyphs (`▶ ✓ ⚠ ✗`), and reusable screen/model primitives.

Versioning uses the `VERSION` file and module tags.

## API

- `tuikit.Model/Msg/Cmd/KeyMsg/WindowSizeMsg/QuitMsg` — Bubble Tea type aliases
- `tuikit.RunAltScreen(model)` — standard alt-screen runner
- `tuikit.ExecProcess(cmd, fn)`, `tuikit.QuitCmd()` — Bubble Tea command helpers
- `tuikit.KeyString`, `KeyRunes`, `KeyEnter`, `AssertQuit` — key/test helpers
- `tuikit.TextInput`, `NewTextInput`, `TextInputBlink` — Bubbles text input wrapper
- `tuikit.KeyBinding`, `NewKeyBinding`, `MatchesKey` — Bubbles key binding wrapper
- `tuikit.KeyHelpModel`, `NewKeyHelpModel` — Bubbles help wrapper
- `tuikit.FilterListItem`, `FilterListModel`, `NewFilterList` — Bubbles filterable list wrapper
- `tuikit.Language`, `MessageCatalog`, `Translator`, `CommonMessages`, `MergeCatalogs` — TUI language/i18n catalog helpers
- `tuikit.KeyLanguage`, `tuikit.KeyDetail` — default language selector key (`l`) and detail key (`?`)
- `tuikit.Screen(width, title, body, footer string) string` — rounded frame with title, rule, body, and footer
- `tuikit.CursorList(items []string, selected int) string` — selectable list with the shared cursor glyph
- `tuikit.Footer(status, help string) string` — status/help footer
- `tuikit.KeyValues(rows []tuikit.KV) string` — compact key/value summary
- `tuikit.ActionList(actions []tuikit.Action, selected int) string` — aligned `key - description` action menu
- `tuikit.SelectItemsFromStrings`, `SelectFromStrings`, `SelectItemsFromActions` — convert strings/actions into select models
- `tuikit.TextScreen(...)`, `tuikit.TailScreen(...)`, `tuikit.TrimLines(...)` — read-only and log-tail screen helpers
- `tuikit.RunTaskMenu(...)`, `RunTaskMenuWithVersion(...)`, `RunTaskMenuWithVersionAndCatalog(...)`, `RunTaskMenuWithVersionCatalogLanguage(...)`, `NewRunMenuModel(...)` — shared `run-kit` task menu with streaming output, language selection, and detail view
- `tuikit.IsQuitKey/IsBackKey/IsUpKey/IsDownKey` — shared key predicates
- `tuikit.HelpSelectRunQuit`, `HelpEnterEscQuit`, `HelpEnterOpenEscQuit`, `HelpAnyHomeQuit` — shared help strings

Reusable models:

- `SelectModel` — selectable list state, navigation, and rendering
- `TextModel` — read-only text screen
- `ResultModel` — key/value result plus artifact list
- `TailModel` — log/tail line accumulation and rendering
- `RunMenuModel` — action runner for task selection, streaming, cancel, completion, language, and detail screens
- `ConfirmModel` — key-based confirmation screen
- `FormModel` — lightweight string form. `FormField.Secret` masks display values, and `FormField.Options` + `CycleValue` support toggle/enum fields.

## Rules

- Tool code should use `tui-kit` public APIs instead of importing Bubble Tea/Bubbles directly.
- Add new Bubble Tea/Bubbles adapters or reusable models to `tui-kit` before consuming them from tools.
- Prefer reusable select/text/result/tail/confirm/form models before writing domain-specific state.
- Keep `internal/tui/theme.go` in tools as a local alias shim only when a tool wants shorter local names.

## Usage

```go
import tuikit "github.com/obsidzen/tool-kit/tui-kit"

view := tuikit.Screen(width, "Title", body, footer)
```

Form toggle:

```go
form := tuikit.FormModel{Fields: []tuikit.FormField{
    {Key: "enabled", Value: "false", Options: []string{"false", "true"}},
}}
if form.FieldHasOptions() {
    form = form.CycleValue()
}
```

Language catalog:

```go
catalog := tuikit.MessageCatalog{
    tuikit.LangKorean: {
        "actions.verify.description": "번들된 assets 검증",
        "actions.verify.detail": "manifest, PMTiles, glyphs, sprites를 검사합니다.",
    },
}
return tuikit.RunTaskMenuWithVersionCatalogLanguage("tool", version, tasks, catalog, tuikit.LangKorean)
```
