package app

import (
	//	"fmt"
	"os/exec"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

// BookmarkApp text ui for managing bookmarks
type BookmarkApp struct {
	app              *tview.Application
	manager          *bookmarks.BookmarkManager
	selectedIndex    int
	selectedBookmark *bookmarks.Bookmark
}

// NewBookmarkApp creates a new bookmark app
func NewBookmarkApp(manager *bookmarks.BookmarkManager) *BookmarkApp {
	app := tview.NewApplication()

	return &BookmarkApp{
		app:           app,
		manager:       manager,
		selectedIndex: 0,
	}
}

//Load loads a list of bookmarks
func (b *BookmarkApp) Load() {
	b.draw()
	b.app.SetInputCapture(b.handleInput)
	if err := b.app.Run(); err != nil {
		panic(err)
	}
}

func (b *BookmarkApp) draw() {
	grid := b.createGrid()
	main := b.createTable()

	b.loadLinks(main)
	if main.GetRowCount() > 0 {
		main.Select(b.selectedIndex, 0)
	}

	grid.AddItem(main, 1, 0, 1, 3, 0, 0, false)

	b.app.SetRoot(grid, true)
	b.app.SetFocus(main)
}

func (b *BookmarkApp) createGrid() *tview.Grid {
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0).
		SetBorders(true).
		AddItem(b.createHeader(), 0, 0, 1, 3, 0, 0, false).
		AddItem(b.createFooter(), 2, 0, 1, 3, 0, 0, false)

	return grid

}

func (b *BookmarkApp) createTable() *tview.Table {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.Select(b.selectedIndex, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				b.app.Stop()
			}
		}).
		SetSelectedFunc(func(row int, column int) {
			b.openLink(b.selectedBookmark)
		})

	return table
}

func (b *BookmarkApp) createHeader() *tview.TextView {
	header := tview.NewTextView()
	header.SetText("bookmark app ver:0.1")

	return header
}

func (b *BookmarkApp) createFooter() *tview.TextView {
	footer := tview.NewTextView()
	footer.SetText("ENTER:Open q:Quit R:Reload d:Delete")

	return footer
}

func (b *BookmarkApp) loadLinks(table *tview.Table) {
	items, err := b.manager.LoadBookmarks()
	if err != nil {
		panic(err)
	}

	for i, bookmark := range items {
		table.SetCell(i, 0, tview.NewTableCell(bookmark.Title).SetAlign(tview.AlignLeft))
		table.SetCell(i, 1, tview.NewTableCell(bookmark.URL).SetAlign(tview.AlignLeft).SetMaxWidth(0))
	}

	table.SetSelectionChangedFunc(func(row, column int) {
		b.selectedBookmark = items[row]
		rowcount := table.GetRowCount() - 1
		if rowcount == 0 {
			b.selectedIndex = 0
			return
		}

		if rowcount-row == 0 {
			b.selectedIndex = row - 1
		} else {
			b.selectedIndex = row
		}

	})

}

func (b *BookmarkApp) openLink(bookmark *bookmarks.Bookmark) {
	cmd := exec.Command("xdg-open", bookmark.URL)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func (b *BookmarkApp) deleteBookmark() {
	b.manager.Remove(b.selectedBookmark.Title)
}

func (b *BookmarkApp) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key == tcell.KeyRune {
		switch ch := event.Rune(); ch {
		case 'q':
			b.app.Stop()
		case 'R':
			b.draw()
		case 'd':
			b.deleteBookmark()
			b.draw()
		}
	}

	return event
}
