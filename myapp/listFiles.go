package myapp

import (
	"os"
	"path/filepath"
)

func ReadAllFiles() []string {
	var files []string

	root := GetEnv("dlna")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
