package myapp

import (
	"encoding/json"
	"fmt"
	"github.com/Machiel/slugify"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	movies = constants.MOVIES
	series = constants.SERIES
)

type myFile struct {
	file           string
	fileWithoutDir string
	complete       string
	name           string
	SearchEngine   string
	transName      string
	serieName      string
	serieNumber    string
	season         string
	year           int
	episode        string
	count          int
}

func (m *myFile) Process() {
	m.count = 0
	_, m.complete = filepath.Split(m.file)
	m.fileWithoutDir = m.complete
	log.Println(logger.Info("complete: ", m.complete))
	m.start("")
}

func (m *myFile) start(serieOrMovieOrBoth string) {
	m.slugFile()
	if m.serieName == "" || serieOrMovieOrBoth == "serie" {
		m.isMovie()
	} else {
		m.isSerie()
	}
}

func (m *myFile) isMovie() {
	extension := filepath.Ext(m.file)
	log.Println(logger.Info("name: ", m.name))
	var movie MoviesDb
	if m.transName != "" {
		movie, _ = m.dbMovies(false, m.transName)
	} else {
		movie, _ = m.dbMovies(false, m.name)
	}
	if len(movie.Results) > 0 {
		if m.SearchEngine != "" {
			m.name = slugify.Slugify(m.SearchEngine)
		}
		var path1 string
		m.complete = m.name + extension
		if m.year != 0 {
			path1 = movies + string(os.PathSeparator) + m.name + "-" + strconv.Itoa(m.year) + extension
		} else {
			path1 = movies + string(os.PathSeparator) + m.complete
		}

		if moveOrRenameFile(m.file, path1) {
			log.Println(logger.Info(m.fileWithoutDir + ", a bien été déplacé dans : " + path1))
			m.createFileForLearning(true)
		}

	} else {
		log.Println(logger.Warn("isMovie : " + m.name))
		m.isNotFindInMovieDb(m.name, "movie")
	}
}

func (m *myFile) isSerie() {
	serie, _ := m.dbSeries(false, m.serieName)
	if len(serie.Results) > 0 {
		if m.SearchEngine != "" {
			m.serieName = slugify.Slugify(m.SearchEngine)
		}
		m.checkFolderSerie()
	} else {
		m.isNotFindInMovieDb(m.serieName, "serie")
	}
}

func (m *myFile) isNotFindInMovieDb(name, serieOrMovie string) {
	if m.count < 1 {
		m.count++
		time.Sleep(2000 * time.Millisecond)
		m.start(serieOrMovie)
	} else if m.count < 2 {
		m.count++
		time.Sleep(2000 * time.Millisecond)
		//m.translateName()
		m.start(serieOrMovie)
	} else {
		log.Println(logger.Warn(name + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + name + ".\n Test manuellement si tu le trouves ;-)"))
		m.createFileForLearning(false)
		m.count = 0
	}
}

func (m *myFile) checkFolderSerie() (string, string) {
	// serieName, exist := folderExist(series, serieName)
	newFolder := string(os.PathSeparator) + m.serieName + string(os.PathSeparator) + "season-" + m.season[1:]
	folderOk := series + string(os.PathSeparator) + m.serieName
	if _, err := os.Stat(folderOk); os.IsNotExist(err) {
		log.Printf(logger.Info("Création du dossier : " + m.serieName))
		createFolder(folderOk)
	}
	if _, err := os.Stat(series + newFolder); os.IsNotExist(err) {
		log.Println(logger.Info("Création du dossier : " + newFolder))
		createFolder(series + newFolder)
	}

	finalFilePath := series + newFolder + string(os.PathSeparator) + m.complete

	if moveOrRenameFile(m.file, finalFilePath) {
		log.Println(logger.Info(m.fileWithoutDir + ", a bien été déplacé dans : " + finalFilePath))
		m.createFileForLearning(true)
	}
	return m.complete, finalFilePath
}

// Ok test
func (m *myFile) translateName() {
	slug := slugify.New(slugify.Configuration{
		ReplaceCharacter: ' ',
	})
	var n = m.name
	n = slug.Slugify(n)
	resp, err := http.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=en&dt=t&q=" + url.PathEscape(n))
	if resp == nil || err != nil {
		time.Sleep(1 * time.Minute)
		m.translateName()
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Request)

	if resp.StatusCode != http.StatusOK {
		log.Println(logger.Warn("response status code was %d\n", resp.StatusCode))
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		log.Println(logger.Warn("response content type was " + ctype + " not text/html"))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(logger.Warn(err))
	}

	var arr [][][]string
	_ = json.Unmarshal(body, &arr)
	if err != nil {
		log.Println(logger.Warn(err))
	}

	m.transName = arr[0][0][0]
}

func createFolder(folder string) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Println(logger.Warn(err))
	}
}

var mu sync.Mutex

func moveOrRenameFile(filePathOld, filePathNew string) bool {
	mu.Lock()
	err := syscall.Rename(filePathOld, filePathNew)
	log.Println(logger.Warn(runtime.GOOS, runtime.GOARCH))
	//err := MoveFile(filePathOld, filePathNew)
	//cmd := exec.Command("/bin/sh", "-c", "mv "+filePathOld+" "+filePathNew)
	//log.Println(logger.Info("Test du mv avec exec.Command"))
	//err := cmd.Run()
	if err != nil {
		log.Println(logger.Warn("Move Or Rename File : ", err))
		mu.Unlock()
		return false
	}
	folder := filepath.Dir(filePathOld)
	if folder != constants.A_TRIER {
		file, _ := ioutil.ReadDir(folder)
		if len(file) == 0 {
			err = watch.Remove(folder)
			if err != nil {
				log.Println(logger.Warn("Erreur sur la suppression du watcher sur le dossier : ", folder))
			}
			log.Println(logger.Info("Suppression du watcher sur le dossier : ", folder))
			err := os.Remove(folder)
			if err != nil {
				log.Println(logger.Warn("Erreur de suppression de dossier : ", folder))
			}
		}
	}
	mu.Unlock()
	return true
}

// création d'un fichier pour le learning
func (m *myFile) createFileForLearning(videosTry bool) {
	f, err := os.OpenFile(path.Clean(constants.LearningFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(logger.Warn(err))
	}
	_, err = f.Write([]byte(fmt.Sprintf("%s;%s;%t\n", m.fileWithoutDir, m.complete, videosTry)))
	if err != nil {
		log.Println(logger.Warn(err))
	}
}
