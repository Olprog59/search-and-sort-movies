package myapp

import (
	"os"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
)

func StartScan() {
	if count, file := fileInFolder(); count > 0 {
		//boucleFiles(file)
		// remove goroutine car je dois tester voir si cela cause le tri non complet des fichiers
		go boucleFiles(file)
	}
}

func fileInFolder() (int, []string) {
	var files []string
	var count int
	err := filepath.Walk(constants.A_TRIER, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		re := regexp.MustCompile(constants.RegexFile)
		if re.MatchString(filepath.Ext(path)) {
			files = append(files, path)
			count++
		}
		return nil
	})
	if err != nil {
		logger.L(logger.Red, "walk error [%v]", err)

	}
	return count, files
}

func boucleFiles(files []string) {
	logger.L(logger.Purple, "Démarrage du tri !")
	for _, f := range files {
		logger.L(logger.Yellow, "File : "+f)
		var m myFile
		m.file = f
		m.Process()
	}
	logger.L(logger.Purple, "Tri terminé !")
}
