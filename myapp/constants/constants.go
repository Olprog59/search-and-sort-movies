package constants

import "os"

const (
	RegexFile = `(?i)(.mkv|.mp4|.avi|.flv|.mov)`

	A_TRIER = path + "/be_sorted"
	MOVIES  = path + "/movies"
	SERIES  = path + "/series"
)

// Premier élément: s == " " sinon -, _ ou .
var FormatFile = func() string {
	env := os.Getenv("FORMAT_FILE")
	if env == "" {
		//return "  - , name"
		// a des fins de test
		return "  - , name, resolution, year"
	}
	return env
}
