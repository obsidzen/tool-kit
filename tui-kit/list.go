package tuikit

import (
	"github.com/charmbracelet/bubbles/list"
)

type FilterListItem interface {
	Title() string
	Description() string
	FilterValue() string
}

type FilterListModel struct {
	model list.Model
}

func NewFilterList(items []FilterListItem) FilterListModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(ColorAccent).BorderForeground(ColorAccent)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(ColorAccent).BorderForeground(ColorAccent)
	delegate.SetSpacing(0)

	model := list.New(filterListItems(items), delegate, 0, 0)
	model.SetShowTitle(false)
	model.SetShowHelp(false)
	model.SetShowStatusBar(false)
	return FilterListModel{model: model}
}

func (m FilterListModel) SetSize(width, height int) FilterListModel {
	m.model.SetSize(width, height)
	return m
}

func (m FilterListModel) SetItems(items []FilterListItem) FilterListModel {
	m.model.SetItems(filterListItems(items))
	return m
}

func (m FilterListModel) Filtering() bool {
	return m.model.FilterState() == list.Filtering
}

func (m FilterListModel) Selected() (FilterListItem, bool) {
	item, ok := m.model.SelectedItem().(FilterListItem)
	return item, ok
}

func (m FilterListModel) Update(msg Msg) (FilterListModel, Cmd) {
	var cmd Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m FilterListModel) View() string {
	return m.model.View()
}

func filterListItems(items []FilterListItem) []list.Item {
	converted := make([]list.Item, len(items))
	for i, item := range items {
		converted[i] = item
	}
	return converted
}
