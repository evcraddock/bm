package tui

import "github.com/charmbracelet/lipgloss"

func (a App) View() string {
	return a.loadMultiPane()
}

func (a App) loadSinglePane() string {
	if a.state == categoryView {
		return a.category.View()
	}

	return a.bookmark.View()
}

func (a App) loadMultiPane() string {
	leftBox := a.category.View()
	rightBox := a.bookmark.View()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
	)
}
