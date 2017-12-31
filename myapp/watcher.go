package myapp

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Watcher(location string) {
	Watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer Watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-Watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					_, file := filepath.Split(event.Name)
					log.Println("Détection de : ", file)
					Process(event.Name)
				}
			case err := <-Watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = Watcher.Add(location)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
