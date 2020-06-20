package update

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"search-and-sort-movies/myapp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/model"
	"strconv"
	"strings"
	"time"
)

type Application struct {
	Version    string `json:"version"`
	OldVersion string `json:"old_version"`
	Name       string `json:"name"`
}

var _firstStart = true

func (a *Application) LaunchAppCheckUpdate(oldVersion string, name string) {
	a.OldVersion = oldVersion
	a.Name = name
	a.ticker()
}

func (a *Application) ticker() {
	if _firstStart {
		a.operationAll()
		_firstStart = false
	}
	tick := time.NewTicker(constants.DURATION)
	go func() {
		for range tick.C {
			a.operationAll()
		}
	}()
}

func (a *Application) operationAll() {
	// envoie des infos en post
	//go send()

	removeFileUpdate()
	checkIfSiteIsOnline()
	a.getVersionOnline()
	different := a.checkIfNewVersion()
	if different {
		log.Println("démarrage de la mise à jour")
		// Début du dl du logiciel de mise à jour
		if downloadApp() {
			log.Println("Ca y est c'est dl!!")
			executeUpdate()
			os.Exit(0)
		}
	} else {
		go myapp.PostInfo(a.OldVersion)
	}
}

var buildInfo model.BuildInfo

func removeFileUpdate() {
	_, err := os.Stat(constants.FileUpdateName)
	if err != nil {
		return
	}
	if err = os.Remove(constants.FileUpdateName); err != nil {
		log.Println(err)
	}
}

func (a *Application) getVersionOnline() {
	url := constants.UrlUpdateURL + "/version?file=" + a.Name
	var netClient = &http.Client{
		Timeout: time.Second * 20,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		log.Println(err)
		time.Sleep(time.Minute * 20)
		a.getVersionOnline()
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&buildInfo)
	a.Version = buildInfo.BuildVersion
}

func checkIfSiteIsOnline() {
	_, err := http.Get(constants.UrlUpdateURL)
	if err != nil {
		log.Printf("Le site n'est pas accessible. Un nouveau test se fera dans %s", constants.DurationRetryConnection.String())
		time.Sleep(constants.DurationRetryConnection)
		checkIfSiteIsOnline()
		return
	}
}

var _count int64

func downloadApp() bool {
	fileUrl := constants.UrlUpdateURL + "/update?file=" + constants.FileUpdateName
	if err := downloadAppUpdate(constants.FileUpdateName, fileUrl); err != nil {
		log.Println("Problème de téléchargement de l'application d'update")
		if _count < 2 {
			time.Sleep(constants.DurationRetryDownload)
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

func (a *Application) checkIfNewVersion() bool {
	var oldV, newV int64
	if strToInt64(a.OldVersion) != 0 {
		oldV = strToInt64(a.OldVersion)
	}
	if strToInt64(a.Version) != 0 {
		newV = strToInt64(a.Version)
	}
	if newV > oldV {
		log.Println("il y a une mise à jour")
		log.Printf("\n    - Ancienne version: %s\n    - Nouvelle version: %s\n\n", a.OldVersion, a.Version)
		return true
	}
	return false
}

func strToInt64(version string) (vv int64) {
	tab := strings.Split(version, ".")
	// TODO : Correction Temporaire pour la mise à jour de 0.9.1.35 à 0.9.1.36
	if version != "0.9.1.36" {
		if len(tab[1]) == 1 {
			tab[1] = "0" + tab[1]
		}
		if len(tab[2]) == 1 {
			tab[2] = "0" + tab[2]
		}
		if len(tab[3]) == 1 {
			tab[3] = "0" + tab[3]
		}
	}
	j := strings.Join(tab, "")
	vv, err := strconv.ParseInt(j, 10, 64)
	if err != nil {
		return 0
	}
	return vv
}
