package myapp

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func executeUpdate() {
	folder, err := os.Getwd()

	binary, lookErr := exec.LookPath("sh")
	if lookErr != nil {
		panic(lookErr)
	}
	err = syscall.Exec(binary, []string{binary, "-c", filepath.Clean(folder + "/" + FileUpdateName)}, os.Environ())
	if err != nil {
		log.Println(err.Error())
	}
}
