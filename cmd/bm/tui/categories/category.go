package categorytui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	tuicommands "github.com/evcraddock/bm/cmd/bm/tui/commands"
	"github.com/evcraddock/bm/pkg/categories"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 1)

	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	title = lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		// Bold(true).
		// Background(lipgloss.Color("62")).
		Padding(0, 1).
		Foreground(lipgloss.Color("230")).
		SetString("Categories")
)

type Model struct {
	categories []categories.Category
	category   string
	cursor     int
	selected   map[int]interface{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New(category string) Model {
	manager := categories.NewCategoryManager()
	l, err := manager.GetCategoryList()
	if err != nil {
		panic(err)
	}

	categoryModel := Model{
		categories: l,
		category:   category,
		selected:   make(map[int]interface{}),
	}

	categoryModel.cursor = categoryModel.getSelectedCategoryIndex()
	return categoryModel
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc", "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.categories)-1 {
				m.cursor++
			}

		case "enter", " ", "o", "l":
			cmd = m.markSelected()
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string
	for i, category := range m.categories {
		cursor := "  "
		name := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
			SetString(category.Name).
			String()

		if m.cursor == i {
			cursor = lipgloss.NewStyle().Foreground(highlight).SetString("->").String()
			name = lipgloss.NewStyle().
				// Foreground(lipgloss.Color("#FFF7DB")).
				// Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
				SetString(category.Name).
				String()
		}

		s += lipgloss.NewStyle().Render(fmt.Sprintf("%s %s \n", cursor, name))
	}

	return docStyle.Render(fmt.Sprintf("%s\n\n%s", title, s))
}

func (m Model) getSelectedCategoryIndex() int {
	for i, category := range m.categories {
		if category.Name == m.category {
			return i
		}
	}

	return 0
}

func (m Model) markSelected() tea.Cmd {
	_, ok := m.selected[m.cursor]
	if ok {
		delete(m.selected, m.cursor)
	} else {
		category := m.categories[m.cursor]
		m.selected[m.cursor] = category

		return tuicommands.SelectCategory(category.Name)
	}

	return nil

}
