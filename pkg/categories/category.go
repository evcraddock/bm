package categories

import (
	"io/ioutil"

	"github.com/evcraddock/bm/pkg/config"
)

type Category struct {
	Name        string
	Description string
}

type CategoryManager struct {
	bookmarkFolder string
}

func NewCategoryManager(cfg *config.Config) *CategoryManager {
	return &CategoryManager{
		bookmarkFolder: cfg.BookmarkFolder,
	}
}

func (c *CategoryManager) GetCategoryList() ([]Category, error) {
	var categoryList []Category
	folders, err := ioutil.ReadDir(c.bookmarkFolder)
	if err != nil {
		return nil, err
	}

	for _, folder := range folders {
		cat := Category{
			Name:        folder.Name(),
			Description: "",
		}

		categoryList = append(categoryList, cat)
	}

	return categoryList, nil
}
