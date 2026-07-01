package tuikit

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type Model = tea.Model
type Msg = tea.Msg
type Cmd = tea.Cmd
type KeyMsg = tea.KeyMsg
type WindowSizeMsg = tea.WindowSizeMsg
type QuitMsg = tea.QuitMsg

func Quit() Msg {
	return tea.Quit()
}

func QuitCmd() Cmd {
	return tea.Quit
}

func RunAltScreen(model Model) error {
	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	return err
}

func ExecProcess(cmd *exec.Cmd, fn func(error) Msg) Cmd {
	return tea.ExecProcess(cmd, fn)
}

func KeyString(msg Msg) (string, bool) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return "", false
	}
	return key.String(), true
}

func KeyRunes(value string) KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}

func KeyEnter() KeyMsg {
	return tea.KeyMsg{Type: tea.KeyEnter}
}

type TestReporter interface {
	Helper()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

func AssertQuit(t TestReporter, cmd Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("expected quit command")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}
