// Package tuikit provides shared Bubble Tea UI components for obsidzen CLI tools.
// It defines theme styles, fixed status glyphs, and the common Screen frame.
package tuikit

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// 상태 글리프 고정셋. 다른 이모지는 쓰지 않는다.
const (
	GlyphStep = "▶"
	GlyphOK   = "✓"
	GlyphWarn = "⚠"
	GlyphErr  = "✗"
)

const (
	KeyQuit  = "q"
	KeyBack  = "esc"
	KeyCtrlC = "ctrl+c"

	HelpSelectRunQuit    = "↑/↓ select · enter run · q quit"
	HelpEnterEscQuit     = "enter select · esc back · q quit"
	HelpEnterOpenEscQuit = "enter open · esc back · q quit"
	HelpAnyHomeQuit      = "any key home · q quit"
)

// ChromeHeight 는 Screen 프레임(테두리+제목+구분선+여백+푸터)이 쓰는 세로 줄 수.
// 본문(리스트 등) 높이 계산에 쓴다.
const ChromeHeight = 7

// 팔레트.
var (
	ColorAccent = lipgloss.Color("12")
	ColorOK     = lipgloss.Color("10")
	ColorErr    = lipgloss.Color("9")
	ColorMuted  = lipgloss.Color("8")
)

// 공통 스타일.
var (
	Frame = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent).
		Padding(0, 1)
	Title  = lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	Rule   = lipgloss.NewStyle().Foreground(ColorMuted)
	Prompt = lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	Cursor = lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	Status = lipgloss.NewStyle().Foreground(ColorMuted)
	Help   = lipgloss.NewStyle().Foreground(ColorMuted)
	Err    = lipgloss.NewStyle().Foreground(ColorErr)
	OK     = lipgloss.NewStyle().Foreground(ColorOK)
)

// Screen 은 모든 화면 공통 레이아웃: 제목 + 구분선 + 본문 + [푸터]를 둥근 테두리로 감싼다.
// 구분선이 가장 넓은 줄이 되어 프레임 폭을 결정하고 나머지가 정렬된다.
func Screen(width int, title, body, footer string) string {
	if width < 24 {
		width = 24
	}
	inner := width - 4 // 테두리(2) + 좌우 패딩(2)
	parts := []string{
		Title.Render(title),
		Rule.Render(strings.Repeat("─", inner)),
		body,
	}
	if footer != "" {
		parts = append(parts, "", footer)
	}
	return Frame.Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

// CursorList renders a dense selectable list with the shared cursor glyph.
func CursorList(items []string, selected int) string {
	if len(items) == 0 {
		return Status.Render("(empty)")
	}
	var b strings.Builder
	for i, item := range items {
		if i == selected {
			b.WriteString(Cursor.Render(GlyphStep+" "+item) + "\n")
		} else {
			b.WriteString("  " + item + "\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// Footer renders optional status plus help text with the shared muted style.
func Footer(status, help string) string {
	if status != "" && help != "" {
		return Status.Render(status) + "\n" + Help.Render(help)
	}
	if status != "" {
		return Status.Render(status)
	}
	return Help.Render(help)
}

// KV is a key/value row for compact status and summary screens.
type KV struct {
	Key   string
	Value string
}

// KeyValues renders aligned key/value rows.
func KeyValues(rows []KV) string {
	if len(rows) == 0 {
		return ""
	}
	width := 0
	for _, row := range rows {
		if len(row.Key) > width {
			width = len(row.Key)
		}
	}
	var b strings.Builder
	for _, row := range rows {
		b.WriteString(fmt.Sprintf("  %-*s  %s\n", width, row.Key+":", row.Value))
	}
	return strings.TrimRight(b.String(), "\n")
}

// Action is a command row for home/action menus.
type Action struct {
	Key         string
	Description string
}

// ActionList renders actions as aligned "key - description" rows.
func ActionList(actions []Action, selected int) string {
	if len(actions) == 0 {
		return Status.Render("(empty)")
	}
	width := 0
	for _, action := range actions {
		if len(action.Key) > width {
			width = len(action.Key)
		}
	}
	items := make([]string, len(actions))
	for i, action := range actions {
		items[i] = fmt.Sprintf("%-*s - %s", width, action.Key, action.Description)
	}
	return CursorList(items, selected)
}

// TrimLines keeps the last max lines of body. max <= 0 keeps body unchanged.
func TrimLines(body string, max int) string {
	if max <= 0 {
		return body
	}
	lines := strings.Split(body, "\n")
	if len(lines) > max {
		lines = lines[len(lines)-max:]
	}
	return strings.Join(lines, "\n")
}

// TextScreen renders a read-only text body with standard screen chrome.
func TextScreen(width int, title, body, status, help string) string {
	return Screen(width, title, body, Footer(status, help))
}

// TailScreen renders streaming tail-style lines, keeping only the last max lines.
func TailScreen(width int, title string, lines []string, empty, status, help string, max int) string {
	body := strings.Join(lines, "\n")
	if body == "" {
		body = Status.Render(empty)
	}
	return TextScreen(width, title, TrimLines(body, max), status, help)
}

func IsQuitKey(key string) bool {
	return key == KeyQuit || key == KeyCtrlC
}

func IsBackKey(key string) bool {
	return key == KeyBack || key == "b"
}

func IsUpKey(key string) bool {
	return key == "up" || key == "k"
}

func IsDownKey(key string) bool {
	return key == "down" || key == "j"
}
