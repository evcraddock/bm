package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// TODO: get left pane width from first load
	leftWidth        = 40
	rightWidth       = getTerminalSize().Width - leftWidth - 3
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

	contentPane := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftBoxStyle.Width(leftWidth).Render(a.category.View()),
		rightBoxStyle.Width(rightWidth).Render(selectView),
	)

	mainPane := lipgloss.JoinVertical(
		lipgloss.Top,
		contentPane,
	)

	return mainPane
}
