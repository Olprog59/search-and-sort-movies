package myapp

import (
	"os"
	"runtime"
	"time"
)

const (
	apiV4 = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJlYTg3Nzk2MzhmMDc4ZjI1ZGFhMzkxM2U4MGZlNDZlYiIsInN1YiI6IjU5Y2Y3NjdiYzNhMzY4MWViMTAxOThjNyIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.mAxfQbzn4WIft74XAooGGiw7PhHxMNTu8TtTvPwhh1c"
	apiV3 = "ea8779638f078f25daa3913e80fe46eb"

	regexFile = `(.mkv|.mp4|.avi|.flv)`

	// Dev
	DURATION                = 1 * time.Minute
	DurationRetryConnection = 1 * time.Minute
	DurationRetryDownload   = 1 * time.Minute
	UrlUpdateURL            = "http://localhost:9999"

	// Prod
	//DURATION                = 5 * time.Hour
	//DurationRetryConnection = 1 * time.Hour
	//DurationRetryDownload   = 1 * time.Hour
	//UrlUpdateURL            = "http://sokys.ddns.net:9999"

	FileUpdateName = "updateSearchAndSortMovies-" + runtime.GOOS
	FolderConfig   = "./searchMoviesConfig"

	LOGFILE    = FolderConfig + string(os.PathSeparator) + "log_SearchAndSort"
	ConfigFile = FolderConfig + string(os.PathSeparator) + ".config.json"
)
