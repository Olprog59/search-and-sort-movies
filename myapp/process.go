package myapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Machiel/slugify"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"search-and-sort-movies/myapp/constants"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	movies = GetEnv("movies")
	series = GetEnv("series")
)

type myFile struct {
	file           string
	fileWithoutDir string
	complete       string
	name           string
	bingName       string
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
	log.Println("complete: ", m.complete)
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
	log.Println("name: ", m.name)
	var movie MoviesDb
	if m.transName != "" {
		movie, _ = m.dbMovies(false, m.transName)
	} else {
		movie, _ = m.dbMovies(false, m.name)
	}
	if len(movie.Results) > 0 {
		if m.bingName != "" {
			m.name = slugify.Slugify(m.bingName)
		}
		var path1 string
		m.complete = m.name + extension
		if m.year != 0 {
			path1 = movies + string(os.PathSeparator) + m.name + "-" + strconv.Itoa(m.year) + extension
		} else {
			path1 = movies + string(os.PathSeparator) + m.complete
		}
		if runtime.GOOS == "windows" {
			copyFile(m.file, movies+string(os.PathSeparator)+path1)
			m.createFileForLearning(true)
		} else {
			if moveOrRenameFile(m.file, path1) {
				log.Printf("%s a bien été déplacé dans %s", m.fileWithoutDir, path1)
				m.createFileForLearning(true)
			}
		}
	} else {
		log.Printf("isMovie : %s", m.name)
		m.isNotFindInMovieDb(m.name, "movie")
	}
}

func (m *myFile) isSerie() {
	serie, _ := m.dbSeries(false, m.serieName)
	if len(serie.Results) > 0 {
		if m.bingName != "" {
			m.serieName = slugify.Slugify(m.bingName)
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
		log.Println(name + ", n'a pas été trouvé sur https://www.themoviedb.org/search?query=" + name + ".\n Test manuellement si tu le trouves ;-)")
		m.createFileForLearning(false)
		m.count = 0
	}
}

func (m *myFile) checkFolderSerie() (string, string) {
	// serieName, exist := folderExist(series, serieName)
	newFolder := string(os.PathSeparator) + m.serieName + string(os.PathSeparator) + "season-" + m.season
	folderOk := series + string(os.PathSeparator) + m.serieName
	if _, err := os.Stat(folderOk); os.IsNotExist(err) {
		log.Printf("Création du dossier : %s\n", m.serieName)
		createFolder(folderOk)
	}
	if _, err := os.Stat(series + newFolder); os.IsNotExist(err) {
		log.Printf("Création du dossier : %s\n", newFolder)
		createFolder(series + newFolder)
	}

	finalFilePath := series + newFolder + string(os.PathSeparator) + m.complete

	if runtime.GOOS == "windows" {
		copyFile(m.file, finalFilePath)
		m.createFileForLearning(true)
	} else {
		if moveOrRenameFile(m.file, finalFilePath) {
			log.Printf("%s a bien été déplacé dans %s", m.fileWithoutDir, finalFilePath)
			m.createFileForLearning(true)
		}
	}
	return m.complete, finalFilePath
}

func (m *myFile) translateName() {
	slugify.New(slugify.Configuration{
		ReplaceCharacter: ' ',
	})
	var n = m.name
	n = slugify.Slugify(n)
	resp, _ := http.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=en&dt=t&q=" + n)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code was %d\n", resp.StatusCode)
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var arr [][][]string
	_ = json.Unmarshal(body, &arr)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	m.transName = arr[0][0][0]
}

func createFolder(folder string) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func moveOrRenameFile(filePathOld, filePathNew string) bool {
	err := os.Rename(filePathOld, filePathNew)
	if err != nil {
		log.Printf("Move Or Rename File : %s", err)
		return false
	}
	folder := filepath.Dir(filePathOld)
	if folder != GetEnv("dlna") {
		file, _ := ioutil.ReadDir(folder)
		if len(file) == 0 {
			_ = watch.Remove(folder)
			log.Println("remove 1 : ", folder)
			err := os.Remove(folder)
			if err != nil {
				log.Println("error de suppression de dossier")
			}
		}
	}
	return true
}

// si l'OS est windows alors je fais une copie et pas un déplacement
func copyFile(oldFile, newFile string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		r, err := os.Open(oldFile)
		if err != nil {
			log.Println(err)
		}
		defer r.Close()

		defer wg.Done()

		w, err := os.Create(newFile)
		if err != nil {
			log.Println(err)
		}
		defer w.Close()

		// do the actual work
		n, err := io.Copy(w, r)
		if err != nil {
			log.Println(err)
		}
		err = os.Chown(newFile, 0, 0)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Copied file : %s to %s - %v bytes\n", oldFile, newFile, n)
		return
	}()

	wg.Wait()
	go func() {
		err := checkIfSizeIsSame(oldFile, newFile)
		if err != nil {
			log.Println(err)
			copyFile(oldFile, newFile)
		} else {
			log.Println("Le fichier est correctement copié. La source va être supprimé !")
			removeAfterCopy(oldFile)
		}
	}()
}

func checkSizeFile(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		log.Println(err)
		return f.Size(), err
	}

	return f.Size(), err
}

func checkIfSizeIsSame(oldFile, newFile string) error {
	newFileSize, err := checkSizeFile(newFile)
	if err != nil {
		return err
	}

	oldFileSize, err := checkSizeFile(oldFile)
	if err != nil {
		return err
	}

	if newFileSize != oldFileSize {
		return errors.New("The files are not the same size. The operation will start again")
	}
	return nil
}

func removeAfterCopy(oldFile string) {
	go func() {
		time.Sleep(time.Second * 10)
		log.Println("remove 2 : ", oldFile)
		err := os.Remove(oldFile)
		if err != nil {
			log.Println(err)
		}
	}()
}

// création d'un fichier pour le learning
func (m *myFile) createFileForLearning(videosTry bool) {
	f, err := os.OpenFile(path.Clean(constants.LearningFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	_, err = f.Write([]byte(fmt.Sprintf("%s, %s, %t\n", m.fileWithoutDir, m.complete, videosTry)))
	if err != nil {
		log.Println(err)
	}
}
