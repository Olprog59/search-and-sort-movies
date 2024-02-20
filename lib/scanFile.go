package lib

import (
	"github.com/sam-docker/media-organizer/constants"
	"github.com/sam-docker/media-organizer/logger"
	"github.com/sam-docker/media-organizer/model"
	"os"
	"path/filepath"
	"regexp"
)

func StartScan(obsSlice *model.ObservableSlice) {
	logger.Debug("Scan files in folder(s) : %s", constants.BE_SORTED)
	err := filepath.Walk(constants.BE_SORTED+"/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		re := regexp.MustCompile(constants.RegexFileExtension)
		if re.MatchString(filepath.Ext(path)) {
			duration, err := GetMediaDuration(path)
			if err != nil {
				logger.Err("Error checking media file: %s", err)
				return nil
			}
			file := model.SliceFile{File: path, Working: false, Duration: duration}
			obsSlice.Add(file)
		}
		return nil
	})
	if err != nil {
		logger.Err("walk error %s", err)
	}
}
