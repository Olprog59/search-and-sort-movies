package controllers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	dlna   = GetEnv("dlna")
	movies = GetEnv("movies")
	series = GetEnv("series")
)

func Process(file string) {
	re := regexp.MustCompile(`(.mkv|.mp4|.avi|.flv)`)
	if !re.MatchString(filepath.Ext(file)) {
		log.Println("Le fichier " + file + " n'est pas reconnu en tant que média")
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
			moveOrRenameFile(dlna+"/"+file, movies+"/"+name)
		} else {
			log.Println(nameClean + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + nameClean + ".\n Test manuellement si tu le trouves ;-)")
		}

	} else {
		serie, _ := dbSeries(false, serieName, strconv.Itoa(year))

		if len(serie.Results) > 0 {
			season, _ := slugSerieSeasonEpisode(serieNumber)
			checkFolderSerie(file, name, serieName, season)
		} else {
			log.Println(serieName + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + serieName + ".\n Test manuellement si tu le trouves ;-)")
		}
	}
}

func folderExist(folder, serieName string) (string, bool) {
	name := searchSimilarFolder(folder, serieName)
	if name == "" {
		return serieName, false
	}
	return name, true
}

func checkFolderSerie(file, name, serieName string, season int) (string, string) {
	serieName, exist := folderExist(series, serieName)
	newFolder := string(os.PathSeparator) + serieName + string(os.PathSeparator) + "season-" + strconv.Itoa(season)
	if !exist {
		createFolder(series + newFolder)
	}
	fmt.Println(dlna+"/"+file, series+newFolder+"/"+name)
	moveOrRenameFile(dlna+"/"+file, series+newFolder+"/"+name)
	return dlna + "/" + file, series + newFolder + "/" + name
}

/*
	TODO :
	Vérifier si un dossier qui ne ressemble pas au nom
*/
func calculatePercentDiffFolder(serieName, folderExist string) float32 {
	folderExist = strings.Replace(folderExist, "-", "", -1)
	serieName = strings.Replace(serieName, "-", "", -1)

	t1 := make(map[string]int)
	t2 := make(map[string]int)

	for _, v := range serieName {
		t1[string(v)] = t1[string(v)] + 1
	}

	for _, v := range folderExist {
		t2[string(v)] = t2[string(v)] + 1
	}

	var count float32
	for k, v := range t1 {
		for l, w := range t2 {
			if k == l {
				if v == w {
					count = count + 1.0
				} else if v > w {
					count = count + (float32(w) / float32(v))
				} else if w > v {
					count = count + (float32(v) / float32(w))
				}
			}
		}
	}

	var percent float32

	if len(t1) > len(t2) {
		percent = (count / float32(len(t1))) * 100
	} else if len(t1) < len(t2) {
		percent = (count / float32(len(t2))) * 100
	} else {
		percent = (count / float32(len(t1))) * 100

	}

	return percent
}

func searchSimilarFolder(currentPath, newFolder string) string {
	var name string
	filepath.Walk(currentPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if calculatePercentDiffFolder(newFolder, f.Name()) > 65 {
				name = f.Name()
			}
		}
		return nil
	})

	return name
}
func createFolder(folder string) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func moveOrRenameFile(filePathOld, filePathNew string) {
	err := os.Rename(filePathOld, filePathNew)
	if err != nil {
		log.Println(err)
	}
}
