package controllers

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

	// moveOrRenameFile(dlna+"/"+file, dlna+"/"+name)

	//TODO : Check name to tvdb

	// Si c'est un film
	if serieName == "" {
		nameClean := name[:len(name)-len(filepath.Ext(name))]
		movie, _ := dbMovies(false, nameClean, strconv.Itoa(year))
		if len(movie.Results) > 0 {
			moveOrRenameFile(dlna+"/"+file, movies+"/"+name)
		} else {
			log.Println(nameClean + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + nameClean + ".\n Test manuellement si tu le trouves ;-)")
		}
		/* TODO : Code qui fonctionne. Juste un petit soucis avec certain caractère lors de la comparaison
		ex : Young & Hungry (sur movieDB) et Young and Hungry (sur le serveur)

		for _, v := range movie.Results {
			v.Title, _, _, _ = slugFile(v.Title)
			v.Title = strings.Replace(v.Title, "-", " ", -1)
			nameClean = strings.Replace(nameClean, "-", " ", -1)
			log.Println(v.Title, nameClean)
			if strings.ToLower(v.Title) == nameClean {
				moveOrRenameFile(dlna+"/"+file, movies+"/"+name)
				return
			}
		} */

	} else {
		serie, _ := dbSeries(false, serieName, strconv.Itoa(year))

		if len(serie.Results) > 0 {
			season, _ := slugSerieSeasonEpisode(serieNumber)
			checkFolderSerie(file, name, serieName, season)
		} else {
			log.Println(serieName + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + serieName + ".\n Test manuellement si tu le trouves ;-)")
		}
		/* TODO : Code qui fonctionne. Juste un petit soucis avec certain caractère lors de la comparaison
		ex : Young & Hungry (sur movieDB) et Young and Hungry (sur le serveur)

		for _, v := range movie.Results {
			log.Println(v.Name)
			if strings.ToLower(v.Name) == strings.Replace(serieName, "-", " ", -1) {
				season, _ := slugSerieSeasonEpisode(serieNumber)
				checkFolderSerie(name, serieName, season)
				return
			}
		}*/
	}
}

func folderExist(folder string) (bool, error) {
	_, err := os.Stat(folder)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, err
}

func checkFolderSerie(file, name, serieName string, season int) (string, string) {
	ok, err := folderExist(series + "/" + serieName)
	newFolder := "/" + serieName + "/season-" + strconv.Itoa(season)
	if ok && err == nil {
		ok, err := folderExist(dlna + newFolder)
		if !ok || err != nil {
			createFolder(series + newFolder)
		}

	} else {
		createFolder(series + newFolder)
	}

	moveOrRenameFile(dlna+"/"+file, series+newFolder+"/"+name)
	return dlna + "/" + file, series + newFolder + "/" + name
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
