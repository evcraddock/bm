package utils

import (
	"os"
	"path"
	"strings"
)

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
