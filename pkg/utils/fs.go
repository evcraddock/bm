package utils

import (
	"os"
	"path"
	"strings"
)

// CreateFolder creates a folder, returns a boolean for if the folder already existed
func CreateFolder(folder string) bool {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		if err != nil {
			panic(err)
		}

		return true
	}

	return false
}

// CreateFile creates a file, return a boolean for if the file was successfully created
func CreateFile(folder, filename string) bool {
	if folderexists := CreateFolder(folder); !folderexists {
		return folderexists
	}

	_, err := os.Create(path.Join(folder, filename))
	if err != nil {
		return false
	}

	return true
}

// RemoveFolder recursively removes a folder and its contents
func RemoveFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return err
	}

	return os.RemoveAll(folder)
}

// ScrubFolder removes special characters from a folder name
func ScrubFolder(folder string) string {

	r := strings.NewReplacer(
		" ", "-",
		",", "",
		"%", "",
		"(", "",
		")", "",
		"/", "-",
		"[", "",
		"]", "",
	)

	return r.Replace(folder)
}
