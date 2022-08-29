package tuicommands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evcraddock/bm/pkg/bookmarks"
)

type BookmarksViewMsg bool
type CategoryViewMsg bool
type CreateBookmarkMsg bool

type ReloadBookmarksMsg struct {
	SelectedIndex int
}

type SaveBookmarkMsg struct {
	SelectedBookmark *bookmarks.Bookmark
}

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

func SaveBookmark(bookmark *bookmarks.Bookmark) tea.Cmd {
	return func() tea.Msg {
		return SaveBookmarkMsg{SelectedBookmark: bookmark}
	}
}

func ReloadBookmarks(index int) tea.Cmd {
	return func() tea.Msg {
		return ReloadBookmarksMsg{SelectedIndex: index}
	}
}
