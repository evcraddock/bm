package categorytui

import (
	"sort"

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
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)

	m.list = list
	return m
}

func (m Model) loadCategoryList(category string) list.Model {
	l, err := m.manager.GetCategoryList()
	if err != nil {
		panic(err)
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].Name < (l[j].Name)
	})

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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = &msg
		m.list.SetSize(m.getWindowSize())

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {

		case "esc", "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ", "o", "tab":
			return m, m.setSelected(true)

		case "ctrl+n":
			return m, m.createBookmark()

		}

	}

	m.list, cmd = m.list.Update(msg)
	return m, tea.Batch(cmd, m.setSelected(false))
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func (b Model) getWindowSize() (int, int) {
	return b.windowSize.Width, b.windowSize.Height - marginHeight
}

func (m Model) GetSelected() *bmcategory {
	item := m.list.SelectedItem()
	if item == nil {
		return nil
	}

	value := item.(bmcategory)
	return &value
}

func (m Model) setSelected(switchView bool) tea.Cmd {
	s := m.GetSelected()
	if s == nil {
		return nil
	}

	selectedItem := m.list.SelectedItem().(bmcategory)
	return tuicommands.SelectCategory(selectedItem.Name, switchView)
}

func (m Model) createBookmark() tea.Cmd {
	return func() tea.Msg {
		return tuicommands.CreateBookmarkMsg(true)
	}
}
