package bookmarks

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/evcraddock/bm/pkg/config"
	"github.com/evcraddock/bm/pkg/utils"
)

type Bookmark struct {
	URL    string   `yaml:"url"`
	Title  string   `yaml:"title"`
	Author string   `yaml:"author"`
	Tags   []string `yaml:"tags"`
}

type Bookmarks []Bookmark

type BookmarkManager struct {
	bookmarkFolder string
	interactive    bool
}

func NewBookmarkManager(cfg *config.Config, interactive bool) *BookmarkManager {
	return &BookmarkManager{
		bookmarkFolder: cfg.BookmarkFolder,
	}
}

func (b *BookmarkManager) Save(bookmark Bookmark, category string) error {
	folderTitle := utils.ScrubFolder(bookmark.Title)
	saveLocation := fmt.Sprintf("%s/%s/%s", b.bookmarkFolder, category, folderTitle)

	if ok := utils.CreateFolder(saveLocation); !ok {
		return fmt.Errorf("Bookmark %s exists \n", bookmark.Title)
		// bookmark = b.Edit(bookmark, category)

	}

	data, err := yaml.Marshal(bookmark)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(saveLocation+"/index.bm", data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Saved Bookmark %s\n", bookmark.Title)
	return nil
}

func (b *BookmarkManager) Edit(bookmark Bookmark, category string) Bookmark {
	bookmark.Author = utils.StringPrompt("Author")
	bookmark.Tags = utils.ListPrompt("Tags")

	//	return b.Save(bookmark, category)
	return bookmark
}
