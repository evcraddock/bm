package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	bookmarktui "github.com/evcraddock/bm/cmd/bm/tui/bookmark"
	bookmarkstui "github.com/evcraddock/bm/cmd/bm/tui/bookmarks"
	categorytui "github.com/evcraddock/bm/cmd/bm/tui/categories"
	tuicommands "github.com/evcraddock/bm/cmd/bm/tui/commands"
)

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.windowSize = msg

	case tuicommands.CategoryViewMsg:
		a.category = categorytui.New(a.selectedCategory, &a.windowSize)
		a.state = categoryView

	case tuicommands.SelectBookmarkMsg:
		bookmark := msg.SelectedBookmark
		a.bookmark = bookmarktui.New(bookmark, &a.windowSize)
		a.state = bookmarkView

	case tuicommands.SelectCategoryMsg:
		a.selectedCategory = msg.SelectedCategory
		a.bookmarks = bookmarkstui.New(msg.SelectedCategory, &a.windowSize)
		if msg.SwitchView {
			a.state = bookmarksView
		}

	case tuicommands.BookmarkViewMsg:
		a.state = bookmarksView

	case tuicommands.ReloadBookmarksMsg:
		a.bookmarks = bookmarkstui.New(a.selectedCategory, &a.windowSize)

	}

	switch a.state {
	case categoryView:
		category, ccmd := a.category.Update(msg)
		categoryModel, ok := category.(categorytui.Model)
		if !ok {
			panic("unable to load categories")
		}

		a.category = categoryModel
		cmd = ccmd

	case bookmarkView:
		bookmark, bcmd := a.bookmark.Update(msg)
		bookmarkModel, ok := bookmark.(bookmarktui.Model)
		if !ok {
			panic("unable to load bookmark")
		}

		a.bookmark = bookmarkModel
		cmd = bcmd

	case bookmarksView:
		bookmarks, bcmd := a.bookmarks.Update(msg)
		bookmarkModel, ok := bookmarks.(bookmarkstui.Model)
		if !ok {
			panic("unable to load bookmarks")
		}

		a.bookmarks = bookmarkModel
		cmd = bcmd
	}

	cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}
