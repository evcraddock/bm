package bookmarktui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

var (
	// Need global constants file
	marginHeight        = 4
	docStyle            = lipgloss.NewStyle().Margin(1, 1)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Model struct {
	bookmark   *bookmarks.Bookmark
	manager    *bookmarks.BookmarkManager
	category   string
	windowSize *tea.WindowSizeMsg
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func New(bookmark *bookmarks.Bookmark, category string, windowSize *tea.WindowSizeMsg) Model {
	m := Model{
		bookmark:   bookmark,
		manager:    bookmarks.NewBookmarkManager(false, category),
		category:   category,
		windowSize: windowSize,
		inputs:     make([]textinput.Model, 5),
	}

	if bookmark == nil {
		m.bookmark = &bookmarks.Bookmark{}
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.SetValue(m.bookmark.Name)
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "URL"
			t.SetValue(m.bookmark.URL)
		case 2:
			t.Placeholder = "Author"
			t.SetValue(m.bookmark.Author)
		case 3:
			t.Placeholder = "Tags"
			t.SetValue(strings.Join(m.bookmark.Tags, ","))
		case 4:
			t.Placeholder = "Category"
			t.SetValue(category)
		}

		m.inputs[i] = t
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
