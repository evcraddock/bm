package tui

import (
	"github.com/charmbracelet/lipgloss"
)

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

	return a.bookmarks.View()
}

func (a App) loadMultiPane() string {
	leftBoxStyle := lipgloss.NewStyle().MarginTop(1).MarginLeft(1)
	rightBoxStyle := lipgloss.NewStyle().MarginTop(1).MarginLeft(-1)
	// selectView := a.bookmarks.View()
	var selectView string

	switch a.state {
	case categoryView:
		leftBoxStyle = selectedBoxStyle
		if selectView == "" {
			selectView = a.bookmarks.View()
		}

	case bookmarkView:
		selectView = a.bookmark.View()
		rightBoxStyle = selectedBoxStyle

	case bookmarksView:
		selectView = a.bookmarks.View()
		rightBoxStyle = selectedBoxStyle
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftBoxStyle.Render(a.category.View()),
			rightBoxStyle.Render(selectView)),
	)
}
