package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	bookmarktui "github.com/evcraddock/bm/cmd/bm/tui/bookmarks"
	categorytui "github.com/evcraddock/bm/cmd/bm/tui/categories"
)

type sessionState int

const (
	categoryView sessionState = iota
	bookmarkView
)

type App struct {
	bookmark         bookmarktui.Model
	category         categorytui.Model
	selectedCategory string
	state            sessionState
	windowSize       tea.WindowSizeMsg
}

func New(category string) App {
	bookmarkModel := bookmarktui.New(category, nil)
	categoryModel := categorytui.New(category, nil)

	return App{
		bookmark:         bookmarkModel,
		category:         categoryModel,
		selectedCategory: category,
		state:            bookmarkView,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}
