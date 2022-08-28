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
	index      int
	selected   categories.Category
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

	m := Model{
		categories: l,
		windowSize: windowSize,
	}

	m.index = m.Index(category)
	m.selected = l[m.index]
	return m
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
			if m.index > 0 {
				m.index--
			}
			cmd = m.setSelected(false)

		case "down", "j":
			if m.index < len(m.categories)-1 {
				m.index++
			}
			cmd = m.setSelected(false)

		case "g":
			m.index = 0
			cmd = m.setSelected(false)

		case "G":
			if len(m.categories) > 0 {
				m.index = len(m.categories) - 1
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

		if m.index == i {
			cursor = itemSelectedStyle.SetString("->").String()
			name = itemSelectedStyle.Background(lipgloss.Color("62")).SetString(category.Name).String()
		}

		s += lipgloss.NewStyle().Render(fmt.Sprintf("%s %s \n", cursor, name))
	}

	if m.windowSize != nil {
		height := m.windowSize.Height - marginHeight
		docStyle = docStyle.Height(height)
	}

	return docStyle.Render(fmt.Sprintf("%s \n\n%s", title, s))
}

func (m Model) Index(category string) int {
	for i, c := range m.categories {
		if c.Name == category {
			return i
		}
	}

	return 0
}

func (m Model) setSelected(switchView bool) tea.Cmd {
	m.selected = m.categories[m.index]
	return tuicommands.SelectCategory(m.selected.Name, switchView)
}
