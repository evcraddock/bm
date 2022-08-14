package bookmarks

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/spf13/viper"
)

func TestNewBookmarkManager(t *testing.T) {
	viper.Set("bookmarkfolder", "test-path")

	manager := NewBookmarkManager(false, "test")
	if manager.category != "test" {
		t.Log("category should be test")
		t.Fail()
	}

	if manager.interactive {
		t.Log("interactive should be false")
		t.Fail()
	}

	if manager.bookmarkFolder != "test-path" {
		t.Logf("bookmarkFolder should be %s", manager.bookmarkFolder)
		t.Fail()
	}

	if manager.location != "test-path/test" {
		t.Logf("location should be %s", manager.location)
		t.Fail()
	}
}

func TestSetCategory(t *testing.T) {
	manager := NewBookmarkManager(false, "test")

	manager.SetCategory("test-cat")

	if manager.category != "test-cat" {
		t.Log("category is incorrect")
		t.Fail()
	}
}

func TestGetBookmarkLocation(t *testing.T) {
	viper.Set("bookmarkFolder", "test-path")
	expected := "test-path/test/title"
	manager := NewBookmarkManager(false, "test")

	location := manager.GetBookmarkLocation("title")

	if location != expected {
		t.Logf("location should be %s", location)
		t.Fail()
	}
}

func TestLoadFailed(t *testing.T) {
	viper.Set("bookmarkFolder", "no-exists")
	location := "bad-location"
	manager := NewBookmarkManager(false, "test")

	bm, err := manager.Load(location)
	if err == nil {
		t.Log("there should be an error")
		t.Fail()
	}

	if bm != nil {
		t.Log("bookmark should be nil")
		t.Fail()
	}
}

func TestLoadBookmarks(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Fail()
	}

	homeDirectory := usr.HomeDir
	viper.Set("bookmarkFolder", fmt.Sprintf("%s/.local/bookmarks", homeDirectory))
	manager := NewBookmarkManager(false, "readlater")

	bookmarks, err := manager.LoadBookmarks()
	if err != nil {
		t.Logf("error loading: %v", err)
		t.Fail()
	}

	if len(bookmarks) == 0 {
		t.Log("bookmarks should not be an empty list")
		t.Fail()
	}
}
