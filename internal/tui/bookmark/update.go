package bookmarktui

import (
	tea "github.com/charmbracelet/bubbletea"
	tuicommands "github.com/evcraddock/bm/internal/tui/commands"
	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/utils"
)

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
				m.inputs[i].CharLimit = 128

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
