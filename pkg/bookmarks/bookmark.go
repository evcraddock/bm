package bookmarks

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/evcraddock/bm/pkg/config"
	"github.com/evcraddock/bm/pkg/utils"
)

type Bookmark struct {
	URL      string   `yaml:"url"`
	Title    string   `yaml:"title"`
	Author   string   `yaml:"author"`
	Tags     []string `yaml:"tags"`
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

func (b *BookmarkManager) Load(title string) (*Bookmark, error) {
	folderTitle := utils.ScrubFolder(title)
	bookmarkLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderTitle)

	bookmarkfile, err := ioutil.ReadFile(bookmarkLocation + "/index.bm")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("The bookmark %s does not exist", title)
		}

		return nil, err
	}

	bookmark := &Bookmark{}
	err = yaml.Unmarshal(bookmarkfile, bookmark)
	bookmark.location = bookmarkLocation

	return bookmark, err
}

func (b *BookmarkManager) Save(bookmark *Bookmark) error {
	// folderTitle := utils.ScrubFolder(bookmark.Title)
	// saveLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, b.category, folderTitle)

	// if ok := utils.CreateFolder(saveLocation); !ok {
	// 	// return fmt.Errorf("Bookmark %s exists \n", bookmark.Title)
	// }

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
	bookmark, err := b.Load(title)
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
