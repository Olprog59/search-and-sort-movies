package myapp

import (
	"encoding/json"
	"github.com/Machiel/slugify"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	movies = constants.MOVIES
	series = constants.SERIES
)

type typeSerieOrMovie uint

const (
	SERIE typeSerieOrMovie = iota
	MOVIE
	NOTHING
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
	logger.L(logger.Yellow, "complete: %s", m.complete)
	m.start(NOTHING)
}

func (m *myFile) start(serieOrMovieOrBoth typeSerieOrMovie) {
	m.slugFile()
	if m.serieName == "" || serieOrMovieOrBoth == MOVIE {
		m.isMovie()
	} else {
		m.isSerie()
	}
}

func (m *myFile) isMovie() {
	extension := filepath.Ext(m.file)
	logger.L(logger.Yellow, "name: %s", m.name)
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
			logger.L(logger.Yellow, m.fileWithoutDir+", a bien été déplacé dans : "+path1)
			m.learningFirestore(true)
		}
	} else {
		m.isNotFindInMovieDb(m.name, MOVIE)
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
		m.isNotFindInMovieDb(m.serieName, SERIE)
	}
}

func (m *myFile) isNotFindInMovieDb(name string, serieOrMovie typeSerieOrMovie) {
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
		logger.L(logger.Yellow, name+", n'a pas été trouvé sur https://www.themoviedb.org/search?query="+name+".\n Test manuellement si tu le trouves ;-)")
		m.learningFirestore(false)
		m.count = 0
	}
}

func (m *myFile) checkFolderSerie() (string, string) {
	// serieName, exist := folderExist(series, serieName)
	newFolder := string(os.PathSeparator) + m.serieName + string(os.PathSeparator) + "season-" + m.season[1:]
	folderOk := series + string(os.PathSeparator) + m.serieName
	if _, err := os.Stat(folderOk); os.IsNotExist(err) {
		logger.L(logger.Yellow, "Création du dossier : "+m.serieName)
		createFolder(folderOk)
	}
	if _, err := os.Stat(series + newFolder); os.IsNotExist(err) {
		logger.L(logger.Yellow, "Création du dossier : "+newFolder)
		createFolder(series + newFolder)
	}

	finalFilePath := series + newFolder + string(os.PathSeparator) + m.complete
	if moveOrRenameFile(m.file, finalFilePath) {
		logger.L(logger.Yellow, m.fileWithoutDir+", a bien été déplacé dans : "+finalFilePath)
		m.learningFirestore(true)
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

	logger.L(logger.Yellow, "%s", resp.Request)

	if resp.StatusCode != http.StatusOK {
		logger.L(logger.Yellow, "response status code was %d", resp.StatusCode)
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		logger.L(logger.Yellow, "response content type was "+ctype+" not text/html")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}

	var arr [][][]string
	_ = json.Unmarshal(body, &arr)
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}

	m.transName = arr[0][0][0]
}

func createFolder(folder string) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
}

var mu sync.Mutex

func moveOrRenameFile(filePathOld, filePathNew string) bool {
	mu.Lock()
	//err := syscall.Rename(filePathOld, strings.ToLower(filePathNew))
	//logger.L(logger.Yellow, "Goos: %s - GoArch: %s", runtime.GOOS, runtime.GOARCH)
	//err := MoveFile(filePathOld, filePathNew)
	cmd := exec.Command("/bin/sh", "-c", "mv "+filePathOld+" "+filePathNew)
	logger.L(logger.Yellow, "mv "+filePathOld+" "+filePathNew)
	err := cmd.Run()
	if err != nil {
		logger.L(logger.Red, "Move Or Rename File : %s", err)
		mu.Unlock()
		return false
	}
	folder := filepath.Dir(filePathOld)

	folder = getAbsolutePathWithRelative(folder)
	absoluteATrier := getAbsolutePathWithRelative(constants.A_TRIER)

	if folder != absoluteATrier {
		file, _ := ioutil.ReadDir(folder)
		if len(file) == 0 {
			err = watch.Remove(folder)
			if err != nil {
				logger.L(logger.Red, "Erreur sur la suppression du watcher sur le dossier : %s", folder)
			}
			logger.L(logger.Yellow, "Suppression du watcher sur le dossier : %s", folder)
			err := os.Remove(folder)
			if err != nil {
				logger.L(logger.Red, "Erreur de suppression de dossier : %s", folder)
			}
		}
	}
	mu.Unlock()
	return true
}

func getAbsolutePathWithRelative(folder string) string {
	abs, err := filepath.Abs(folder)
	if err == nil {
		return abs
	}
	return ""
}
