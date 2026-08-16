package tuikit

import (
	"context"

	runkit "github.com/obsidzen/tool-kit/run-kit"
)

type RunMenuModel struct {
	Width          int
	Title          string
	Version        string
	Tasks          []runkit.Task
	Select         SelectModel
	Languages      []Language
	LanguageSelect SelectModel
	Translator     Translator
	Tail           TailModel
	Height         int
	Status         string
	mode           runMenuMode
	prev           runMenuMode
	cancel         context.CancelFunc
	running        bool
	Runner         runkit.Runner
}

type runMenuMode int

const (
	runMenuModeMenu runMenuMode = iota
	runMenuModeTail
	runMenuModeLanguage
	runMenuModeDetail
)

func NewRunMenuModel(title string, tasks []runkit.Task, runner runkit.Runner) RunMenuModel {
	return NewRunMenuModelWithOptions(RunMenuOptions{Title: title, Tasks: tasks, Runner: runner})
}

type RunMenuOptions struct {
	Title           string
	Version         string
	Tasks           []runkit.Task
	Runner          runkit.Runner
	Languages       []Language
	Catalog         MessageCatalog
	InitialLanguage string
}

func NewRunMenuModelWithOptions(opts RunMenuOptions) RunMenuModel {
	if opts.Runner == nil {
		opts.Runner = runkit.ExecRunner{}
	}
	languages := opts.Languages
	if len(languages) == 0 {
		languages = DefaultLanguages()
	}
	translator := NewTranslator(opts.InitialLanguage, MergeCatalogs(CommonMessages(), opts.Catalog))
	actions := make([]Action, len(opts.Tasks))
	for i, task := range opts.Tasks {
		actions[i] = Action{Key: task.Key, Description: translator.T("actions."+task.Key+".description", task.Description)}
	}
	languageSelect := SelectLanguageByTag(SelectModel{Items: LanguageSelectItems(languages)}, translator.Language)
	return RunMenuModel{
		Title:          opts.Title,
		Version:        opts.Version,
		Tasks:          opts.Tasks,
		Select:         SelectModel{Items: SelectItemsFromActions(actions)},
		Languages:      languages,
		LanguageSelect: languageSelect,
		Translator:     translator,
		Runner:         opts.Runner,
	}
}

func RunTaskMenu(title string, tasks []runkit.Task) error {
	return RunAltScreen(NewRunMenuModel(title, tasks, runkit.ExecRunner{}))
}

func RunTaskMenuWithVersion(title, version string, tasks []runkit.Task) error {
	return RunAltScreen(NewRunMenuModelWithOptions(RunMenuOptions{Title: title, Version: version, Tasks: tasks, Runner: runkit.ExecRunner{}}))
}

func RunTaskMenuWithVersionAndCatalog(title, version string, tasks []runkit.Task, catalog MessageCatalog) error {
	return RunAltScreen(NewRunMenuModelWithOptions(RunMenuOptions{
		Title:   title,
		Version: version,
		Tasks:   tasks,
		Runner:  runkit.ExecRunner{},
		Catalog: catalog,
	}))
}

func RunTaskMenuWithVersionCatalogLanguage(title, version string, tasks []runkit.Task, catalog MessageCatalog, language string) error {
	return RunAltScreen(NewRunMenuModelWithOptions(RunMenuOptions{
		Title:           title,
		Version:         version,
		Tasks:           tasks,
		Runner:          runkit.ExecRunner{},
		Catalog:         catalog,
		InitialLanguage: language,
	}))
}

func (m RunMenuModel) Init() Cmd { return nil }

func (m RunMenuModel) Update(msg Msg) (Model, Cmd) {
	switch msg := msg.(type) {
	case MouseMsg:
		if m.mode == runMenuModeTail {
			switch msg.Button {
			case MouseButtonWheelUp:
				m.Tail = m.Tail.ScrollUp(3)
			case MouseButtonWheelDown:
				m.Tail = m.Tail.ScrollDown(3)
			}
		}
		return m, nil
	case WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// 로그 본문에 실제로 쓸 수 있는 줄 수. 프레임이 먹는 만큼 빼고,
		// 스크롤 위치 안내 한 줄을 위해 하나 더 남긴다.
		if body := msg.Height - ChromeHeight - 1; body > 0 {
			m.Tail.Height = body
		}
	case KeyMsg:
		key := msg.String()
		if m.mode == runMenuModeDetail {
			switch {
			case IsQuitKey(key):
				return m, QuitCmd()
			case IsBackKey(key) || key == KeyDetail:
				m.mode = runMenuModeMenu
				return m, nil
			}
			return m, nil
		}
		if m.mode == runMenuModeLanguage {
			switch {
			case IsQuitKey(key):
				return m, QuitCmd()
			case IsBackKey(key):
				m.mode = m.prev
				return m, nil
			case IsUpKey(key) || IsDownKey(key):
				m.LanguageSelect = m.LanguageSelect.UpdateKey(key)
				return m, nil
			case key == "enter":
				if item, ok := m.LanguageSelect.Selected(); ok {
					m.Translator.Language = item.Key
					m.Select = m.localizedTaskSelect()
					m.Status = m.Translator.T("common.status.language", "language") + ": " + m.Translator.LanguageLabel(m.Languages)
				}
				m.mode = m.prev
				return m, nil
			}
			return m, nil
		}
		if m.mode == runMenuModeTail {
			if m.running {
				if IsQuitKey(key) {
					if m.cancel != nil {
						m.cancel()
					}
					return m, QuitCmd()
				}
				return m, nil
			}
			if key == KeyLanguage {
				m.prev = runMenuModeMenu
				m.mode = runMenuModeLanguage
				return m, nil
			}
			if IsQuitKey(key) {
				return m, QuitCmd()
			}
			// 스크롤 키는 화면을 닫지 않고 소비한다. 되짚어 보려고 위를 눌렀는데
			// 로그가 닫히면 볼 방법 자체가 없어진다.
			if next, ok := scrollTail(m.Tail, key); ok {
				m.Tail = next
				return m, nil
			}
			m.mode = runMenuModeMenu
			m.cancel = nil
			return m, nil
		}
		switch {
		case IsQuitKey(key):
			return m, QuitCmd()
		case key == KeyLanguage:
			m.prev = m.mode
			m.mode = runMenuModeLanguage
			return m, nil
		case key == KeyDetail:
			if _, ok := m.selectedTask(); ok {
				m.mode = runMenuModeDetail
			}
			return m, nil
		case IsUpKey(key) || IsDownKey(key):
			m.Select = m.Select.UpdateKey(key)
		case key == "enter":
			task, ok := m.selectedTask()
			if !ok {
				return m, nil
			}
			ctx, cancel := context.WithCancel(context.Background())
			lines, err := runkit.StreamTaskLines(ctx, m.Runner, task)
			m.mode = runMenuModeTail
			m.cancel = cancel
			m.running = true
			m.Status = GlyphStep + " running " + task.Key
			m.Tail = TailModel{
				Title: m.Title + " · " + task.Key,
				Max:   400,
				Empty: m.Translator.T("common.empty.running", "(running)"),
				Help:  m.Translator.T("common.help.running_quit", "running · q quit"),
			}
			if err != nil {
				m.running = false
				m.Status = GlyphErr + " " + err.Error()
				m.Tail = m.Tail.Append(m.Status)
				return m, nil
			}
			return m, waitRunMenuLine(lines)
		}
	case runMenuLineMsg:
		if msg.err != nil {
			m.Status = GlyphErr + " " + msg.err.Error()
			m.Tail = m.Tail.Append(m.Status)
			m.running = false
			m.cancel = nil
			m.Tail.Help = m.Translator.T("common.help.tail_scroll_quit", HelpTailScrollQuit)
			return m, nil
		}
		if msg.done {
			m.Status = GlyphOK + " " + m.Translator.T("common.status.completed", "completed")
			m.running = false
			m.cancel = nil
			m.Tail.Help = m.Translator.T("common.help.tail_scroll_quit", HelpTailScrollQuit)
			return m, nil
		}
		m.Tail = m.Tail.Append(msg.text)
		return m, waitRunMenuLine(msg.lines)
	}
	return m, nil
}

func (m RunMenuModel) View() string {
	if m.mode == runMenuModeDetail {
		task, _ := m.selectedTask()
		title := m.Title + " · " + task.Key + " · " + m.Translator.T("common.detail.title", "Detail")
		body := m.Translator.T("actions."+task.Key+".detail", task.Detail)
		if body == "" {
			body = m.Translator.T("actions."+task.Key+".description", task.Description)
		}
		return TextScreen(m.Width, title, body, "", m.Translator.T("common.help.detail", "esc back · q quit"))
	}
	if m.mode == runMenuModeLanguage {
		return Screen(
			m.Width,
			m.Translator.T("common.language.title", "Language"),
			m.LanguageSelect.View(),
			Footer("", m.Translator.T("common.help.language", "↑/↓ language · enter apply · esc back · q quit")),
		)
	}
	if m.mode == runMenuModeTail {
		visible, below := m.Tail.Window()
		return TailWindowScreen(m.Width, m.Tail.Title, visible, below, m.Tail.Empty, m.Status, m.Tail.Help)
	}
	status := m.Status
	if status == "" {
		status = m.Version
	}
	return Screen(m.Width, m.Title, m.Select.View(), Footer(status, m.Translator.T("common.help.select_run_detail_language_quit", "↑/↓ select · enter run · ? detail · l language · q quit")))
}

func (m RunMenuModel) selectedTask() (runkit.Task, bool) {
	item, ok := m.Select.Selected()
	if !ok {
		return runkit.Task{}, false
	}
	for _, task := range m.Tasks {
		if task.Key == item.Key {
			return task, true
		}
	}
	return runkit.Task{}, false
}

func (m RunMenuModel) localizedTaskSelect() SelectModel {
	actions := make([]Action, len(m.Tasks))
	for i, task := range m.Tasks {
		actions[i] = Action{Key: task.Key, Description: m.Translator.T("actions."+task.Key+".description", task.Description)}
	}
	next := SelectModel{Items: SelectItemsFromActions(actions), Index: m.Select.Index}
	if next.Index >= len(next.Items) {
		next.Index = len(next.Items) - 1
	}
	if next.Index < 0 {
		next.Index = 0
	}
	return next
}

type runMenuLineMsg struct {
	text  string
	err   error
	done  bool
	lines <-chan runkit.Line
}

func waitRunMenuLine(lines <-chan runkit.Line) Cmd {
	return func() Msg {
		line, ok := <-lines
		if !ok {
			return runMenuLineMsg{done: true, lines: lines}
		}
		return runMenuLineMsg{text: line.Text, err: line.Err, done: line.Done, lines: lines}
	}
}

// scrollTail 은 로그 스크롤 키를 처리한다. 처리했으면 두 번째 값이 true 이고,
// 호출자는 화면을 닫지 말아야 한다.
func scrollTail(tail TailModel, key string) (TailModel, bool) {
	page := tail.Height
	if page <= 1 {
		page = 10
	}
	switch {
	case IsUpKey(key):
		return tail.ScrollUp(1), true
	case IsDownKey(key):
		return tail.ScrollDown(1), true
	case key == "pgup":
		return tail.ScrollUp(page), true
	case key == "pgdown":
		return tail.ScrollDown(page), true
	case key == "home":
		return tail.ScrollTop(), true
	case key == "end":
		return tail.ScrollBottom(), true
	}
	return tail, false
}
