package myapp

import (
	"log"

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
					log.Println("Détection du fichier : ", event.Name)
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
