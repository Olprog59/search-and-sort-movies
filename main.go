package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"search-and-sort-movies/myapp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/flags"
	"search-and-sort-movies/myapp/logger"
)

var (
	BuildVersion string
	BuildHash    string
	BuildDate    string
	BuildClean   string
	BuildName    = "search-and-sort-movies-" + runtime.GOOS + "-" + runtime.GOARCH
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(0)
}

func main() {
	flags.Flags(BuildName, BuildVersion, BuildHash, BuildDate, BuildClean)

	// Write log to file : log_SearchAndSort
	f, err := os.OpenFile(constants.LOGFILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.L(logger.Red, "error opening file: %s", err)
	}
	defer f.Close()
	log.SetOutput(f)

	logger.L(logger.Teal, "\n\nBuild Version: %s\nBuild Date: %s\n", BuildVersion, BuildDate)

	checkIfFolderExistAndCreate(constants.A_TRIER)
	checkIfFolderExistAndCreate(constants.MOVIES)
	checkIfFolderExistAndCreate(constants.SERIES)

	logger.L(logger.Magenta, "Start :-D")
	fmt.Print(logger.Magenta("Start :-D"))
	fmt.Print(logger.Magenta("Ecoute sur le dossier : " + constants.A_TRIER))

	myapp.MyWatcher(constants.A_TRIER)

}

func checkIfFolderExistAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0766)
	}
}
