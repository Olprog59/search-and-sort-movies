package lib

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/sam-docker/media-organizer/constants"
	"github.com/sam-docker/media-organizer/logger"
	"github.com/sam-docker/media-organizer/model"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

var watch *fsnotify.Watcher
var err error

func MyWatcher(location string, obsSlice *model.ObservableSlice) {
	watch, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Err("%s", err)
	}
	defer func(watch *fsnotify.Watcher) {
		err := watch.Close()
		if err != nil {
			logger.Err("Error closing watcher: %s", err)
		}
	}(watch)

	sem := make(chan struct{}, 4) // Limite à 4 goroutines simultanées

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watch.Events:
				if !ok {
					return
				}
				sem <- struct{}{}
				go func(e fsnotify.Event) {
					defer func() { <-sem }()
					//logger.Debug("Event: %s", e)
					handleEvent(e, obsSlice)
				}(event)

			case err, ok := <-watch.Errors:
				if !ok {
					return
				}
				logger.Err("Watcher error: %s", err)
			}
		}
	}()

	logger.Debug("Add watcher to folder: %s", location)
	if len(location) > 0 && strings.Contains(location, constants.BE_SORTED) {
		if err := watch.Add(location); err != nil {
			logger.Err("location: %s %s", location, err)
		}
	}

	<-done
	logger.Debug("Watcher finished")
}

func handleEvent(e fsnotify.Event, obsSlice *model.ObservableSlice) {
	if e.Op&fsnotify.Write == fsnotify.Write || e.Op&fsnotify.Create == fsnotify.Create || e.Op&fsnotify.Rename == fsnotify.Rename {
		if obsSlice.SameItem(e.Name) {
			//logger.Debug("File already in slice: %s", e.Name)
			return
		}
		if isWriteComplete(e.Name) {
			isDir, isNil := _checkIfDir(e)
			if isNil {
				//logger.Debug("File is nil: %s", e.Name)
				return
			}

			if !isDir {
				re := regexp.MustCompile(constants.RegexFileExtension)
				if re.MatchString(filepath.Ext(e.Name)) {
					duration, err := GetMediaDuration(e.Name)
					if err != nil {
						logger.Err("Erreur lors de la vérification du fichier multimédia : %s", err)
						return
					}
					file := model.SliceFile{File: e.Name, Duration: duration, Working: false}
					if !obsSlice.SameItem(e.Name) {
						obsSlice.Add(file)
						logger.Debug("Le fichier a terminé de s'écrire : %s duration: %s", e.Name, duration)
					}
				}
			}
		}
	} else if e.Op&fsnotify.Remove == fsnotify.Remove || e.Op&fsnotify.Rename == fsnotify.Rename {
		obsSlice.Remove(e.Name)
		logger.Debug("File removed or rename: %s", e.Name)
	}
}

func _stat(event fsnotify.Event) (os.FileInfo, fsnotify.Event) {
	f, err := os.Stat(event.Name)
	if err != nil {
		logger.Err("%s", err)
	}
	return f, event
}

func _checkIfDir(event fsnotify.Event) (isDir bool, isNil bool) {
	f, e := _stat(event)
	//log.Printf("f: %#+v, e: %#+v", f, e)
	//Ajout d'une sécurité si le fichier a déjà été déplacé
	if f == nil {
		return false, true
	}
	if f.IsDir() && filepath.Dir(f.Name()) != constants.BE_SORTED {
		err := watch.Add(e.Name)
		logger.Debug("Add watcher : " + e.Name)
		if err != nil {
			logger.Err("%s", err)
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

// GetMediaDuration utilise ffprobe pour récupérer la durée d'un fichier média et la retourne sous forme de chaîne de caractères.
func GetMediaDuration(filePath string) (string, error) {
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

func StartProcessing(slice []model.SliceFile, obsSlice *model.ObservableSlice) {
	for k, v := range slice {
		if !v.Working {
			obsSlice.Lock.Lock()
			obsSlice.Slice[k].Working = true
			obsSlice.Lock.Unlock()
			go ProcessFile(v, obsSlice)
		}
	}
}

func ProcessFile(k model.SliceFile, obsSlice *model.ObservableSlice) {
	time.Sleep(10 * time.Second)
	//logger.Debug("Processing file: %s", k)
	var m myFile
	m.file = k.File
	duration, err := strconv.ParseFloat(k.Duration, 10)
	if err != nil {
		logger.Err("Error parsing duration: %s", err)
		duration = 0
	}
	m.duration = duration
	if k.Force {
		m.ForceType = k.TypeMedia
	}
	m.Process()
	obsSlice.Lock.Lock()
	if len(obsSlice.Slice) > 0 {
		if i := slices.Index(obsSlice.Slice, k); i != -1 {
			slices.Delete(obsSlice.Slice, i, i+1)
		}
	}
	obsSlice.Lock.Unlock()
	//logger.Debug("File processed: %s", k)
}
