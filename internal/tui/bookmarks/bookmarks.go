package bookmarkstui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	tuicommands "github.com/evcraddock/bm/internal/tui/commands"
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

type bookmark struct {
	bookmarks.Bookmark
}

func (i bookmark) Title() string {
	var b strings.Builder
	b.WriteString(i.Name)

	if i.Author != "" {
		b.WriteString(fmt.Sprintf(" (by %s)", i.Author))
	}

	return b.String()
}

func (i bookmark) Description() string {
	var b strings.Builder

	for _, tag := range i.Tags {
		b.WriteString(fmt.Sprintf("[%s] ", strings.Trim(tag, " ")))
	}

	b.WriteString(i.URL)
	return b.String()
}

func (i bookmark) FilterValue() string {
	return i.Name
}

func New(category, selectedBookmark string, selectedIndex int, windowSize *tea.WindowSizeMsg) Model {
	bookmarkModel := Model{
		manager:    bookmarks.NewBookmarkManager(false, category),
		category:   category,
		windowSize: windowSize,
	}

	m := bookmarkModel.loadBookmarksList(category, selectedBookmark, selectedIndex)
	m.Title = category
	m.SetShowFilter(false)
	m.SetStatusBarItemName("bookmark", "bookmarks")
	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(len(m.Items()) > 0)

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

		case "ctrl+o":
			b.openSelectedUrl()

		case "ctrl+r":
			return b, tuicommands.ReloadBookmarks(b.list.Index())

		case "ctrl+d":
			return b, b.deleteBookmark()

		case "ctrl+e":
			return b, b.getSelectedBookmark()

		case "ctrl+n":
			return b, func() tea.Msg {
				return tuicommands.CreateBookmarkMsg(true)
			}

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
	w, h := 0, 0
	if b.windowSize != nil {
		w = b.windowSize.Width
		h = b.windowSize.Height - marginHeight
	}

	return w, h
}

func (b Model) getSelectedBookmark() tea.Cmd {
	selectedItem := b.list.SelectedItem().(bookmark)
	bookmark, err := b.manager.Load(b.manager.GetBookmarkLocation(selectedItem.Name))
	if err != nil {
		return func() tea.Msg {
			return errMsg{err}
		}
	}

	return tuicommands.SelectBookmark(bookmark)
}

func (b Model) deleteBookmark() tea.Cmd {
	bookmark := b.list.SelectedItem().(bookmark)
	if bookmark.Name != "" {
		err := b.manager.Remove(bookmark.Name)
		if err != nil {
			return func() tea.Msg {
				return errMsg{err}
			}
		}
	}

	index := b.list.Index()
	lastIndex := len(b.list.Items()) - 1
	b.list.RemoveItem(index)
	if index > 0 && lastIndex == index {
		index--
	}

	return tuicommands.ReloadBookmarks(index)
}

func (b Model) loadBookmarksList(category, bookmarkName string, selectedIndex int) list.Model {
	l, err := b.manager.LoadBookmarks()
	if err != nil {
		panic(err)
	}

	index := selectedIndex
	items := []list.Item{}
	for i := 0; i < len(l); i++ {
		bm := l[i]
		items = append(items, bookmark{bm})
		if selectedIndex > 0 && bm.Name == bookmarkName {
			index = i
		}
	}

	// w, h := 0, 0
	// if b.windowSize != nil {
	w, h := b.getWindowSize()
	// }

	m := list.New(items, list.NewDefaultDelegate(), w, h)
	m.Select(index)
	return m
}

func (b Model) openSelectedUrl() tea.Msg {
	bookmark := b.list.SelectedItem().(bookmark)
	cmd := exec.Command("xdg-open", bookmark.URL)
	err := cmd.Start()
	if err != nil {
		return errMsg{err}
	}

	return nil
}
