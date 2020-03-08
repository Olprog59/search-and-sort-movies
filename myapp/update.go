package myapp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func ticker() {
	tick := time.NewTicker(10 * time.Second)
	go func() {
		for range tick.C {
			getVersionOnline()
			same := checkIfNewVersion()
			if same {
				log.Println("démarrage de la mise à jour")
				log.Println("il faut couper le ticker après avoir download l'app de mise à jour")
				if downloadApp() {
					log.Println("Ca y est c'est dl!!")
					tick.Stop()
					executeUpdate()
					os.Exit(0)
				}
			}
		}
	}()
}

const UrlUpdateURL = "http://sokys.ddns.net:9999"
const FileUpdateName = "updateSearchAndSortMovies"

//func getVersionOnline() {
//	resp, err := http.Get(UrlUpdateURL)
//	if err != nil {
//		log.Println("Problème de http get pour l'update de l'app\n: " + err.Error())
//	}
//
//	defer resp.Body.Close()
//
//	_ = json.NewDecoder(resp.Body).Decode(&app)
//
//	log.Println(app)
//}

var buildInfo BuildInfo

func getVersionOnline() {
	url := UrlUpdateURL + "/version?file=" + app.Name
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		ticker()
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&buildInfo)
	app.Version = buildInfo.BuildVersion
	log.Println(&app)
	log.Println(&buildInfo)
}

func executeUpdate() {
	cmd := exec.Command("nohup", "./"+FileUpdateName, "&")
	err := cmd.Start()
	log.Println(cmd)
	if err != nil {
		log.Println("Erreur à l'éxécution de 'searchAndSortMoviesUpdate' !!!")
	}
}

func downloadApp() bool {
	fileUrl := UrlUpdateURL + "/update?file=" + FileUpdateName
	if err := downloadAppUpdate("searchAndSortMoviesUpdate", fileUrl); err != nil {
		log.Println("Problème de téléchargement de l'application d'update")
		return false
	}
	return true
}

func downloadAppUpdate(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
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

	// Write the body to file
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
		return true
	}

	log.Println("Pas de mise à jour")
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
