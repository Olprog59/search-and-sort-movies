package myapp

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const regexFile = `(.mkv|.mp4|.avi|.flv)`

func startScan() {
	if count, file := fileInFolder(); count > 0 {
		go boucleFiles(file)
	}
}

func fileInFolder() (int, []string) {
	var files []string
	var count int
	err := filepath.Walk(GetEnv("dlna"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		re := regexp.MustCompile(regexFile)
		if re.MatchString(filepath.Ext(path)) {
			files = append(files, path)
			count++
		}
		return nil
	})
	if err != nil {
		log.Printf("walk error [%v]\n", err)
	}
	return count, files
}

func boucleFiles(files []string) {
	log.Println("Démarrage du tri !")
	for _, f := range files {
		log.Println("File : " + f)
		Process(f)
	}
	log.Println("Tri terminé !")
}
