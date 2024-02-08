package lib

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/sam-docker/media-organizer/constants"
	"github.com/sam-docker/media-organizer/logger"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

var watch *fsnotify.Watcher
var err error

//func MyWatcher(location string) {
//	watch, err = fsnotify.NewWatcher()
//	if err != nil {
//		logger.L(logger.Red, "%s", err)
//	}
//	defer watch.Close()
//
//	done := make(chan bool)
//
//	go func() {
//		for {
//			select {
//			case event, ok := <-watch.Events:
//				if !ok {
//					return
//				}
//				if event.Op&fsnotify.Write == fsnotify.Write {
//					go func(e fsnotify.Event) {
//						isDir, isNil := _checkIfDir(e)
//						if isNil {
//							return
//						}
//
//						if !isDir {
//							re := regexp.MustCompile(constants.RegexFileExtension)
//							if re.MatchString(filepath.Ext(e.Name)) {
//								_fsNotifyCreateFile(e.Name, re)
//							}
//						} else if len(e.Name) > 0 {
//							if err := watch.Add(e.Name); err != nil {
//								logger.L(logger.Red, "Error adding watcher: %s %s", e.Name, err)
//							}
//						}
//					}(event)
//				}
//			case err, ok := <-watch.Errors:
//				if !ok {
//					return
//				}
//				logger.L(logger.Red, "Watcher error: %s", err)
//			}
//		}
//	}()
//
//	logger.L(logger.Purple, "Add watcher to folder: %s", location)
//	if len(location) > 0 {
//		if err := watch.Add(location); err != nil {
//			logger.L(logger.Red, "location: %s %s", location, err)
//		}
//	}
//
//	<-done
//	logger.L(logger.Purple, "Watcher finished")
//}

func MyWatcher(location string, obsSlice *ObservableSlice) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
	defer watch.Close()

	done := make(chan bool)

	go func() {
	first:
		for {
			select {
			case event, ok := <-watch.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					if isWriteComplete(event.Name) {
						duration, err := getMediaDuration(event.Name)
						if err != nil {
							logger.L(logger.Red, "Error checking media file: %s", err)
							continue
						}
						if obsSlice.SameItem(event.Name) {
							continue first
						}

						isDir, isNil := _checkIfDir(event)
						if isNil {
							continue
						}

						if !isDir {
							re := regexp.MustCompile(constants.RegexFileExtension)
							if re.MatchString(filepath.Ext(event.Name)) {
								file := sliceFile{file: event.Name, duration: duration, working: false}
								obsSlice.Add(file)
								//logger.L(logger.Purple, "File write complete: %s duration: %s", event.Name, duration)
							}
						} else if len(event.Name) > 0 {
							if err := watch.Add(event.Name); err != nil {
								logger.L(logger.Red, "Error adding watcher: %s %s", event.Name, err)
							}
						}
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					obsSlice.Remove(event.Name)
					//logger.L(logger.Purple, "File removed or rename: %s", event.Name)
				}
			case err, ok := <-watch.Errors:
				if !ok {
					return
				}
				logger.L(logger.Red, "Watcher error: %s", err)
			}
		}
	}()

	log.Println("Add watcher to folder:", location)
	if len(location) > 0 {
		if err := watch.Add(location); err != nil {
			logger.L(logger.Red, "location: %s %s", location, err)
		}
	}

	<-done
	logger.L(logger.Purple, "Watcher finished")
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
	if f.IsDir() && filepath.Dir(f.Name()) != constants.BE_SORTED {
		err := watch.Add(e.Name)
		logger.L(logger.Purple, "Add watcher : "+e.Name)
		if err != nil {
			logger.L(logger.Red, "%s", err)
		} else {
			return true, true
		}
	}
	return false, false
}

func isWriteComplete(filePath string) bool {
	var lastSize int64 = -1
	const checks = 3 // Nombre de vérifications à effectuer
	for i := 0; i < checks; i++ {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			logger.L(logger.Red, "%s", err)
			return false
		}
		currentSize := fileInfo.Size()
		if lastSize == currentSize && i > 0 { // Si la taille est constante et que ce n'est pas la première vérification
			return true
		}
		lastSize = currentSize
		time.Sleep(2 * time.Second) // Attendre 2 secondes avant la prochaine vérification
	}
	return false
}

// getMediaDuration utilise ffprobe pour récupérer la durée d'un fichier média et la retourne sous forme de chaîne de caractères.
func getMediaDuration(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	duration := strings.TrimSpace(out.String()) // Nettoie la sortie pour obtenir uniquement la durée
	return duration, nil
}

func StartProcessing(slice []sliceFile, obsSlice *ObservableSlice) {
	for k, v := range slice {
		if !v.working {
			obsSlice.lock.Lock()
			obsSlice.slice[k].working = true
			obsSlice.lock.Unlock()
			go ProcessFile(v, obsSlice)
		}
	}
}

func ProcessFile(k sliceFile, obsSlice *ObservableSlice) {
	time.Sleep(10 * time.Second)
	logger.L(logger.Purple, "Processing file: %s", k)
	var m myFile
	m.file = k.file
	m.duration = k.duration
	m.Process()
	obsSlice.lock.Lock()
	if len(obsSlice.slice) > 0 {
		if i := slices.Index(obsSlice.slice, k); i != -1 {
			slices.Delete(obsSlice.slice, i, i+1)
		}
	}
	//obsSlice.slice = append(obsSlice.slice[:0], obsSlice.slice[1:]...)
	obsSlice.lock.Unlock()
	logger.L(logger.Purple, "File processed: %s", k)
}
