package categorytui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	tuicommands "github.com/evcraddock/bm/internal/tui/commands"
	"github.com/evcraddock/bm/pkg/categories"
)

var (
	docStyle     = lipgloss.NewStyle().Margin(1, 1)
	leftWidth    = 40
	marginHeight = 4

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
)

type Model struct {
	list       list.Model
	manager    *categories.CategoryManager
	index      int
	selected   categories.Category
	windowSize *tea.WindowSizeMsg
}

type bmcategory struct {
	categories.Category
}

func (i bmcategory) Title() string {
	return i.Name
}

func (i bmcategory) Description() string {
	return ""
}

func (i bmcategory) FilterValue() string {
	return i.Name
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New(category string, windowSize *tea.WindowSizeMsg) Model {
	m := Model{
		manager:    categories.NewCategoryManager(),
		windowSize: windowSize,
	}

	list := m.loadCategoryList(category)
	list.Title = "Bookmark Manager"
	list.SetShowFilter(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)

	m.list = list
	return m
}

func (m Model) loadCategoryList(category string) list.Model {
	l, err := m.manager.GetCategoryList()
	if err != nil {
		panic(err)
	}

	index := 0
	items := []list.Item{}
	for i := 0; i < len(l); i++ {
		cat := l[i]
		items = append(items, bmcategory{cat})
		if cat.Name == category {
			index = i
		}
	}

	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(1)

	ci := list.New(items, d, leftWidth, m.windowSize.Height-marginHeight)
	ci.Select(index)

	return ci
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc", "ctrl+c", "q":
			return m, tea.Quit

		case "up", "down", "k", "j", "g", "G":
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, m.setSelected(false))
			cmds = append(cmds, cmd)

		case "enter", " ", "o", "tab":
			cmd = m.setSelected(true)

		case "ctrl+n":
			cmd = m.createBookmark()

		default:
			m.list, cmd = m.list.Update(msg)

		}

	case tea.WindowSizeMsg:
		m.windowSize = &msg
		m.list.SetSize(m.getWindowSize())
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func (b Model) getWindowSize() (int, int) {
	return b.windowSize.Width, b.windowSize.Height - marginHeight
}

func (m Model) setSelected(switchView bool) tea.Cmd {
	selectedItem := m.list.SelectedItem().(bmcategory)
	return tuicommands.SelectCategory(selectedItem.Name, switchView)
}

func (m Model) createBookmark() tea.Cmd {
	return func() tea.Msg {
		return tuicommands.CreateBookmarkMsg(true)
	}
}
