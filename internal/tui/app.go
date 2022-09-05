package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	bookmarktui "github.com/evcraddock/bm/internal/tui/bookmark"
	bookmarkstui "github.com/evcraddock/bm/internal/tui/bookmarks"
	categorytui "github.com/evcraddock/bm/internal/tui/categories"
)

type sessionState int

const (
	categoryView sessionState = iota
	bookmarkView
	bookmarksView
)

type App struct {
	bookmark         bookmarktui.Model
	bookmarks        bookmarkstui.Model
	category         categorytui.Model
	selectedCategory string
	state            sessionState
	windowSize       tea.WindowSizeMsg
}

func New(category string) App {
	bookmarkModel := bookmarktui.New(nil, category, nil)
	bookmarksModel := bookmarkstui.New(category, "", 0, nil)
	categoryModel := categorytui.New(category, getTerminalSize())

	return App{
		bookmark:         bookmarkModel,
		bookmarks:        bookmarksModel,
		category:         categoryModel,
		selectedCategory: category,
		state:            bookmarksView,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func getTerminalSize() *tea.WindowSizeMsg {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return &tea.WindowSizeMsg{Width: w, Height: h}
}
