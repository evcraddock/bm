package tui

import "github.com/charmbracelet/lipgloss"

var (
	boxStyle         = lipgloss.NewStyle().MarginTop(1)
	selectedBoxStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, true, true, true)
)

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
	leftBoxStyle := lipgloss.NewStyle().MarginTop(1).MarginLeft(1)
	rightBoxStyle := lipgloss.NewStyle().MarginTop(1).MarginLeft(-1)

	switch a.state {
	case categoryView:
		leftBoxStyle = selectedBoxStyle

	case bookmarkView:
		rightBoxStyle = selectedBoxStyle
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftBoxStyle.Render(a.category.View()),
			rightBoxStyle.Render(a.bookmark.View())),
	)
}
