package myapp

import (
	"bufio"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Movie struct {
	gorm.Model
	Name          string `json:"name"`
	Path          string `json:"path"`
	Image         string `json:"image"`
	OriginalTitle string `json:"original_title"`
	Description   string `json:"description"`
}

type Serie struct {
	gorm.Model
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Image         string   `json:"image"`
	OriginalTitle string   `json:"original_title"`
	Description   string   `json:"description"`
	Seasons       []Season `json:"seasons" gorm:"foreignkey:SerieID"`
}

type Season struct {
	gorm.Model
	Name    string `json:"name"`
	Path    string `json:"path"`
	Number  int    `json:"number"`
	Files   []File `json:"files" gorm:"foreignkey:SeasonID"`
	SerieID uint
}

type File struct {
	gorm.Model
	Name          string `json:"name"`
	Path          string `json:"path"`
	Image         string `json:"image"`
	OriginalTitle string `json:"original_title"`
	Description   string `json:"description"`
	SeasonID      uint
}

//func ReadAllFiles() []string {
//
//	root := GetEnv("dlna")
//	f, err := ioutil.ReadDir(root)
//	if err != nil {
//		log.Println(err)
//	}
//	var files []string
//	for _, v := range f {
//		if v.IsDir() {
//			continue
//		}
//		files = append(files, v.Name())
//	}
//	return files
//}

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

type mapLog struct {
	Date string
	Log  []logString
}

type logString struct {
	Hour        string
	File        string
	Description string
}

func ReadFileLog() (data []mapLog) {
	reHour, _ := regexp.Compile(`(?m)\d{2}\:\d{2}\:\d{2}`)
	reDate, _ := regexp.Compile(`(?m)\d{4}\/\d{2}\/\d{2}`)
	re1, _ := regexp.Compile(`(?m)\S+\.go\:\d+\:`)
	var re2 = regexp.MustCompile(`(?m)\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2} \S+\.go\:\d+\:\s`)

	file, err := os.Open(LOGFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var date string
	var oldDate string
	var mapL *mapLog
	tabLog := []logString{}

	scanner := bufio.NewScanner(file)
	var count int

	for scanner.Scan() {
		log := logString{}
		text := scanner.Text()

		date = reDate.FindString(text)
		if count == 0 {
			oldDate = date
		}

		if date == oldDate {
			log.Hour = reHour.FindString(text)
			log.File = re1.FindString(text)
			log.Description = strings.Join(re2.Split(text, -1), "")
			if log.File != "" {
				tabLog = append(tabLog, log)
			}
		} else if count > 0 && date != "" && len(tabLog) > 0 {
			mapL = new(mapLog)
			mapL.Date = oldDate
			mapL.Log = tabLog
			data = append(data, *mapL)
			tabLog = []logString{}
		}

		if date != "" {
			oldDate = date
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}

//func ifMatchVideoFile(file string) bool {
//	re := regexp.MustCompile(`(.mkv|.mp4|.avi|.flv)`)
//	if !re.MatchString(filepath.Ext(file)) {
//		return false
//	}
//	return true
//}

func saveMovie(movieOnline movieDBMovie, name, path string) {
	var posterPath string
	var originalTitle string
	var description string
	if len(movieOnline.Results) > 0 {
		posterPath = "https://image.tmdb.org/t/p/w200" + movieOnline.Results[0].PosterPath
		originalTitle = movieOnline.Results[0].OriginalTitle
		description = movieOnline.Results[0].Overview
	}
	movie := &Movie{Name: name, Path: path, OriginalTitle: originalTitle, Image: posterPath, Description: description}
	testDb(func(db *gorm.DB) {
		res := db.Where("path = ? and name = ?", path, name).First(&movie)
		if res.RowsAffected == 0 {
			db.Create(&movie)
		}
	})
}

//
//func SaveAllMovies() bool {
//	err := filepath.Walk(GetEnv("movies"), func(path string, info os.FileInfo, err error) error {
//		if !info.IsDir() {
//			if !ifMatchVideoFile(info.Name()) {
//				return nil
//			}
//			name, _, _, year := slugFile(info.Name())
//			nameClean := name[:len(name)-len(filepath.Ext(name))]
//			movieOnline, _ := dbMovies(false, nameClean, strconv.Itoa(year))
//			var posterPath string
//			var originalTitle string
//			var description string
//			if len(movieOnline.Results) > 0 {
//				posterPath = "https://image.tmdb.org/t/p/w500" + movieOnline.Results[0].PosterPath
//				originalTitle = movieOnline.Results[0].OriginalTitle
//				description = movieOnline.Results[0].Overview
//			}
//			movie := &Movie{Name: info.Name(), Path: path, OriginalTitle: originalTitle, Image: posterPath, Description: description}
//			testDb(func(db *gorm.DB) {
//				res := db.Where("path = ? and name = ?", path, name).First(&movie)
//				if res.RowsAffected == 0 {
//					db.Create(&movie)
//				}
//			})
//		}
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	return true
//}

func splitSeriePath(path string) (serieName, serieSeason, fileName string) {
	a := path[len(GetEnv("series")):]
	if a[0] == '/' {
		a = a[1:]
	}
	tab := strings.Split(a, string(os.PathSeparator))
	serieName = tab[0]
	serieSeason = tab[1]
	fileName = tab[2]
	return serieName, serieSeason, fileName
}

func saveSerie(serieOnline movieDBTv, name, path string) {
	var serie Serie
	var season Season
	var nameSerie string
	var fileName string
	nameSerie, season.Name, fileName = splitSeriePath(path)

	serie.Name = name
	serie.Path = GetEnv("series") + string(os.PathSeparator) + nameSerie
	if len(serieOnline.Results) > 0 {
		serie.Image = "https://image.tmdb.org/t/p/w200" + serieOnline.Results[0].PosterPath
		serie.OriginalTitle = serieOnline.Results[0].OriginalName
		serie.Description = serieOnline.Results[0].Overview
	}
	season.Path = serie.Path + string(os.PathSeparator) + season.Name
	season.Number, _ = strconv.Atoi(strings.Split(season.Name, "-")[1])
	fmt.Println(season.Number)

	var file File
	file.Name = fileName
	file.Path = season.Path + string(filepath.Separator) + file.Name

	season.Files = append(season.Files, file)
	serie.Seasons = append(serie.Seasons, season)

	testDb(func(db *gorm.DB) {
		resSerie := db.Where("path = ?", serie.Path).First(&serie)
		db.Where("path = ?", season.Path).First(&season)
		resFile := db.Where("path = ?", file.Path).First(&file)
		if resSerie.RowsAffected == 0 {
			db.Save(&serie)
			return
		} else if resFile.RowsAffected == 0 {
			serie.Seasons = nil
			serie.Seasons = append(serie.Seasons, season)
			db.Save(&serie)
			return
		}
	})
}

//func SaveAllSeries() bool {
//	path := GetEnv("series")
//	for _, s := range read(path) {
//		if s.IsDir() {
//			var serie Serie
//			serie.Name = s.Name()
//			serie.Path = path + string(filepath.Separator) + serie.Name
//			//serieOnline, _ := dbMovies(false, serie.Name)
//			serieOnline, _ := dbSeries(false, serie.Name, "")
//			if len(serieOnline.Results) > 0 {
//				serie.Image = "https://image.tmdb.org/t/p/w500" + serieOnline.Results[0].PosterPath
//				serie.OriginalTitle = serieOnline.Results[0].OriginalName
//				serie.Description = serieOnline.Results[0].Overview
//			}
//			for _, se := range read(serie.Path) {
//				if se.IsDir() {
//					var season Season
//					season.Name = se.Name()
//					season.Path = serie.Path + string(filepath.Separator) + season.Name
//
//					for _, fi := range read(season.Path) {
//						if !fi.IsDir() {
//							var file File
//							file.Name = fi.Name()
//							file.Path = season.Path + string(filepath.Separator) + file.Name
//							season.Files = append(season.Files, file)
//						}
//					}
//					serie.Seasons = append(serie.Seasons, season)
//				}
//			}
//
//			testDb(func(db *gorm.DB) {
//				res := db.Where("path = ?", path).First(&serie)
//				if res.RowsAffected == 0 {
//					db.Create(&serie)
//				}
//			})
//
//		}
//	}
//	return true
//}

//func readJSON(file string) []byte {
//	f, err := ioutil.ReadFile(file)
//
//	if err != nil {
//		log.Println(err)
//	}
//
//	return f
//}

func ReadAllMovies() []Movie {
	var movie []Movie
	//json.Unmarshal(readJSON(MOVIESFILE), &movie)
	testDb(func(db *gorm.DB) {
		db.Find(&movie)
	})
	fmt.Println(movie)
	return movie
}

func ReadAllSeries() []Serie {
	var serie []Serie
	//json.Unmarshal(readJSON(SERIESFILE), &serie)
	testDb(func(db *gorm.DB) {
		db.Set("gorm:auto_preload", true).Find(&serie)
	})
	fmt.Println(serie)
	return serie
}

//func read(path string) []os.FileInfo {
//	file, err := ioutil.ReadDir(path)
//	if err != nil {
//		log.Println(err)
//	}
//	return file
//}
