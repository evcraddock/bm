package tuicommands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evcraddock/bm/pkg/bookmarks"
)

type BookmarkViewMsg bool
type CategoryViewMsg bool
type ReloadBookmarksMsg bool

type SelectBookmarkMsg struct {
	SelectedBookmark *bookmarks.Bookmark
}

type SelectCategoryMsg struct {
	SelectedCategory string
	SwitchView       bool
}

func SelectCategory(category string, switchView bool) tea.Cmd {
	return func() tea.Msg {
		return SelectCategoryMsg{SelectedCategory: category, SwitchView: switchView}
	}
}

func SelectBookmark(bookmark *bookmarks.Bookmark) tea.Cmd {
	return func() tea.Msg {
		return SelectBookmarkMsg{SelectedBookmark: bookmark}
	}
}
