package myapp

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func MyWatcher(location string) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()

	done := make(chan bool)

	go watchStart(watch)

	err = watch.Add(location)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func watchStart(watch *fsnotify.Watcher) {
	var folder string
	for {
		select {
		case event := <-watch.Events:
			fmt.Println(event)


			if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println(event.Name)
				re := regexp.MustCompile(regexFile)
				if !re.MatchString(filepath.Ext(event.Name)) {
					log.Println("Détection de : ", event.Name)
					os.RemoveAll(filepath.Dir(event.Name))
				}
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				f, _ := os.Stat(event.Name)
				if f.IsDir() {
					folder = event.Name
					fmt.Println("ecoute sur : " + folder)
					MyWatcher(folder)
				}
				_, file := filepath.Split(event.Name)
				re := regexp.MustCompile(regexFile)
				if !re.MatchString(filepath.Ext(file)) {
					log.Println("Détection de : ", file)
					//Process(event.Name)
				}
			}
		case err := <-watch.Errors:
			log.Println("error:", err)
		}
	}
}
