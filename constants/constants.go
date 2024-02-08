package constants

import (
	"github.com/sam-docker/media-organizer/model"
	"os"
	"strconv"
)

const (
	RegexFileExtension = `(?i)(.mkv|.mp4|.avi|.flv|.mov)`
)

var (
	ObsSlice    = model.NewObservableSlice()
	BE_SORTED   = GetEnv("BE_SORTED", "/medias/be_sorted")
	MOVIES      = GetEnv("MOVIES", "/medias/movies")
	SERIES      = GetEnv("SERIES", "/medias/series")
	REGEX_MOVIE = GetEnv("REGEX_MOVIE", "")
	REGEX_SERIE = GetEnv("REGEX_SERIE", "")
	UID         = GetEnvInt("UID", "0")
	GID         = GetEnvInt("GID", "0")
	CHMOD       = GetEnv("CHMOD", "0755")
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key, fallback string) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	if i, err := strconv.Atoi(fallback); err == nil {
		return i
	}
	return 0
}
