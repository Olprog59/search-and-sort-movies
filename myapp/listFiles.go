package myapp

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type Movie struct {
	Files []File `json:"files"`
}

type Serie struct {
	Series []Series `json:"series"`
}

type Series struct {
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Image         string   `json:"image"`
	OriginalTitle string   `json:"original_title"`
	Description   string   `json:"description"`
	Seasons       []Season `json:"seasons"`
}

type Season struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Files []File `json:"files"`
}

type File struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Image         string `json:"image"`
	OriginalTitle string `json:"original_title"`
	Description   string `json:"description"`
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
	file, err := os.Open(LOGFILE)
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

func ifMatchVideoFile(file string) bool {
	re := regexp.MustCompile(`(.mkv|.mp4|.avi|.flv)`)
	if !re.MatchString(filepath.Ext(file)) {
		return false
	}
	return true
}

func SaveAllMovies() bool {
	var movie Movie
	err := filepath.Walk(GetEnv("movies"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if !ifMatchVideoFile(info.Name()) {
				return nil
			}
			name, _, _, year := slugFile(info.Name())
			nameClean := name[:len(name)-len(filepath.Ext(name))]
			movieOnline, _ := dbMovies(false, nameClean, strconv.Itoa(year))
			var posterPath string
			var originalTitle string
			var description string
			if len(movieOnline.Results) > 0 {
				posterPath = "https://image.tmdb.org/t/p/w500" + movieOnline.Results[0].PosterPath
				originalTitle = movieOnline.Results[0].OriginalTitle
				description = movieOnline.Results[0].Overview
			}
			movie.Files = append(movie.Files, File{Name: info.Name(), Path: path, OriginalTitle: originalTitle, Image: posterPath, Description: description})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	j, err := json.MarshalIndent(&movie, "", " ")
	if err != nil {
		log.Println(err)
	}

	writeJSONFile(MOVIESFILE, j)
	return true
}

func SaveAllSeries() bool {
	var serie Serie
	path := GetEnv("series")
	for _, s := range read(path) {
		if s.IsDir() {
			var series Series
			series.Name = s.Name()
			series.Path = path + string(filepath.Separator) + series.Name
			//serieOnline, _ := dbMovies(false, series.Name)
			serieOnline, _ := dbSeries(false, series.Name, "")
			if len(serieOnline.Results) > 0 {
				series.Image = "https://image.tmdb.org/t/p/w500" + serieOnline.Results[0].PosterPath
				series.OriginalTitle = serieOnline.Results[0].OriginalName
				series.Description = serieOnline.Results[0].Overview
			}
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
	j, err := json.MarshalIndent(&serie, "", " ")
	if err != nil {
		log.Println(err)
	}

	writeJSONFile(SERIESFILE, j)
	return true
}

func readJSON(file string) []byte {
	f, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println(err)
	}

	return f
}

func ReadAllMovies() Movie {
	var movie Movie
	json.Unmarshal(readJSON(MOVIESFILE), &movie)
	return movie
}

func ReadAllSeries() Serie {
	var serie Serie
	json.Unmarshal(readJSON(SERIESFILE), &serie)
	return serie
}

func read(path string) []os.FileInfo {
	file, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	return file
}
