package myapp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func Process(file string) {
	re := regexp.MustCompile(`(.mkv|.mp4|.avi|.flv)`)
	if !re.MatchString(filepath.Ext(file)) {
		// log.Println("Le fichier " + file + " n'est pas reconnu en tant que média")
		return
	}

	_, file = filepath.Split(file)
	go start(file)
}

func start(file string) {
	name, serieName, serieNumber, year := slugFile(file)

	// Si c'est un film
	if serieName == "" {
		extension := filepath.Ext(name)
		nameClean := name[:len(name)-len(extension)]
		movie, _ := dbMovies(false, nameClean, strconv.Itoa(year))
		if len(movie.Results) > 0 {
			path := movies+string(os.PathSeparator)+nameClean + "-"+ strconv.Itoa(year) +"-"+extension
			if runtime.GOOS == "windows" {
				copyFile(dlna+string(os.PathSeparator)+file, movies+string(os.PathSeparator)+path)
			} else {
				moveOrRenameFile(dlna+string(os.PathSeparator)+file, path)
				log.Printf("%s a bien été déplacé dans %s", name, path)
			}
		} else {
			if count < 3 {
				count++
				time.Sleep(2000 * time.Millisecond)
				start(file)
			} else {
				message := fmt.Sprintln(name[:len(name)-len(filepath.Ext(name))] + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + name[:len(name)-len(filepath.Ext(name))] + ". Test manuellement si tu le trouves ;-)")
				//EnvoiDeMail("Search and sort movies Problem", message)
				log.Println(message)
				count = 0
			}
		}

	} else {
		serie, _ := dbSeries(false, serieName, strconv.Itoa(year))

		if len(serie.Results) > 0 {
			_, season, _ := slugSerieSeasonEpisode(serieNumber)
			checkFolderSerie(file, name, serieName, season)
		} else {
			if count < 3 {
				count++
				time.Sleep(2000 * time.Millisecond)
				start(file)
			} else {
				log.Println(serieName + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + serieName + ".\n Test manuellement si tu le trouves ;-)")
				count = 0
			}
		}
	}
	SaveAllMovies()
	SaveAllSeries()
}

func checkFolderSerie(file, name, serieName string, season int) (string, string) {
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
	oldFilePath := dlna + string(os.PathSeparator) + file

	if runtime.GOOS == "windows" {
		copyFile(oldFilePath, finalFilePath)
	} else {
		if moveOrRenameFile(oldFilePath, finalFilePath) {
			log.Printf("%s a bien été déplacé dans %s", name, finalFilePath)
		}
	}
	return oldFilePath, finalFilePath
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

		err := os.Remove(oldFile)
		if err != nil {
			log.Println(err)
		}
	}()
}
