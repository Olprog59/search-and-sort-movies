package myapp

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
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

func ReadFileLog() (data []string) {
	file, err := os.Open("./log_SearchAndSort")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}
