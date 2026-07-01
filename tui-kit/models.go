package tuikit

import (
	"fmt"
	"strings"
)

type SelectItem struct {
	Key         string
	Label       string
	Description string
}

type SelectModel struct {
	Items []SelectItem
	Index int
}

func SelectItemsFromStrings(values []string) []SelectItem {
	items := make([]SelectItem, len(values))
	for i, value := range values {
		items[i] = SelectItem{Key: value, Label: value}
	}
	return items
}

func SelectFromStrings(values []string) SelectModel {
	return SelectModel{Items: SelectItemsFromStrings(values)}
}

func SelectItemsFromActions(actions []Action) []SelectItem {
	items := make([]SelectItem, len(actions))
	for i, action := range actions {
		items[i] = SelectItem{Key: action.Key, Label: action.Key, Description: action.Description}
	}
	return items
}

func (m SelectModel) MoveUp() SelectModel {
	if m.Index > 0 {
		m.Index--
	}
	return m
}

func (m SelectModel) MoveDown() SelectModel {
	if m.Index < len(m.Items)-1 {
		m.Index++
	}
	return m
}

func (m SelectModel) UpdateKey(key string) SelectModel {
	if IsUpKey(key) {
		return m.MoveUp()
	}
	if IsDownKey(key) {
		return m.MoveDown()
	}
	return m
}

func (m SelectModel) Selected() (SelectItem, bool) {
	if m.Index < 0 || m.Index >= len(m.Items) {
		return SelectItem{}, false
	}
	return m.Items[m.Index], true
}

func (m SelectModel) Labels() []string {
	labels := make([]string, len(m.Items))
	width := 0
	for _, item := range m.Items {
		label := item.Label
		if label == "" {
			label = item.Key
		}
		if item.Description != "" && len(label) > width {
			width = len(label)
		}
	}
	for i, item := range m.Items {
		label := item.Label
		if label == "" {
			label = item.Key
		}
		if item.Description != "" {
			label = fmt.Sprintf("%-*s - %s", width, label, item.Description)
		}
		labels[i] = label
	}
	return labels
}

func (m SelectModel) View() string {
	return CursorList(m.Labels(), m.Index)
}

type TextModel struct {
	Title string
	Body  string
	Help  string
	Max   int
}

func (m TextModel) View(width int) string {
	help := m.Help
	if help == "" {
		help = HelpAnyHomeQuit
	}
	return TextScreen(width, m.Title, TrimLines(m.Body, m.Max), "", help)
}

type ResultModel struct {
	Title     string
	Status    string
	Rows      []KV
	Artifacts []string
	Help      string
}

func (m ResultModel) View(width int) string {
	body := KeyValues(m.Rows)
	if len(m.Artifacts) > 0 {
		if body != "" {
			body += "\n\n"
		}
		body += "artifacts:\n"
		for _, artifact := range m.Artifacts {
			body += "  - " + artifact + "\n"
		}
	}
	help := m.Help
	if help == "" {
		help = HelpAnyHomeQuit
	}
	return Screen(width, m.Title, strings.TrimSpace(body), Footer(m.Status, help))
}

type TailModel struct {
	Title string
	Lines []string
	Max   int
	Empty string
	Help  string
}

func (m TailModel) Append(line string) TailModel {
	if m.Max <= 0 {
		m.Max = 400
	}
	if len(m.Lines) >= m.Max {
		m.Lines = m.Lines[len(m.Lines)-m.Max+1:]
	}
	m.Lines = append(m.Lines, line)
	return m
}

func (m TailModel) View(width int) string {
	empty := m.Empty
	if empty == "" {
		empty = "(waiting)"
	}
	help := m.Help
	if help == "" {
		help = HelpAnyHomeQuit
	}
	return TailScreen(width, m.Title, m.Lines, empty, "", help, m.Max)
}

type ConfirmOption struct {
	Key         string
	Label       string
	Description string
}

type ConfirmModel struct {
	Title   string
	Message string
	Options []ConfirmOption
	Help    string
}

func (m ConfirmModel) Choice(key string) (ConfirmOption, bool) {
	for _, option := range m.Options {
		if key == option.Key {
			return option, true
		}
	}
	return ConfirmOption{}, false
}

func (m ConfirmModel) View(width int) string {
	var parts []string
	if m.Message != "" {
		parts = append(parts, m.Message)
	}
	for _, option := range m.Options {
		label := option.Label
		if label == "" {
			label = option.Key
		}
		if option.Description != "" {
			label += " - " + option.Description
		}
		parts = append(parts, "  "+option.Key+"  "+label)
	}
	help := m.Help
	if help == "" {
		help = "choose option · q quit"
	}
	return Screen(width, m.Title, strings.Join(parts, "\n"), Footer("", help))
}

type FormField struct {
	Key     string
	Label   string
	Value   string
	Secret  bool
	Options []string
}

type FormModel struct {
	Title  string
	Fields []FormField
	Index  int
	Help   string
}

func (m FormModel) UpdateKey(key string) FormModel {
	if IsUpKey(key) {
		return m.FocusPrev()
	}
	if IsDownKey(key) {
		return m.FocusNext()
	}
	return m
}

func (m FormModel) FocusNext() FormModel {
	if m.Index < len(m.Fields)-1 {
		m.Index++
	}
	return m
}

func (m FormModel) FocusPrev() FormModel {
	if m.Index > 0 {
		m.Index--
	}
	return m
}

func (m FormModel) SetValue(value string) FormModel {
	if m.Index >= 0 && m.Index < len(m.Fields) {
		m.Fields[m.Index].Value = value
	}
	return m
}

func (m FormModel) SetValueByKey(key, value string) FormModel {
	for i := range m.Fields {
		if m.Fields[i].Key == key {
			m.Fields[i].Value = value
			return m
		}
	}
	return m
}

func (m FormModel) CycleValue() FormModel {
	if m.Index < 0 || m.Index >= len(m.Fields) {
		return m
	}
	options := m.Fields[m.Index].Options
	if len(options) == 0 {
		return m
	}
	current := m.Fields[m.Index].Value
	next := 0
	for i, option := range options {
		if option == current {
			next = (i + 1) % len(options)
			break
		}
	}
	m.Fields[m.Index].Value = options[next]
	return m
}

func (m FormModel) FieldHasOptions() bool {
	if m.Index < 0 || m.Index >= len(m.Fields) {
		return false
	}
	return len(m.Fields[m.Index].Options) > 0
}

func (m FormModel) Field() (FormField, bool) {
	if m.Index < 0 || m.Index >= len(m.Fields) {
		return FormField{}, false
	}
	return m.Fields[m.Index], true
}

func (m FormModel) Value(key string) string {
	for _, field := range m.Fields {
		if field.Key == key {
			return field.Value
		}
	}
	return ""
}

func (m FormModel) Values() map[string]string {
	values := map[string]string{}
	for _, field := range m.Fields {
		values[field.Key] = field.Value
	}
	return values
}

func (m FormModel) View(width int) string {
	rows := make([]KV, len(m.Fields))
	for i, field := range m.Fields {
		key := field.Label
		if key == "" {
			key = field.Key
		}
		if i == m.Index {
			key = GlyphStep + " " + key
		}
		value := field.Value
		if field.Secret && value != "" {
			value = strings.Repeat("*", len(value))
		}
		rows[i] = KV{Key: key, Value: value}
	}
	help := m.Help
	if help == "" {
		help = "↑/↓ field · enter submit · esc back · q quit"
	}
	return Screen(width, m.Title, KeyValues(rows), Footer("", help))
}
