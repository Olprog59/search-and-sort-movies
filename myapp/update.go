package myapp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Application struct {
	Version    string `json:"version"`
	OldVersion string `json:"old_version"`
	Name       string `json:"name"`
}

var app Application

func LaunchAppCheckUpdate(oldVersion string, name string) {
	app.OldVersion = oldVersion
	app.Name = name
	ticker()
}

const duration = 1 * time.Hour

func ticker() {
	tick := time.NewTicker(duration)
	go func() {
		for range tick.C {
			removeFileUpdate()
			checkIfSiteIsOnline()
			getVersionOnline()
			same := checkIfNewVersion()
			if same {
				log.Println("démarrage de la mise à jour")
				if downloadApp() {
					log.Println("Ca y est c'est dl!!")
					executeUpdate()
					os.Exit(0)
				}
			}
		}
	}()
}

const UrlUpdateURL = "http://sokys.ddns.net:9999"
const FileUpdateName = "updateSearchAndSortMovies-" + runtime.GOOS

var buildInfo BuildInfo

func removeFileUpdate() {
	_, err := os.Stat(FileUpdateName)
	if err != nil {
		return
	}
	if err = os.Remove(FileUpdateName); err != nil {
		log.Println(err)
	}
}

func getVersionOnline() {
	url := UrlUpdateURL + "/version?file=" + app.Name
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&buildInfo)
	app.Version = buildInfo.BuildVersion
}

func checkIfSiteIsOnline() {
	_, err := http.Get(UrlUpdateURL)
	if err != nil {
		log.Println("Le site n'est pas accessible. Un nouveau test se fera dans 1 minute")
		time.Sleep(1 * time.Minute)
		checkIfSiteIsOnline()
	}
}

var _count int64
func downloadApp() bool {
	fileUrl := UrlUpdateURL + "/update?file=" + FileUpdateName
	if err := downloadAppUpdate(FileUpdateName, fileUrl); err != nil {
		log.Println("Problème de téléchargement de l'application d'update")
		if _count < 2 {
			downloadApp()
		}
		_count++
		return false
	}
	return true
}

func downloadAppUpdate(filepath string, url string) error {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	err = out.Chmod(0755)
	return err
}

func checkIfNewVersion() bool {
	var oldV, newV int64
	if strToInt64(app.OldVersion) != 0 {
		oldV = strToInt64(app.OldVersion)
	}
	if strToInt64(app.Version) != 0 {
		newV = strToInt64(app.Version)
	}
	if newV > oldV {
		log.Println("il y a une mise à jour")
		log.Printf("\n    - Ancienne version: %s\n    - Nouvelle version: %s\n\n", app.OldVersion, app.Version)
		return true
	}
	//log.Println("Pas de mise à jour")
	return false

}

func strToInt64(version string) (vv int64) {
	tab := strings.Split(version, ".")
	j := strings.Join(tab, "")
	vv, err = strconv.ParseInt(j, 10, 64)
	if err != nil {
		return 0
	}
	return vv
}
