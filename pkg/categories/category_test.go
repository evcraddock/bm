package categories

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/spf13/viper"
)

func TestGetCategoryList(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Fail()
	}

	homeDirectory := usr.HomeDir
	viper.Set("bookmarkFolder", fmt.Sprintf("%s/.local/bookmarks", homeDirectory))
	manager := NewCategoryManager()

	categories, err := manager.GetCategoryList()
	if err != nil {
		t.Logf("error loading categories: %v", err)
		t.Fail()
	}

	if len(categories) == 0 {
		t.Log("should be at least one category")
		t.Fail()
	}
}

func TestGetCategoryListFail(t *testing.T) {
	viper.Set("bookmarkFolder", "bad-location")
	manager := NewCategoryManager()

	categories, err := manager.GetCategoryList()
	if err == nil {
		t.Log("should be error loading categories")
		t.Fail()
	}

	if categories != nil {
		t.Log("should be no categories")
		t.Fail()
	}
}
