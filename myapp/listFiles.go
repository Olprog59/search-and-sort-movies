package myapp

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/model"
)

func getVideos() model.Video {
	return model.Video{
		Movie: getMovies(),
		Serie: getSeries(),
	}
}

var videos model.Video

func getMovies() []model.File {
	videos.Movie = []model.File{}
	err := filepath.Walk(GetEnv("movies"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() || info.Name()[0] == '.' {
			return nil
		}
		var file model.File
		var id int64
		file.Name = info.Name()
		file.Image, id = model.GetImage(file.Name, false)
		file.Date = info.ModTime()
		file.Taille = info.Size()
		file.TrailerKey, file.TrailerSite = model.GetTrailer(id, false)
		videos.Movie = append(videos.Movie, file)
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return videos.Movie
}

func getSeries() []model.Serie {
	videos.Serie = []model.Serie{}
	var serie model.Serie
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
			if info.Name() != serie.Name && serie.Seasons != nil {
				if len(seasons.Files) > 0 {
					serie.Seasons = append(serie.Seasons, seasons)
					seasons = model.Season{}
				}
				videos.Serie = append(videos.Serie, serie)
				serie = model.Serie{}
				serie.Name = info.Name()
				serie = getMovieDb(serie)
				return nil
			} else if info.Name() != serie.Name && len(seasons.Files) > 0 {
				serie.Seasons = append(serie.Seasons, seasons)
				seasons = model.Season{}
				videos.Serie = append(videos.Serie, serie)
				serie = model.Serie{}
				serie.Name = info.Name()
				serie = getMovieDb(serie)
				return nil
			} else {
				serie.Name = info.Name()
				serie = getMovieDb(serie)
				return nil
			}
		} else if info.IsDir() && re.MatchString(info.Name()) {
			if info.Name() != seasons.Name && seasons.Files != nil {
				serie.Seasons = append(serie.Seasons, seasons)
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
			file.Taille = info.Size()
			seasons.Files = append(seasons.Files, file)
			return nil
		}
	})
	if err != nil {
		log.Println(err)
	}
	serie.Seasons = append(serie.Seasons, seasons)
	videos.Serie = append(videos.Serie, serie)
	return videos.Serie
}

func getMovieDb(serie model.Serie) model.Serie {
	var id int64
	if serie.Name != "" {
		serie.Image, id = model.GetImage(serie.Name, true)
		serie.TrailerKey, serie.TrailerSite = model.GetTrailer(id, true)
	}
	return serie
}
