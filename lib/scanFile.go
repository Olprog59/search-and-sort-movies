package lib

import (
	"github.com/sam-docker/media-organizer/constants"
	"github.com/sam-docker/media-organizer/logger"
	"os"
	"path/filepath"
	"regexp"
)

func StartScan() {
	count, file := fileInFolder()
	if count > 0 {
		go boucleFiles(file)
	}
}

func fileInFolder() (int, []string) {
	var files []string
	var count int
	logger.L(logger.Purple, "%s", constants.BE_SORTED)
	err := filepath.Walk(constants.BE_SORTED+"/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		re := regexp.MustCompile(constants.RegexFileExtension)
		if re.MatchString(filepath.Ext(path)) {
			files = append(files, path)
			count++
		}
		return nil
	})
	if err != nil {
		logger.L(logger.Red, "walk error %s", err)
	}
	return count, files
}

func boucleFiles(files []string) {
	logger.L(logger.Purple, "Start sorting videos !!")
	for _, f := range files {
		logger.L(logger.Yellow, "file : "+f)
		var m myFile
		m.file = f
		m.Process()
	}
	logger.L(logger.Purple, "Sorting completed !")
}
