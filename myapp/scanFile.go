package myapp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func startScan(auto bool) {
	if count, file := fileInFolder(); count > 0 {
		if auto {
			fmt.Println("Scan automatique")
			go boucleFiles(file)
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Je vois qu'il y a des fichiers vidéos actuellement dans ton dossier source.")
			fmt.Println("Veux tu faire le tri? (O/n)")
			text, _ := reader.ReadString('\n')
			fmt.Println(text)
			if strings.TrimSpace(text) == "n" || strings.TrimSpace(text) == "N" {
				return
			}

			go boucleFiles(file)
		}
	}
}

func fileInFolder() (int, []os.FileInfo) {
	files, err := ioutil.ReadDir(GetEnv("dlna"))
	if err != nil {
		log.Fatal(err)
	}

	var count int
	for _, f := range files {
		if !f.IsDir() {
			re := regexp.MustCompile(`(.mkv|.mp4|.avi|.flv)`)
			if re.MatchString(filepath.Ext(f.Name())) {
				count++
			}
		}
	}
	return count, files
}

func boucleFiles(files []os.FileInfo) {
	log.Println("Démarrage du tri !")
	for _, f := range files {
		if !f.IsDir() {
			log.Println("File : " + f.Name())
			Process(f.Name())
		}
	}
	log.Println("Tri terminé !")
}
