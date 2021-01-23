package myapp

import (
	"encoding/json"
	"fmt"
	"github.com/Machiel/slugify"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
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
	fmt.Println(logger.Info("complete: ", m.complete))
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
	fmt.Println(logger.Info("name: ", m.name))
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
			fmt.Println(logger.Info("%s a bien été déplacé dans %s", m.fileWithoutDir, path1))
			m.createFileForLearning(true)
		}

	} else {
		fmt.Println(logger.Warn("isMovie : %s", m.name))
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
		m.translateName()
		m.start(serieOrMovie)
	} else {
		fmt.Println(logger.Warn(name + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + name + ".\n Test manuellement si tu le trouves ;-)"))
		m.createFileForLearning(false)
		m.count = 0
	}
}

func (m *myFile) checkFolderSerie() (string, string) {
	// serieName, exist := folderExist(series, serieName)
	newFolder := string(os.PathSeparator) + m.serieName + string(os.PathSeparator) + "season-" + m.season[1:]
	folderOk := series + string(os.PathSeparator) + m.serieName
	if _, err := os.Stat(folderOk); os.IsNotExist(err) {
		fmt.Printf(logger.Info("Création du dossier : %s", m.serieName))
		createFolder(folderOk)
	}
	if _, err := os.Stat(series + newFolder); os.IsNotExist(err) {
		fmt.Println(logger.Info("Création du dossier : %s", newFolder))
		createFolder(series + newFolder)
	}

	finalFilePath := series + newFolder + string(os.PathSeparator) + m.complete

	if moveOrRenameFile(m.file, finalFilePath) {
		fmt.Println(logger.Info("%s a bien été déplacé dans %s", m.fileWithoutDir, finalFilePath))
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
		fmt.Println(logger.Warn("response status code was %d\n", resp.StatusCode))
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		fmt.Println(logger.Warn("response content type was %s not text/html\n", ctype))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(logger.Warn(err))
	}

	var arr [][][]string
	_ = json.Unmarshal(body, &arr)
	if err != nil {
		fmt.Println(logger.Warn(err))
	}

	m.transName = arr[0][0][0]
}

func createFolder(folder string) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		fmt.Println(logger.Warn(err))
	}
}

var mu sync.Mutex

func moveOrRenameFile(filePathOld, filePathNew string) bool {
	mu.Lock()
	//err := os.Rename(filePathOld, filePathNew)
	//err := MoveFile(filePathOld, filePathNew)
	cmd := exec.Command("/bin/sh", "-c", "mv "+filePathOld+" "+filePathNew)
	fmt.Println(logger.Info("Test du mv avec exec.Command"))
	err := cmd.Run()
	if err != nil {
		fmt.Println(logger.Warn("Move Or Rename File : %s", err))
		mu.Unlock()
		return false
	}
	folder := filepath.Dir(filePathOld)
	if folder != constants.A_TRIER {
		file, _ := ioutil.ReadDir(folder)
		if len(file) == 0 {
			err = watch.Remove(folder)
			if err != nil {
				fmt.Println(logger.Warn("Erreur sur la suppression du watcher sur le dossier : ", folder))
			}
			fmt.Println(logger.Info("Suppression du watcher sur le dossier : ", folder))
			err := os.Remove(folder)
			if err != nil {
				fmt.Println(logger.Warn("Erreur de suppression de dossier : ", folder))
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
		fmt.Println(logger.Warn(err))
	}
	_, err = f.Write([]byte(fmt.Sprintf("%s;%s;%t\n", m.fileWithoutDir, m.complete, videosTry)))
	if err != nil {
		fmt.Println(logger.Warn(err))
	}
}
