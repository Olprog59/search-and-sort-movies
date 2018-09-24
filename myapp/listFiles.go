package myapp

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Movie struct {
	Name  string
	Path  string
	Files []File
}

type Serie struct {
	Name   string
	Path   string
	Series []Series
}

type Series struct {
	Name    string
	Path    string
	Seasons []Season
}

type Season struct {
	Name  string
	Path  string
	Files []File
}

type File struct {
	Name string
	Path string
}

func ReadAllFiles() []string {

	root := GetEnv("dlna")
	f, err := ioutil.ReadDir(root)
	if err != nil {
		log.Println(err)
	}
	var files []string
	for _, v := range f {
		if v.IsDir() {
			continue
		}
		files = append(files, v.Name())
	}
	return files
}

func ReadFileLog() (data []string) {
	file, err := os.Open("./log_SearchAndSort")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}

func AllMovies() Movie {
	var movie Movie
	err := filepath.Walk(GetEnv("movies"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		movie.Files = append(movie.Files, File{Name: info.Name(), Path: path})
		return nil
	})
	if err != nil {
		panic(err)
	}
	return movie
}

func AllSeries() *Serie {
	var serie Serie
	path := GetEnv("series")
	for _, s := range read(path) {
		if s.IsDir() {
			var series Series
			series.Name = s.Name()
			series.Path = path + string(filepath.Separator) + series.Name
			for _, se := range read(series.Path) {
				if se.IsDir() {
					var season Season
					season.Name = se.Name()
					season.Path = series.Path + string(filepath.Separator) + season.Name
					for _, fi := range read(season.Path) {
						if !fi.IsDir() {
							var file File
							file.Name = fi.Name()
							file.Path = season.Path + string(filepath.Separator) + file.Name
							season.Files = append(season.Files, file)
						}
					}
					series.Seasons = append(series.Seasons, season)
				}
			}
			serie.Series = append(serie.Series, series)
		}
	}
	return &serie
}

func read(path string) []os.FileInfo {
	file, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	return file
}
