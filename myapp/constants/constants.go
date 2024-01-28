package constants

import "os"

const (
	RegexFile = `(?i)(.mkv|.mp4|.avi|.flv|.mov)`
)

var (
	A_TRIER     = GetEnv("A_TRIER", "/mnt/medias/a_trier")
	MOVIES      = GetEnv("MOVIES", "/mnt/medias/movies")
	SERIES      = GetEnv("SERIES", "/mnt/medias/series")
	ALL         = GetEnv("ALL", "")
	REGEX_MOVIE = GetEnv("REGEX_MOVIE", "{name}-{resolution} ({year})")
	REGEX_SERIE = GetEnv("REGEX_SERIE", "{name}-s{season}e{episode}-{resolution} ({year})")
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
