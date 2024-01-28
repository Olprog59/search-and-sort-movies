package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"search-and-sort-movies/myapp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.SetFlags(0)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	fmt.Println("path: " + constants.A_TRIER)
	if constants.ALL == "" {
		checkIfFolderExistAndCreate(constants.ALL)
		constants.A_TRIER = constants.ALL + "/a_trier"
		constants.MOVIES = constants.ALL + "/movies"
		constants.SERIES = constants.ALL + "/series"
	}
	checkIfFolderExistAndCreate(constants.A_TRIER)
	checkIfFolderExistAndCreate(constants.MOVIES)
	checkIfFolderExistAndCreate(constants.SERIES)
}

func main() {
	logger.L(logger.Magenta, "Start :-D")

	myapp.MyWatcher(constants.A_TRIER)
}

func checkIfFolderExistAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0766)
	}
}
