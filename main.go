package main

import (
	"log"
	"os"
	"runtime"
	"search-and-sort-movies/myapp"
	"time"
)

var (
	BuildVersion string
	BuildHash    string
	BuildDate    string
	BuildClean   string
	BuildName    = "search-and-sort-movies"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if _, err := os.Stat(myapp.FOLDER_CONFIG); os.IsNotExist(err) {
		os.Mkdir(myapp.FOLDER_CONFIG, os.ModePerm)
	}
}

func main() {
	myapp.Flags(BuildName, BuildVersion, BuildHash, BuildDate, BuildClean)
	// Write log to file : log_SearchAndSort
	f, err := os.OpenFile(myapp.LOGFILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	myapp.StartServerWeb()

	for {
		if myapp.GetEnv("dlna") != "" || myapp.GetEnv("movies") != "" || myapp.GetEnv("series") != "" {
			break
		}
		log.Println("En attente de configuration : va sur http://localhost:1515")
		time.Sleep(30 * time.Second)
	}

	checkFolderExists(myapp.GetEnv("dlna"))
	checkFolderExists(myapp.GetEnv("movies"))
	checkFolderExists(myapp.GetEnv("series"))

	log.Println("Start :-D")
	log.Println("Ecoute sur le dossier : " + myapp.GetEnv("dlna"))

	myapp.Watcher(myapp.GetEnv("dlna"))

}
func checkFolderExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}
}
