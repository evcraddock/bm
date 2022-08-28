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

	title = lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Foreground(lipgloss.Color("230")).
		SetString("Categories")

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	itemSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230"))

	marginHeight = 4
)

type Model struct {
	categories []categories.Category
	category   string
	cursor     int
	selected   map[int]interface{}
	windowSize *tea.WindowSizeMsg
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New(category string, windowSize *tea.WindowSizeMsg) Model {
	manager := categories.NewCategoryManager()
	l, err := manager.GetCategoryList()
	if err != nil {
		panic(err)
	}

	categoryModel := Model{
		categories: l,
		category:   category,
		selected:   make(map[int]interface{}),
		windowSize: windowSize,
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
			cmd = m.setSelected(false)

		case "down", "j":
			if m.cursor < len(m.categories)-1 {
				m.cursor++
			}
			cmd = m.setSelected(false)

		case "g":
			m.cursor = 0
			cmd = m.setSelected(false)

		case "G":
			if len(m.categories) > 0 {
				m.cursor = len(m.categories) - 1
			}
			cmd = m.setSelected(false)

		case "enter", " ", "o", "l", "tab":
			cmd = m.setSelected(true)
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string
	for i, category := range m.categories {
		cursor := "  "
		name := itemStyle.SetString(category.Name).String()

		if m.cursor == i {
			cursor = itemSelectedStyle.SetString("->").String()
			name = itemSelectedStyle.Background(lipgloss.Color("62")).SetString(category.Name).String()
		}

		s += lipgloss.NewStyle().Render(fmt.Sprintf("%s %s \n", cursor, name))
	}

	if m.windowSize != nil {
		height := m.windowSize.Height - marginHeight
		docStyle = docStyle.Height(height)
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

func (m Model) setSelected(switchView bool) tea.Cmd {
	_, ok := m.selected[m.cursor]
	if ok {
		delete(m.selected, m.cursor)
	} else {
		category := m.categories[m.cursor]
		m.selected[m.cursor] = category

		return tuicommands.SelectCategory(category.Name, switchView)
	}

	return nil
}
