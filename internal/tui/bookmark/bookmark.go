package bookmarktui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	tuicommands "github.com/evcraddock/bm/internal/tui/commands"
	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/utils"
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return m, func() tea.Msg {
				return tuicommands.BookmarksViewMsg(true)
			}

		case "ctrl+c":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, m.updateBookmark()
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					m.inputs[i].SetCursor(0)
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.windowSize != nil {
		height := m.windowSize.Height - marginHeight
		docStyle = docStyle.Height(height)
	}

	return docStyle.Render(b.String())
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) updateBookmark() tea.Cmd {
	updated := m.bookmark

	if m.bookmark != nil && m.bookmark.Name != "" {
		if err := m.manager.Remove(m.bookmark.Name); err != nil {
			panic(err)
		}
	}

	updated.Name = m.inputs[0].Value()
	updated.URL = m.inputs[1].Value()
	updated.Author = m.inputs[2].Value()
	updated.Tags = utils.ToList(m.inputs[3].Value())
	updatedCategory := m.inputs[4].Value()

	if m.category != updatedCategory {
		m.manager = bookmarks.NewBookmarkManager(false, updatedCategory)
	}

	if ok, err := m.manager.Upsert(updated); err != nil {
		panic(err)
	} else if ok {
		return tuicommands.SaveBookmark(updated)
	}

	return nil
}
