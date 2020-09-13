package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"search-and-sort-movies/myapp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/flags"
	"search-and-sort-movies/myapp/update"
)

var (
	BuildVersion string
	BuildHash    string
	BuildDate    string
	BuildClean   string
	BuildName    = "search-and-sort-movies-" + runtime.GOOS + "-" + runtime.GOARCH
	//count        int
	//ticker       *time.Ticker
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	//isFlags := make(chan bool)
	flags.Flags(BuildName, BuildVersion, BuildHash, BuildDate, BuildClean)
	//<- isFlags

	// Write log to file : log_SearchAndSort
	f, err := os.OpenFile(constants.LOGFILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("\n\nBuild Version: %s\nBuild Date: %s\n\n", BuildVersion, BuildDate)
	// Start test update application auto

	// TODO ne pas oublier d'activer pour l'auto update
	var applicationUpdate update.Application
	go applicationUpdate.LaunchAppCheckUpdate(BuildVersion, BuildName)

	// End test update application auto

	// Check if it's the first connection
	if firstConnect() {
		firstConfig()
	} else {
		//for {
		if myapp.GetEnv("dlna") == "" || myapp.GetEnv("movies") == "" || myapp.GetEnv("series") == "" || myapp.GetEnv("port") == "" {
			firstConfig()
		}
		//}
	}

	checkFolderExists(myapp.GetEnv("dlna"))
	checkFolderExists(myapp.GetEnv("movies"))
	checkFolderExists(myapp.GetEnv("series"))

	log.Println("Start :-D")
	log.Println("Ecoute sur le dossier : " + myapp.GetEnv("dlna"))
	fmt.Println("Start :-D")
	fmt.Println("Ecoute sur le dossier : " + myapp.GetEnv("dlna"))

	go myapp.ServerHttp()

	myapp.MyWatcher(myapp.GetEnv("dlna"))

}

func checkFolderExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, os.ModePerm)
	}
}

func firstConnect() bool {
	_, err := os.Stat(constants.ConfigFile)

	if os.IsNotExist(err) {
		log.Println(err)
		return true
	}
	return false
}

func readJSONFileConsole() {
	f, err := ioutil.ReadFile(constants.ConfigFile)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%s\n", string(f))
}

func firstConfig() {
	fmt.Printf("Click to open web site : http://%s:%s/config\n", myapp.IpLocal(), myapp.GetEnv("port"))
}
