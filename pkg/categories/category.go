package categories

import (
	"io/ioutil"

	"github.com/spf13/viper"
)

type Category struct {
	Name        string
	Description string
}

type CategoryManager struct {
	bookmarkFolder string
}

func NewCategoryManager() *CategoryManager {
	bookmarkFolder := viper.GetString("BookmarkFolder")
	return &CategoryManager{
		bookmarkFolder: bookmarkFolder,
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
