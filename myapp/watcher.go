package myapp

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var watch *fsnotify.Watcher
var err error

func MyWatcher(location string) {
	watch, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watch.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					re := regexp.MustCompile(regexFile)
					if !_checkIfDir(event) {
						if re.MatchString(filepath.Ext(event.Name)) {
							go _fsNotifyCreateFile(event, re)
						}
					}

				}

			case err := <-watch.Errors:
				log.Println("error:", err)
			}
		}
	}()

	if err := watch.Add(location); err != nil {
		log.Fatal(err)
	}

	<-done
}

func _ticker(event fsnotify.Event, c *chan bool) {
	ticker := time.NewTicker(5 * time.Second)
	var size int64
	go func() {
		for range ticker.C {
			f, err := os.Stat(event.Name)
			if err != nil {
				log.Println(err)
			}
			//log.Printf("Name: %s\n\tInfo size: %d - Size: %d\n\n", event.Name, f.Size(), size)
			if f.Size() != size || f.Size() < 100000 {
				size = f.Size()
				continue
			}
			log.Println()
			ticker.Stop()
			*c <- true
		}
	}()
}

func _stat(event fsnotify.Event) (os.FileInfo, fsnotify.Event) {
	f, err := os.Stat(event.Name)
	if err != nil {
		log.Println(err)
	}
	return f, event
}

func _checkIfDir(event fsnotify.Event) bool {
	f, e := _stat(event)
	if f.IsDir() {
		err := watch.Add(e.Name)
		if err != nil {
			log.Println(err)
		}
		return false
	}
	return false
}

func _fsNotifyCreateFile(event fsnotify.Event, re *regexp.Regexp) {
	f, _ := _stat(event)

	finish := make(chan bool)
	go _ticker(event, &finish)
	<-finish

	if re.MatchString(filepath.Ext(f.Name())) {
		log.Println("Détection de : ", filepath.Base(f.Name()))
		folder := filepath.Dir(f.Name())
		if folder != GetEnv("dlna") {
			files, _ := ioutil.ReadDir(filepath.Dir(f.Name()))
			if len(files) == 1 {
				err := watch.Remove(f.Name())
				if err != nil {
					log.Println(err)
				}
			}
		}
		Process(event.Name)
	}

}
