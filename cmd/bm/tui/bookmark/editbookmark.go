package bookmarktui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	tuicommands "github.com/evcraddock/bm/cmd/bm/tui/commands"
	"github.com/evcraddock/bm/pkg/bookmarks"
)

var (
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
	windowSize *tea.WindowSizeMsg
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func New(bookmark *bookmarks.Bookmark, windowSize *tea.WindowSizeMsg) Model {
	m := Model{
		bookmark:   bookmark,
		windowSize: windowSize,
		inputs:     make([]textinput.Model, 3),
	}

	if bookmark != nil {
		var t textinput.Model
		for i := range m.inputs {
			t = textinput.New()
			t.CursorStyle = cursorStyle
			t.CharLimit = 32

			switch i {
			case 0:
				t.Placeholder = "Nickname"
				t.SetValue(bookmark.Title)
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
			case 1:
				t.Placeholder = "URL"
				t.SetValue(bookmark.URL)
			case 2:
				t.Placeholder = "Author"
				t.SetValue(bookmark.Author)
			}

			m.inputs[i] = t
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return m, func() tea.Msg {
				return tuicommands.BookmarkViewMsg(true)
			}

		// case "shift+tab":
		// 	return m, func() tea.Msg {
		// 		return tuicommands.CategoryViewMsg(true)
		// 	}

		case "ctrl+c":
			return m, tea.Quit

		// case "ctrl+r":
		// 	m.cursorMode++
		// 	if m.cursorMode > textinput.CursorHide {
		// 		m.cursorMode = textinput.CursorBlink
		// 	}
		// 	cmds := make([]tea.Cmd, len(m.inputs))
		// 	for i := range m.inputs {
		// 		cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
		// 	}
		// 	return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
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
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		m.windowSize = &msg
		// b.list.SetSize(b.getWindowSize())
	}

	// cmds = append(cmds, cmd)
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

	// b.WriteString(helpStyle.Render("cursor mode is "))
	// b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	// b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
