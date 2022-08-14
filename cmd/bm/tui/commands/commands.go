package tuicommands

import tea "github.com/charmbracelet/bubbletea"

type CategoryViewMsg bool
type ReloadBookmarksMsg bool
type SelectCategoryMsg struct {
	SelectedCategory string
}

func SelectCategory(category string) tea.Cmd {
	return func() tea.Msg {
		return SelectCategoryMsg{SelectedCategory: category}
	}
}
