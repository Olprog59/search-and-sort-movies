package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"search-and-sort-movies/myapp"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

var (
	buildVersion string
	buildHash    string
	buildDate    string
	buildClean   string
	buildName    = "search-and-sort-movies"
)

func main() {

	vers := flag.Bool("v", false, "Indique la version de l'application")
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	flag.Parse()

	if *vers {
		// flag.PrintDefaults()
		fmt.Printf("Name: %s\n", buildName)
		fmt.Printf("Version: %s\n", buildVersion)
		fmt.Printf("Git Commit Hash: %s\n", buildHash)
		fmt.Printf("Build Date: %s\n", buildDate)
		fmt.Printf("Built from clean source tree: %s\n", buildClean)
		fmt.Printf("OS: %s\n", runtime.GOOS)
		fmt.Printf("Architecture: %s\n", runtime.GOARCH)
		os.Exit(1)
	}

	if *scan {
		myapp.StartScan(true)
	}

	// Write log to file : log_SearchAndSort
	f, err := os.OpenFile("log_SearchAndSort", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Check if it's the first connection
	if myapp.FirstConnect() {
		myapp.FirstConfig()
	} else {
		for {
			if myapp.GetEnv("dlna") == "" || myapp.GetEnv("movies") == "" || myapp.GetEnv("series") == "" {
				myapp.FirstConfig()
			} else {
				break
			}
		}
	}

	myapp.CheckFolderExists(myapp.GetEnv("dlna"))
	myapp.CheckFolderExists(myapp.GetEnv("movies"))
	myapp.CheckFolderExists(myapp.GetEnv("series"))

	fmt.Println("Start :-D")

	// startScan(false)

	fmt.Println("Ecoute sur le dossier : " + myapp.GetEnv("dlna"))
	myapp.Watcher(myapp.GetEnv("dlna"))

}
