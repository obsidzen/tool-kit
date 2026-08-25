package tuikit

import (
	"context"
	"io"
	"strings"
	"testing"

	runkit "github.com/obsidzen/tool-kit/run-kit"
)

func TestSelectModel(t *testing.T) {
	m := SelectFromStrings([]string{"a", "b"})
	m = m.UpdateKey("down")
	item, ok := m.Selected()
	if !ok || item.Key != "b" {
		t.Fatalf("selected = %#v, %v", item, ok)
	}
	m = m.UpdateKey("up")
	item, _ = m.Selected()
	if item.Key != "a" {
		t.Fatalf("selected = %#v", item)
	}
}

func TestSelectItemsFromActions(t *testing.T) {
	items := SelectItemsFromActions([]Action{{Key: "run", Description: "Run command"}})
	if len(items) != 1 || items[0].Key != "run" || items[0].Description != "Run command" {
		t.Fatalf("items = %#v", items)
	}
}

func TestSelectModelLabelsAlignDescriptions(t *testing.T) {
	m := SelectModel{Items: SelectItemsFromActions([]Action{
		{Key: "run", Description: "Run command"},
		{Key: "deploy-long", Description: "Deploy command"},
	})}
	labels := m.Labels()
	if labels[0] != "run         - Run command" {
		t.Fatalf("labels[0] = %q", labels[0])
	}
	if labels[1] != "deploy-long - Deploy command" {
		t.Fatalf("labels[1] = %q", labels[1])
	}
}

func TestTailModelAppendTrims(t *testing.T) {
	m := TailModel{Max: 2}
	m = m.Append("one")
	m = m.Append("two")
	m = m.Append("three")
	if got := len(m.Lines); got != 2 {
		t.Fatalf("len = %d", got)
	}
	if m.Lines[0] != "two" || m.Lines[1] != "three" {
		t.Fatalf("lines = %#v", m.Lines)
	}
}

func TestRunMenuUsesSharedEventRendering(t *testing.T) {
	lines := make(chan runkit.Line, 1)
	lines <- runkit.EventLine(runkit.Event{
		SchemaVersion: runkit.EventSchemaVersion,
		PhaseID:       "database",
		Status:        runkit.StatusRunning,
		Message:       "Check database schema",
	})
	close(lines)
	msg := waitRunMenuLine(lines)().(runMenuLineMsg)
	if got, want := msg.text, "running database — Check database schema"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}
}

func TestConfirmModelChoice(t *testing.T) {
	m := ConfirmModel{Options: []ConfirmOption{{Key: "y"}, {Key: "n"}}, Help: "y yes · n no"}
	if option, ok := m.Choice("y"); !ok || option.Key != "y" {
		t.Fatalf("choice = %#v, %v", option, ok)
	}
	if _, ok := m.Choice("x"); ok {
		t.Fatal("unexpected choice")
	}
}

func TestFormModelValues(t *testing.T) {
	m := FormModel{Fields: []FormField{{Key: "name"}, {Key: "mode"}}}
	m = m.SetValue("app")
	m = m.FocusNext()
	m = m.SetValue("debug")
	m = m.SetValueByKey("mode", "release")
	values := m.Values()
	if values["name"] != "app" || values["mode"] != "release" {
		t.Fatalf("values = %#v", values)
	}
	if got := m.Value("mode"); got != "release" {
		t.Fatalf("mode = %q", got)
	}
	field, ok := m.Field()
	if !ok || field.Key != "mode" {
		t.Fatalf("field = %#v, %v", field, ok)
	}
}

func TestFormModelCycleValue(t *testing.T) {
	model := FormModel{Fields: []FormField{{Key: "flag", Value: "false", Options: []string{"false", "true"}}}}
	model = model.CycleValue()
	if got := model.Value("flag"); got != "true" {
		t.Fatalf("cycled value = %q", got)
	}
	if !model.FieldHasOptions() {
		t.Fatal("field should report options")
	}
	model = model.CycleValue()
	if got := model.Value("flag"); got != "false" {
		t.Fatalf("cycled value = %q", got)
	}
}

func TestFormModelMasksSecretFields(t *testing.T) {
	m := FormModel{Fields: []FormField{{Key: "password", Value: "secret", Secret: true}}}
	view := m.View(80)
	if !strings.Contains(view, "******") {
		t.Fatalf("secret value is not masked: %q", view)
	}
	if strings.Contains(view, "secret") {
		t.Fatalf("secret value leaked: %q", view)
	}
}

type testFilterItem string

func (i testFilterItem) Title() string       { return string(i) }
func (i testFilterItem) Description() string { return "desc " + string(i) }
func (i testFilterItem) FilterValue() string { return string(i) }

func TestFilterListModelSelectedAndSetItems(t *testing.T) {
	m := NewFilterList([]FilterListItem{testFilterItem("one")})
	item, ok := m.Selected()
	if !ok || item.Title() != "one" {
		t.Fatalf("selected = %#v, %v", item, ok)
	}
	m = m.SetItems([]FilterListItem{testFilterItem("two")})
	item, ok = m.Selected()
	if !ok || item.Title() != "two" {
		t.Fatalf("selected after set = %#v, %v", item, ok)
	}
}

func TestRunMenuModelStartsSelectedTask(t *testing.T) {
	m := NewRunMenuModel("tool", []runkit.Task{
		{Key: "one", Description: "first", Spec: runkit.CommandSpec{Name: "echo", Args: []string{"ok"}}},
	}, fakeRunMenuRunner{output: "ok\n"})
	updated, cmd := m.Update(KeyEnter())
	if cmd == nil {
		t.Fatal("expected command")
	}
	next := updated.(RunMenuModel)
	if !next.running || next.Status != GlyphStep+" running one" {
		t.Fatalf("model = %#v", next)
	}
}

func TestRunMenuLanguageSelectionLocalizesActions(t *testing.T) {
	m := NewRunMenuModelWithOptions(RunMenuOptions{
		Title: "tool",
		Tasks: []runkit.Task{{Key: "run", Description: "Run command"}},
		Catalog: MessageCatalog{
			LangKorean: {
				"actions.run.description": "명령 실행",
			},
		},
	})
	updated, _ := m.Update(KeyRunes("l"))
	language := updated.(RunMenuModel)
	if language.mode != runMenuModeLanguage {
		t.Fatalf("expected language mode, got %v", language.mode)
	}
	language.LanguageSelect = language.LanguageSelect.UpdateKey("down")
	updated, _ = language.Update(KeyEnter())
	got := updated.(RunMenuModel)
	if got.Translator.Language != LangKorean {
		t.Fatalf("language = %q", got.Translator.Language)
	}
	if got.Select.Items[0].Description != "명령 실행" {
		t.Fatalf("description = %q", got.Select.Items[0].Description)
	}
}

func TestRunMenuInitialLanguageAndDetail(t *testing.T) {
	m := NewRunMenuModelWithOptions(RunMenuOptions{
		Title:           "tool",
		InitialLanguage: LangKorean,
		Tasks: []runkit.Task{{
			Key:         "verify",
			Description: "Verify files",
			Detail:      "Checks files.",
		}},
		Catalog: MessageCatalog{
			LangKorean: {
				"actions.verify.description": "파일 검증",
				"actions.verify.detail":      "필수 파일과 manifest를 검사합니다.",
			},
		},
	})
	if m.Translator.Language != LangKorean {
		t.Fatalf("language = %q", m.Translator.Language)
	}
	if m.Select.Items[0].Description != "파일 검증" {
		t.Fatalf("description = %q", m.Select.Items[0].Description)
	}
	updated, _ := m.Update(KeyRunes("?"))
	detail := updated.(RunMenuModel)
	if detail.mode != runMenuModeDetail {
		t.Fatalf("expected detail mode, got %v", detail.mode)
	}
	if view := detail.View(); !strings.Contains(view, "필수 파일과 manifest를 검사합니다.") {
		t.Fatalf("detail view missing localized text: %q", view)
	}
}

type fakeRunMenuRunner struct {
	output string
}

func (r fakeRunMenuRunner) Run(ctx context.Context, spec runkit.CommandSpec) ([]byte, error) {
	return []byte(r.output), nil
}

func (r fakeRunMenuRunner) Stream(ctx context.Context, spec runkit.CommandSpec) (io.ReadCloser, func() error, error) {
	return io.NopCloser(strings.NewReader(r.output)), func() error { return nil }, nil
}

func TestFormHandleKeyCoversNavigationAndSubmit(t *testing.T) {
	form := FormModel{Fields: []FormField{
		{Key: "name", Label: "Name"},
		{Key: "mode", Label: "Mode", Options: []string{"a", "b"}, Value: "a"},
	}}

	next, submitted, handled := form.HandleKey("tab")
	if !handled || submitted || next.Index != 1 {
		t.Fatalf("tab should move focus: index=%d submitted=%v handled=%v", next.Index, submitted, handled)
	}

	cycled, _, _ := next.HandleKey("space")
	if cycled.Value("mode") != "b" {
		t.Fatalf("space should cycle an option field, got %q", cycled.Value("mode"))
	}

	_, submitted, _ = form.HandleKey("enter")
	if !submitted {
		t.Fatalf("enter should report submission")
	}
}

func TestFormHandleKeyLeavesUnknownKeysToTheCaller(t *testing.T) {
	form := FormModel{Fields: []FormField{{Key: "name", Label: "Name", Value: "ab"}}}

	// The caller still owns text editing, so a rune key must come back unhandled
	// rather than being silently swallowed or double-applied.
	next, submitted, handled := form.HandleKey("c")

	if handled || submitted {
		t.Fatalf("handled=%v submitted=%v, want both false", handled, submitted)
	}
	if next.Value("name") != "ab" {
		t.Fatalf("value = %q, want it untouched", next.Value("name"))
	}
}
