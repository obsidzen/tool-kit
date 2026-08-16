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
type MouseMsg = tea.MouseMsg

const (
	MouseButtonWheelUp   = tea.MouseButtonWheelUp
	MouseButtonWheelDown = tea.MouseButtonWheelDown
)

type QuitMsg = tea.QuitMsg

func Quit() Msg {
	return tea.Quit()
}

func QuitCmd() Cmd {
	return tea.Quit
}

func RunAltScreen(model Model) error {
	// 대체 화면은 터미널 자체 스크롤백을 쓸 수 없게 하므로, 되짚어 보는 수단을
	// 앱이 직접 제공해야 한다. 휠 이벤트를 받아 로그 스크롤에 쓴다.
	_, err := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
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
