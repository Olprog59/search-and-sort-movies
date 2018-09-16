package myapp

import (
	"io/ioutil"
	"log"
)

func ReadAllFiles() []string {

	root := GetEnv("dlna")
	f, err := ioutil.ReadDir(root)
	if err != nil {
		log.Println(err)
	}
	var files []string
	for _, v := range f {
		if v.IsDir() {
			continue
		}
		files = append(files, v.Name())
	}
	return files
}
