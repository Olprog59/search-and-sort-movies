package myapp

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/model"
)

func getAllFiles() model.AllFiles {
	return model.AllFiles{
		Movie: getMovies(),
		Serie: getSeries(),
	}
}

func getMovies() model.Movie {
	var movies model.Movie
	err := filepath.Walk(GetEnv("movies"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		var file model.File
		file.Name = info.Name()
		file.Date = info.ModTime()
		movies.Files = append(movies.Files, file)

		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return movies
}

func getSeries() model.Serie {
	var series model.Serie
	var folders model.Folder
	var seasons model.Season
	var re = regexp.MustCompile(`(?mi)season-\d+`)
	err := filepath.Walk(GetEnv("series"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.Name() == filepath.Base(GetEnv("series")) || info.Name()[0] == '.' {
			return nil
		}
		if info.IsDir() && !re.MatchString(info.Name()) {
			log.Printf("series folder: %d", len(series.Folder))
			if info.Name() != folders.Name && series.Folder != nil {
				series.Folder = append(series.Folder, folders)
				folders = model.Folder{}
				folders.Name = info.Name()
				return nil
			} else {
				folders.Name = info.Name()
				return nil
			}
		} else if info.IsDir() && re.MatchString(info.Name()) {
			log.Printf("folders seasons: %d", len(seasons.Files))
			if info.Name() != seasons.Name && seasons.Files != nil {
				folders.Seasons = append(folders.Seasons, seasons)
				seasons = model.Season{}
				seasons.Name = info.Name()
				return nil
			} else {
				seasons.Name = info.Name()
				return nil
			}
		} else {
			var file model.File
			file.Name = info.Name()
			file.Date = info.ModTime()
			seasons.Files = append(seasons.Files, file)
			return nil
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	folders.Seasons = append(folders.Seasons, seasons)
	series.Folder = append(series.Folder, folders)
	return series
}
