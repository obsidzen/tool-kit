package tuikit

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
)

type TextInput = textinput.Model

func NewTextInput(placeholder, value string) TextInput {
	input := textinput.New()
	input.Placeholder = placeholder
	input.SetValue(value)
	return input
}

var TextInputBlink Cmd = textinput.Blink

type KeyBinding struct {
	binding key.Binding
}

func NewKeyBinding(keys []string, helpKey, helpText string) KeyBinding {
	return KeyBinding{binding: key.NewBinding(key.WithKeys(keys...), key.WithHelp(helpKey, helpText))}
}

func MatchesKey(msg KeyMsg, binding KeyBinding) bool {
	return key.Matches(msg, binding.binding)
}

type KeyHelpModel struct {
	model   help.Model
	ShowAll bool
}

func NewKeyHelpModel() KeyHelpModel {
	return KeyHelpModel{model: help.New()}
}

func (m *KeyHelpModel) SetWidth(width int) {
	m.model.Width = width
}

func (m *KeyHelpModel) Toggle() {
	m.ShowAll = !m.ShowAll
	m.model.ShowAll = m.ShowAll
}

func (m KeyHelpModel) View(short []KeyBinding, full [][]KeyBinding) string {
	m.model.ShowAll = m.ShowAll
	return m.model.View(keyHelpMap{short: short, full: full})
}

type keyHelpMap struct {
	short []KeyBinding
	full  [][]KeyBinding
}

func (m keyHelpMap) ShortHelp() []key.Binding {
	return unwrapKeyBindings(m.short)
}

func (m keyHelpMap) FullHelp() [][]key.Binding {
	groups := make([][]key.Binding, len(m.full))
	for i, group := range m.full {
		groups[i] = unwrapKeyBindings(group)
	}
	return groups
}

func unwrapKeyBindings(bindings []KeyBinding) []key.Binding {
	unwrapped := make([]key.Binding, len(bindings))
	for i, binding := range bindings {
		unwrapped[i] = binding.binding
	}
	return unwrapped
}
