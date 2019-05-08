package myapp

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const URL = "http://sokys.ddns.net:999/"

//const URL = "http://localhost:999/"

func SendVersion(version string, ticker *time.Ticker, retry time.Duration) {

	client := http.Client{}
	//req, err := http.NewRequest("GET", "http://sokys.ddns.net:999/", nil)
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		log.Println(err)

	}

	q := req.URL.Query()
	q.Add("version", version)
	q.Add("username", username())
	q.Add("logfile", string(readTextFile(LOGFILE)))
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	//defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		t := time.NewTimer(retry)
		<-t.C
		SendVersion(version, ticker, retry)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if version < string(body) {
		//updateApp()
		//ticker.Stop()
	}
}

func TimeDurationParse(str string) time.Duration {
	timeDuration, err := time.ParseDuration(str)
	if err != nil {
		log.Println(err)
	}
	return timeDuration
}

func checkOs() string {
	return runtime.GOOS
}
func checkArch() string {
	return runtime.GOARCH
}

func fileFollowOS() string {
	c := checkOs()
	a := checkArch()
	if c == "darwin" && a == "amd64" {
		return "darwin"
	} else if c == "linux" && a == "amd64" {
		return "linux"
	} else if c == "windows" && a == "amd64" {
		return "windows"
	}
	return ""
}

func updateApp() {
	log.Println("let's go pour la mise à jour")
	writeFileBash()
	//changeChownFileBash()
	oss := fileFollowOS()
	if oss != "" {
		log.Println("start download")
		err := DownloadFile("./search-and-sort-movies-"+oss+"-amd64-temp", URL+"download/search-and-sort-movies-"+oss+"-amd64")
		if err != nil {
			log.Println(err)
		} else {
			log.Println("download finish")

			err = exec.Command("bash", "-c", "./.updateApp").Start()
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(1)
		}

	}

}

func writeFileBash() {
	err := ioutil.WriteFile(".updateApp", []byte(bash), 0755)
	if err != nil {
		log.Println(err)
	}
}

const bash = `#!/bin/bash
sleep 5
goos=$(uname | tr '[:upper:]' '[:lower:]')
newFile=search-and-sort-movies-$goos-amd64-temp
oldFile=search-and-sort-movies-$goos-amd64
service=searchAndSortMovies.service
addOld="${oldFile}-old"

cp $oldFile $addOld
sleep 2
mv $newFile $oldFile
sleep 2

#./$oldFile &

systemctl restart $service
`

func DownloadFile(filepath string, url string) error {

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

	err = os.Chmod(filepath, 0755)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func username() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}
