package bookmarks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/evcraddock/bm/pkg/utils"
	readability "github.com/philipjkim/goreadability"
	"github.com/spf13/viper"
)

// Bookmark bookmark data
type Bookmark struct {
	Name      string    `yaml:"title"`
	URL       string    `yaml:"url"`
	Author    string    `yaml:"author,omitempty"`
	Tags      []string  `yaml:"tags,omitempty"`
	DateAdded time.Time `yaml:"dateAdded,omitempty"`
	location  string
}

type ByName []Bookmark

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type ByDateAdded []Bookmark

func (b ByDateAdded) Len() int           { return len(b) }
func (b ByDateAdded) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDateAdded) Less(i, j int) bool { return b[i].DateAdded.Before(b[j].DateAdded) }

// BookmarkManager manages a bookmark
type BookmarkManager struct {
	bookmarkFolder string
	interactive    bool
	category       string
	location       string
}

// NewBookmarkManager creates new BookmarkManager
func NewBookmarkManager(interactive bool, category string) *BookmarkManager {
	bookmarkFolder := viper.GetString("BookmarkFolder")

	return &BookmarkManager{
		bookmarkFolder: bookmarkFolder,
		interactive:    interactive,
		category:       category,
		location:       fmt.Sprintf("%s/%s", bookmarkFolder, category),
	}
}

func (b *BookmarkManager) SetCategory(category string) {
	b.category = category
}

// GetBookmarkLocation return the folder location where bookmarks are stored
func (b *BookmarkManager) GetBookmarkLocation(title string) string {
	folderName := utils.ScrubFolder(title)
	bookmarkLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderName)
	return bookmarkLocation
}

// Load loads a Bookmark from a file
func (b *BookmarkManager) Load(bookmarkLocation string) (*Bookmark, error) {
	bookmarkfile, err := ioutil.ReadFile(bookmarkLocation + "/index.bm")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("bookmark %s does not exist", bookmarkLocation)
		}

		return nil, err
	}

	bookmark := &Bookmark{}
	err = yaml.Unmarshal(bookmarkfile, bookmark)
	bookmark.location = bookmarkLocation

	return bookmark, err
}

// LoadBookmarks returns a list of bookmarks
func (b *BookmarkManager) LoadBookmarks() ([]Bookmark, error) {
	bookmarkLocation := fmt.Sprintf("%s/%s", b.bookmarkFolder, b.category)
	var bookmarks []Bookmark

	err := filepath.Walk(bookmarkLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			bookmark, err := b.Load(path)
			if err == nil {
				bookmarks = append(bookmarks, *bookmark)
			}
		}

		return nil
	})

	sort.Sort(ByDateAdded(bookmarks))
	return bookmarks, err

}

// Save saves a bookmark
func (b *BookmarkManager) Save(bookmark *Bookmark) (bool, error) {
	data, err := yaml.Marshal(bookmark)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(bookmark.location+"/index.bm", data, 0644)
	if err != nil {
		return false, err
	}

	b.savePreview(bookmark.URL, bookmark.location)
	return true, nil
}

func (b *BookmarkManager) Upsert(bookmark *Bookmark) (bool, error) {
	if ok, location := b.createFolder(bookmark.Name); ok {
		bookmark.location = location
		return b.Save(bookmark)
	}

	folderName := utils.ScrubFolder(bookmark.Name)
	bookmarkLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderName)
	updated, err := b.Load(bookmarkLocation)
	if err != nil {
		return false, err
	}

	updated.URL = bookmark.URL
	updated.Author = bookmark.Author
	updated.Tags = bookmark.Tags

	return b.Save(updated)
}

// Create save a new bookmark
func (b *BookmarkManager) Create(url, title string) (bool, error) {
	if title == "" {
		opt := readability.NewOption()
		opt.ImageRequestTimeout = 3000
		content, err := readability.Extract(url, opt)
		if err != nil {
			return false, err
		}

		if content.Title != "" {
			title = content.Title
		}
	}

	bookmark := &Bookmark{
		Name:      title,
		URL:       url,
		DateAdded: time.Now(),
	}

	if b.interactive {
		bookmark = b.Edit(bookmark)
	}

	return b.Upsert(bookmark)
}

// Update updates and existing bookmark
func (b *BookmarkManager) UpdatePrompt(title string) (bool, error) {
	folderName := utils.ScrubFolder(title)
	bookmarkLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderName)

	bookmark, err := b.Load(bookmarkLocation)
	if err != nil {
		return false, err
	}

	newbookmark := b.Edit(bookmark)
	if title != newbookmark.Name {
		if ok, location := b.moveFolder(bookmark.location, newbookmark.Name); ok {
			newbookmark.location = location
		} else {
			return false, fmt.Errorf("Error moving folder %s", bookmark.location)
		}
	}

	return b.Save(newbookmark)
}

// Edit prompts for changes to a bookmark
func (b *BookmarkManager) Edit(bookmark *Bookmark) *Bookmark {
	fmt.Printf("To keep current value press return\n")
	title := utils.StringPrompt("Name", bookmark.Name)
	url := utils.StringPrompt("Link", bookmark.URL)
	author := utils.StringPrompt("Author", bookmark.Author)
	tags := utils.ListPrompt("Tags", strings.Join(bookmark.Tags, ","))

	newBookMark := &Bookmark{
		Name:     title,
		URL:      url,
		Author:   author,
		Tags:     tags,
		location: bookmark.location,
	}

	return newBookMark
}

// Remove removes a bookmark
func (b *BookmarkManager) Remove(title string) error {
	if title == "" {
		// nothing to remove
		// already have desired result; return nil
		return nil
	}

	folderName := utils.ScrubFolder(title)
	bookmarklocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderName)

	return utils.RemoveFolder(bookmarklocation)
}

func (b *BookmarkManager) createFolder(title string) (bool, string) {
	folderName := utils.ScrubFolder(title)
	saveLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderName)

	if ok := utils.CreateFolder(saveLocation); !ok {
		return false, saveLocation
	}

	return true, saveLocation
}

func (b *BookmarkManager) moveFolder(oldlocation, newtitle string) (bool, string) {
	err := utils.RemoveFolder(oldlocation)
	if err != nil {
		return false, oldlocation
	}

	return b.createFolder(newtitle)
}
