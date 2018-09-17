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
		nameClean := name[:len(name)-len(filepath.Ext(name))]
		movie, _ := dbMovies(false, nameClean, strconv.Itoa(year))
		if len(movie.Results) > 0 {
			if runtime.GOOS == "windows" {
				copyFile(dlna+string(os.PathSeparator)+file, movies+string(os.PathSeparator)+name)
			} else {
				moveOrRenameFile(dlna+string(os.PathSeparator)+file, movies+string(os.PathSeparator)+name)
				log.Printf("%s a bien été déplacé dans %s", name, movies+string(os.PathSeparator)+name)
			}
		} else {
			movie, _ := dbMovies(false, nameClean)
			if len(movie.Results) > 0 {
				if runtime.GOOS == "windows" {
					copyFile(dlna+string(os.PathSeparator)+file, movies+string(os.PathSeparator)+name)
				} else {
					if moveOrRenameFile(dlna+string(os.PathSeparator)+file, movies+string(os.PathSeparator)+name) {
						log.Printf("%s a bien été déplacé dans %s", name, movies+string(os.PathSeparator)+name)
					}
				}
			}
			if count < 3 {
				count++
				time.Sleep(2000 * time.Millisecond)
				start(file)
			} else {
				message := fmt.Sprintln(name[:len(name)-len(filepath.Ext(name))] + ", n'a pas été trouvé sur <a href='https://www.themoviedb.org/search?query=" + name[:len(name)-len(filepath.Ext(name))] + "'>cliques ici pour vérifier sur moviedb</a>.\n Test manuellement si tu le trouves ;-)")
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
}

// func folderExist(folder, serieName string) (string, bool) {
// 	name := searchSimilarFolder(folder, serieName)
// 	if name == "" {
// 		return serieName, false
// 	}
// 	return name, true
// }

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

/*
	TODO :
	Vérifier si un dossier qui ne ressemble pas au nom
*/
// func calculatePercentDiffFolder(serieName, folderExist string) float32 {
// 	folderExist = strings.Replace(folderExist, "-", "", -1)
// 	serieName = strings.Replace(serieName, "-", "", -1)

// 	t1 := make(map[string]int)
// 	t2 := make(map[string]int)

// 	for _, v := range serieName {
// 		t1[string(v)] = t1[string(v)] + 1
// 	}

// 	for _, v := range folderExist {
// 		t2[string(v)] = t2[string(v)] + 1
// 	}

// 	var count float32
// 	for k, v := range t1 {
// 		for l, w := range t2 {
// 			if k == l {
// 				if v == w {
// 					count = count + 1.0
// 				} else if v > w {
// 					count = count + (float32(w) / float32(v))
// 				} else if w > v {
// 					count = count + (float32(v) / float32(w))
// 				}
// 			}
// 		}
// 	}

// 	var percent float32

// 	if len(t1) > len(t2) {
// 		percent = (count / float32(len(t1))) * 100
// 	} else if len(t1) < len(t2) {
// 		percent = (count / float32(len(t2))) * 100
// 	} else {
// 		percent = (count / float32(len(t1))) * 100

// 	}
// 	return percent
// }

// func searchSimilarFolder(currentPath, newFolder string) string {
// 	var name string
// 	filepath.Walk(currentPath, func(path string, f os.FileInfo, err error) error {
// 		if f.IsDir() {
// 			if calculatePercentDiffFolder(newFolder, f.Name()) > 80 {
// 				name = f.Name()
// 			}
// 		}
// 		return nil
// 	})

// 	return name
// }
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

		return
		log.Printf("Copied file : %s to %s - %v bytes\n", oldFile, newFile, n)
		log.Printf("Copied file : %s to %s - %v bytes\n", oldFile, newFile, n)
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
