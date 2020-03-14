package myapp

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"search-and-sort-movies/myapp/constants"
	"syscall"
)

func executeUpdate() {
	folder, err := os.Getwd()

	binary, lookErr := exec.LookPath("bash")
	if lookErr != nil {
		panic(lookErr)
	}
	err = syscall.Exec(binary, []string{binary, "-c", filepath.Clean(folder + "/" + constants.FileUpdateName)}, os.Environ())
	if err != nil {
		log.Println(err.Error())
	}
}
