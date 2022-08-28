package bookmarkstui

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	tuicommands "github.com/evcraddock/bm/cmd/bm/tui/commands"
	"github.com/evcraddock/bm/pkg/bookmarks"
)

var (
	docStyle     = lipgloss.NewStyle().Margin(1, 1)
	marginHeight = 4
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type Model struct {
	category   string
	list       list.Model
	manager    *bookmarks.BookmarkManager
	windowSize *tea.WindowSizeMsg
}

type item struct {
	name, description string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.name }

func New(category string, windowSize *tea.WindowSizeMsg) Model {
	bookmarkModel := Model{
		manager:    bookmarks.NewBookmarkManager(false, category),
		category:   category,
		windowSize: windowSize,
	}

	m := bookmarkModel.loadBookmarksList(category)
	m.Title = category

	if windowSize != nil {
		width, height := bookmarkModel.getWindowSize()
		m.SetSize(width, height)
	}

	bookmarkModel.list = m
	return bookmarkModel
}

func (b Model) Init() tea.Cmd {
	return nil
}

func (b Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc", "shift+tab":
			return b, func() tea.Msg {
				return tuicommands.CategoryViewMsg(true)
			}

		case "ctrl+c":
			return b, tea.Quit

		case "ctrl+o", "enter":
			b.openSelectedUrl()

		case "ctrl+r":
			return b, func() tea.Msg {
				return tuicommands.ReloadBookmarksMsg(true)
			}

		case "ctrl+d":
			b.deleteBookmark()
			return b, func() tea.Msg {
				return tuicommands.ReloadBookmarksMsg(true)
			}

		case "ctrl+e":
			selectedItem := b.list.SelectedItem().(item)
			bookmark, err := b.manager.Load(b.manager.GetBookmarkLocation(selectedItem.Title()))
			if err != nil {
				// TODO return error msg instead of panic
				panic(err)
			}
			return b, tuicommands.SelectBookmark(bookmark)

		default:
			b.list, cmd = b.list.Update(msg)

		}

	case tea.WindowSizeMsg:
		b.windowSize = &msg
		b.list.SetSize(b.getWindowSize())
	}

	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b Model) View() string {
	return docStyle.Render(b.list.View())
}

func (b Model) getWindowSize() (int, int) {
	return b.windowSize.Width, b.windowSize.Height - marginHeight
}

func (b Model) getSelectedBookmark() tea.Msg {
	selectedItem := b.list.SelectedItem().(item)
	bookmark, err := b.manager.Load(b.manager.GetBookmarkLocation(selectedItem.Title()))
	if err != nil {
		return errMsg{err}
	}

	return tuicommands.SelectBookmark(bookmark)
}

func (b Model) deleteBookmark() tea.Msg {
	bookmark := b.list.SelectedItem().(item)
	if bookmark.Title() != "" {
		err := b.manager.Remove(bookmark.Title())
		if err != nil {
			return errMsg{err}
		}
	}

	b.list.RemoveItem(b.list.Index())

	return nil
}

func (b Model) loadBookmarksList(category string) list.Model {
	l, err := b.manager.LoadBookmarks()
	if err != nil {
		panic(err)
	}

	items := []list.Item{}
	for _, bm := range l {
		items = append(items, item{name: bm.Title, description: bm.URL})
	}

	m := list.New(items, list.NewDefaultDelegate(), 0, 0)

	return m
}

func (b Model) openSelectedUrl() tea.Msg {
	selectedItem := b.list.SelectedItem().(item)
	cmd := exec.Command("xdg-open", selectedItem.description)
	err := cmd.Start()
	if err != nil {
		return errMsg{err}
	}

	return nil
}
