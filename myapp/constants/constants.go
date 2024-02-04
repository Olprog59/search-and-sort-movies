package constants

import (
	"os"
	"strconv"
)

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
	UID         = GetEnvInt("UID", "0")
	GID         = GetEnvInt("GID", "0")
	CHMOD       = GetEnvUInt32("CHMOD", "0755")
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvBool(key string) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}
	return false
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

func GetEnvUInt32(key, fallback string) uint32 {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint32(i)
		}
	}
	if i, err := strconv.ParseUint(fallback, 10, 32); err == nil {
		return uint32(i)
	}
	return 0
}
