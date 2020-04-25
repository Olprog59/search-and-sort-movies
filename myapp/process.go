package myapp

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	dlna   = GetEnv("dlna")
	movies = GetEnv("movies")
	series = GetEnv("series")
	count  = 0
)

func Process(complete string) {
	dir, file := filepath.Split(complete)
	go start(complete, dir, file)
}

func start(complete, dir, file string) {
	name, serieName, serieNumber, year := slugFile(file)
	// Si c'est un film
	if serieName == "" {
		extension := filepath.Ext(name)
		nameClean := name[:len(name)-len(extension)]
		originalName := file[:len(file)-len(extension)]
		originalName = url.QueryEscape(originalName)
		if count > 1 {
			originalName = ""
		}
		movie, _ := dbMovies(false, nameClean, originalName)
		if len(movie.Results) > 0 {
			var path string
			if year != 0 {
				path = movies + string(os.PathSeparator) + nameClean + "-" + strconv.Itoa(year) + extension
			} else {
				path = movies + string(os.PathSeparator) + nameClean + extension
			}
			if runtime.GOOS == "windows" {
				copyFile(complete, movies+string(os.PathSeparator)+path)
			} else {
				if moveOrRenameFile(complete, path) {
					log.Printf("%s a bien été déplacé dans %s", name, path)
				}
			}
		} else {
			if count < 3 {
				count++
				time.Sleep(2000 * time.Millisecond)
				start(complete, dir, file)
			} else {
				message := fmt.Sprintln(name[:len(name)-len(filepath.Ext(name))] + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + name[:len(name)-len(filepath.Ext(name))] + ". Test manuellement si tu le trouves ;-)")
				log.Println(message)
				count = 0
			}
		}

	} else {
		originalName := file[:len(file)-len(filepath.Ext(name))-len(serieName)]
		originalName = url.QueryEscape(originalName)
		if count > 1 {
			originalName = ""
		}
		serie, _ := dbSeries(false, serieName, originalName)
		if len(serie.Results) > 0 {
			_, season, _ := slugSerieSeasonEpisode(serieNumber)
			checkFolderSerie(complete, file, name, serieName, season)
		} else {
			if count < 3 {
				count++
				time.Sleep(2000 * time.Millisecond)
				start(complete, dir, file)
			} else {
				log.Println(serieName + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + serieName + ".\n Test manuellement si tu le trouves ;-)")
				count = 0
			}
		}
	}
}

func checkFolderSerie(complete, file, name, serieName string, season int) (string, string) {
	// serieName, exist := folderExist(series, serieName)
	newFolder := string(os.PathSeparator) + serieName + string(os.PathSeparator) + "season-" + strconv.Itoa(season)
	folderOk := series + string(os.PathSeparator) + serieName
	if _, err := os.Stat(folderOk); os.IsNotExist(err) {
		log.Printf("Création du dossier : %s\n", serieName)
		createFolder(folderOk)
	}
	if _, err := os.Stat(series + newFolder); os.IsNotExist(err) {
		log.Printf("Création du dossier : %s\n", newFolder)
		createFolder(series + newFolder)
	}

	finalFilePath := series + newFolder + string(os.PathSeparator) + name

	if runtime.GOOS == "windows" {
		copyFile(complete, finalFilePath)
	} else {
		if moveOrRenameFile(complete, finalFilePath) {
			log.Printf("%s a bien été déplacé dans %s", name, finalFilePath)
		}
	}
	return complete, finalFilePath
}

func createFolder(folder string) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func moveOrRenameFile(filePathOld, filePathNew string) bool {
	err := os.Rename(filePathOld, filePathNew)
	if err != nil {
		log.Printf("Move Or Rename File : %s", err)
		return false
	}
	folder := filepath.Dir(filePathOld)
	if folder != GetEnv("dlna") {
		file, _ := ioutil.ReadDir(folder)
		if len(file) == 0 {
			_ = watch.Remove(folder)
			log.Println("remove 1 : ", folder)
			err := os.Remove(folder)
			if err != nil {
				log.Println("error de suppression de dossier")
			}
		}
	}
	return true
}

func copyFile(oldFile, newFile string) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		r, err := os.Open(oldFile)
		if err != nil {
			log.Println(err)
		}
		defer r.Close()

		defer wg.Done()

		w, err := os.Create(newFile)
		if err != nil {
			log.Println(err)
		}
		defer w.Close()

		// do the actual work
		n, err := io.Copy(w, r)
		if err != nil {
			log.Println(err)
		}
		err = os.Chown(newFile, 0, 0)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Copied file : %s to %s - %v bytes\n", oldFile, newFile, n)
		return
	}()

	wg.Wait()
	go func() {
		err := checkIfSizeIsSame(oldFile, newFile)
		if err != nil {
			log.Println(err)
			copyFile(oldFile, newFile)
		} else {
			log.Println("Le fichier est correctement copié. La source va être supprimé !")
			removeAfterCopy(oldFile)
		}
	}()
}

func checkSizeFile(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		log.Println(err)
		return f.Size(), err
	}

	return f.Size(), err
}

func checkIfSizeIsSame(oldFile, newFile string) error {
	newFileSize, err := checkSizeFile(newFile)
	if err != nil {
		return err
	}

	oldFileSize, err := checkSizeFile(oldFile)
	if err != nil {
		return err
	}

	if newFileSize != oldFileSize {
		return errors.New("The files are not the same size. The operation will start again")
	}
	return nil
}

func removeAfterCopy(oldFile string) {
	go func() {
		time.Sleep(time.Second * 10)
		log.Println("remove 2 : ", oldFile)
		err := os.Remove(oldFile)
		if err != nil {
			log.Println(err)
		}
	}()
}
