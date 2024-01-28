package myapp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	movies = constants.MOVIES
	series = constants.SERIES
)

type typeSerieOrMovie uint

const (
	MOVIE typeSerieOrMovie = iota
	NOTHING
)

type myFile struct {
	file           string
	ext            string
	resolution     string
	fileWithoutDir string
	complete       string
	completeSlug   string
	name           string
	serieName      string
	serieNumber    string
	season         int
	year           int
	episode        int
	episodeRaw     string
	count          int
	language       string
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

	var path1 string
	m.complete = m.name + extension
	if m.year != 0 {
		path1 = movies + string(os.PathSeparator) + m.name + "-" + strconv.Itoa(m.year) + extension
	} else {
		path1 = movies + string(os.PathSeparator) + m.complete
	}
	if moveOrRenameFile(m.file, path1) {
		logger.L(logger.Yellow, m.fileWithoutDir+", has been moved to: "+path1)
	}
}

func (m *myFile) isSerie() {
	m.checkFolderSerie()
}

func (m *myFile) checkFolderSerie() (string, string) {
	// serieName, exist := folderExist(series, serieName)
	ss := func() string {
		if m.season == 0 {
			return "00"
		}
		return strconv.Itoa(m.season)
	}()

	newFolder := string(os.PathSeparator) + m.serieName + string(os.PathSeparator) + "season-" + ss
	folderOk := series + string(os.PathSeparator) + m.serieName
	if _, err := os.Stat(folderOk); os.IsNotExist(err) {
		logger.L(logger.Yellow, "Create folder: "+m.serieName)
		createFolder(folderOk)
	}
	if _, err := os.Stat(series + newFolder); os.IsNotExist(err) {
		logger.L(logger.Yellow, "Create folder : "+newFolder)
		createFolder(series + newFolder)
	}

	finalFilePath := series + newFolder + string(os.PathSeparator) + m.complete
	if moveOrRenameFile(m.file, finalFilePath) {
		logger.L(logger.Yellow, m.fileWithoutDir+", has been moved to: "+finalFilePath)
	}
	return m.complete, finalFilePath
}

func (m *myFile) formatageSerie() {
	format := constants.REGEX_SERIE
	re := regexp.MustCompile(`\{(\w+)}`)
	result := re.ReplaceAllStringFunc(format, func(serie string) string {
		switch serie {
		case "{name}":
			return m.name
		case "{season}":
			return fmt.Sprintf("%s", oneToNine(m.season))
		case "{episode}":
			return fmt.Sprintf("%s", oneToNine(m.episode))
		case "{resolution}":
			return m.resolution
		case "{year}":
			if m.year == 0 {
				return ""
			}
			return fmt.Sprintf("%d", m.year)
		default:
			return serie
		}
	})

	result = strings.ReplaceAll(result, " - ", " ")
	result = strings.ReplaceAll(result, "- ", " ")
	result = strings.ReplaceAll(result, "()", "")
	result = strings.TrimSpace(result)

	m.complete = result + m.ext
	m.name = result
	m.serieNumber = fmt.Sprintf("s%se%s", oneToNine(m.season), oneToNine(m.episode))
}

func oneToNine(number int) string {
	if number < 10 {
		return "0" + strconv.Itoa(number)
	}
	return strconv.Itoa(number)
}

func (m *myFile) formatageMovie() {
	// constants.REGEX_MOVIE
	format := constants.REGEX_MOVIE
	re := regexp.MustCompile(`\{(\w+)}`)
	result := re.ReplaceAllStringFunc(format, func(movie string) string {
		switch movie {
		case "{name}":
			return m.name
		case "{resolution}":
			return m.resolution
		case "{year}":
			if m.year == 0 {
				return ""
			}
			return fmt.Sprintf("%d", m.year)
		default:
			return movie
		}
	})

	result = strings.ReplaceAll(result, " - ", " ")
	result = strings.ReplaceAll(result, "- ", " ")
	result = strings.ReplaceAll(result, "()", "")
	result = strings.TrimSpace(result)

	m.complete = result + m.ext
	m.name = result
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
	err := syscall.Rename(filePathOld, strings.ToLower(filePathNew))
	if err != nil {
		logger.L(logger.Red, "Failed Rename file => %s", filePathOld)
		cmd := exec.Command("/bin/sh", "-c", "mv \""+filePathOld+"\" "+filePathNew)
		logger.L(logger.Yellow, "mv \""+filePathOld+"\" "+filePathNew)
		err = cmd.Run()
		if err != nil {
			logger.L(logger.Red, "Move Or Rename File : %s", err)
			mu.Unlock()
			return false
		}
	} else {
		logger.L(logger.Yellow, "Rename file => %s", filePathOld)
	}

	folder := filepath.Dir(filePathOld)

	folder = getAbsolutePathWithRelative(folder)
	absoluteATrier := getAbsolutePathWithRelative(constants.A_TRIER)

	if folder != absoluteATrier {
		file, _ := os.ReadDir(folder)
		if len(file) == 0 {
			err = watch.Remove(folder)
			if err != nil {
				logger.L(logger.Red, "Error. Can't delete watcher to folder: %s", folder)
			}
			logger.L(logger.Yellow, "Delete watcher to folder: %s", folder)
			err := os.Remove(folder)
			if err != nil {
				logger.L(logger.Red, "Error to delete folder: %s", folder)
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
