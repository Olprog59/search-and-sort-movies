package myapp

import (
	"github.com/Olprog59/fsnotify"
	"os"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
	"time"
)

var watch *fsnotify.Watcher
var err error

func MyWatcher(location string) {
	watch, err = fsnotify.NewWatcher()
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
	//defer watch.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watch.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.MovedTo == fsnotify.MovedTo {
					//logger.L(logger.Red, "%s", event.Op)
					re := regexp.MustCompile(constants.RegexFile)
					//Ajout d'une sécurité si le fichier a déjà été déplacé
					if isDir, isNil := _checkIfDir(event); !isDir && !isNil {
						if re.MatchString(filepath.Ext(event.Name)) {
							go _fsNotifyCreateFile(event, re)
						}
					}
				}

			case err := <-watch.Errors:
				logger.L(logger.Red, "%s", err)
			}
		}
	}()

	logger.L(logger.Purple, "Ajout d'un watcher sur le dossier : %s", location)
	if len(location) > 0 {
		if err := watch.Add(location); err != nil {
			logger.L(logger.Red, "location: %s %s", location, err)
		}
	}

	<-done
}

func _ticker(event fsnotify.Event, c *chan bool) {
	ticker := time.NewTicker(5 * time.Second)
	var size int64 = -1
	go func() {
		for range ticker.C {
			f, err := os.Stat(event.Name)
			if err != nil {
				logger.L(logger.Red, "%s", err)
			}
			//log.Printf("Name: %s\n\tInfo size: %d - Size: %d\n\n", event.Name, f.Size(), size)
			if f.Size() != size {
				size = f.Size()
				continue
			}
			ticker.Stop()
			*c <- true
		}
	}()
}

func _stat(event fsnotify.Event) (os.FileInfo, fsnotify.Event) {
	f, err := os.Stat(event.Name)
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
	return f, event
}

func _checkIfDir(event fsnotify.Event) (isDir bool, isNil bool) {
	f, e := _stat(event)
	//log.Printf("f: %v, e: %v", f, e)
	//Ajout d'une sécurité si le fichier a déjà été déplacé
	if f == nil {
		return false, true
	}
	if f.IsDir() && filepath.Dir(f.Name()) != constants.A_TRIER {
		err := watch.Add(e.Name)
		logger.L(logger.Purple, "Ajout d'un watcher sur "+e.Name)
		if err != nil {
			logger.L(logger.Red, "%s", err)
		} else {
			return true, false
		}
	}
	return false, false
}

//var wg sync.Mutex

func _fsNotifyCreateFile(event fsnotify.Event, re *regexp.Regexp) {
	_, e := _stat(event)

	finish := make(chan bool)
	go _ticker(event, &finish)
	<-finish

	if re.MatchString(filepath.Ext(e.Name)) {
		logger.L(logger.Purple, "Détection de : %s", filepath.Base(e.Name))
		var m myFile
		m.file = event.Name
		//wg.Lock()
		m.Process()
		//wg.Unlock()
	}
}
