package bookmarks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/evcraddock/bm/pkg/config"
	"github.com/evcraddock/bm/pkg/utils"
)

type Bookmark struct {
	Title    string   `yaml:"title"`
	URL      string   `yaml:"url"`
	Author   string   `yaml:"author,omitempty"`
	Tags     []string `yaml:"tags,omitempty"`
	location string
}

type Bookmarks []Bookmark

type BookmarkManager struct {
	bookmarkFolder string
	interactive    bool
	category       string
}

func NewBookmarkManager(cfg *config.Config, interactive bool, category string) *BookmarkManager {
	return &BookmarkManager{
		bookmarkFolder: cfg.BookmarkFolder,
		interactive:    interactive,
		category:       category,
	}
}

func (b *BookmarkManager) Load(bookmarkLocation string) (*Bookmark, error) {
	bookmarkfile, err := ioutil.ReadFile(bookmarkLocation + "/index.bm")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("The bookmark %s does not exist", bookmarkLocation)
		}

		return nil, err
	}

	bookmark := &Bookmark{}
	err = yaml.Unmarshal(bookmarkfile, bookmark)
	bookmark.location = bookmarkLocation

	return bookmark, err
}

func (b *BookmarkManager) LoadBookmarks() ([]*Bookmark, error) {
	bookmarkLocation := fmt.Sprintf("%s/%s", b.bookmarkFolder, b.category)
	var bookmarks []*Bookmark

	_ = filepath.Walk(bookmarkLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if bookmark, err := b.Load(path); err == nil {
				bookmarks = append(bookmarks, bookmark)
			}
		}

		return nil
	})

	return bookmarks, fmt.Errorf("%s does not exist\n", bookmarkLocation)

}

func (b *BookmarkManager) Save(bookmark *Bookmark) error {
	data, err := yaml.Marshal(bookmark)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(bookmark.location+"/index.bm", data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Saved Bookmark %s\n", bookmark.Title)
	return nil
}

func (b *BookmarkManager) Create(title, url string) error {
	bookmark := &Bookmark{
		Title: title,
		URL:   url,
	}

	if b.interactive {
		bookmark = b.Edit(bookmark)
	}

	if ok, location := b.createFolder(bookmark.Title); ok {
		bookmark.location = location
		if err := b.Save(bookmark); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Bookmark %s already exists \n", title)
}

func (b *BookmarkManager) Update(title string) error {
	folderTitle := utils.ScrubFolder(title)
	bookmarkLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderTitle)

	bookmark, err := b.Load(bookmarkLocation)
	if err != nil {
		return err
	}

	newbookmark := b.Edit(bookmark)
	if title != newbookmark.Title {
		if ok, location := b.moveFolder(bookmark.location, newbookmark.Title); ok {
			newbookmark.location = location
		} else {
			return fmt.Errorf("Error moving folder %s", bookmark.location)
		}
	}

	if err := b.Save(newbookmark); err != nil {
		return err
	}

	return nil
}

func (b *BookmarkManager) Edit(bookmark *Bookmark) *Bookmark {
	fmt.Printf("To keep current value press return\n")
	title := utils.StringPrompt("Title", bookmark.Title)
	url := utils.StringPrompt("Link", bookmark.URL)
	author := utils.StringPrompt("Author", bookmark.Author)
	tags := utils.ListPrompt("Tags", strings.Join(bookmark.Tags, ","))

	newBookMark := &Bookmark{
		Title:    title,
		URL:      url,
		Author:   author,
		Tags:     tags,
		location: bookmark.location,
	}

	return newBookMark
}

func (b *BookmarkManager) Remove(title string) error {
	folderTitle := utils.ScrubFolder(title)
	bookmarklocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderTitle)

	return utils.RemoveFolder(bookmarklocation)
}

func (b *BookmarkManager) createFolder(title string) (bool, string) {
	folderTitle := utils.ScrubFolder(title)
	saveLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderTitle)

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
