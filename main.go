package main

import (
	"fmt"
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
}

func main() {
	flags.Flags(BuildName, BuildVersion, BuildHash, BuildDate, BuildClean)

	fmt.Printf("\n\nBuild Version: %s\nBuild Date: %s\n\n", BuildVersion, BuildDate)

	checkIfFolderExistAndCreate(constants.A_TRIER)
	checkIfFolderExistAndCreate(constants.MOVIES)
	checkIfFolderExistAndCreate(constants.SERIES)

	fmt.Println(logger.Info("Start :-D"))
	fmt.Println(logger.Info("Ecoute sur le dossier : " + constants.A_TRIER))

	myapp.MyWatcher(constants.A_TRIER)

}

func checkIfFolderExistAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0766)
	}
}
