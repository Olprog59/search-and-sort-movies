package myapp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"time"
)

const URL = "http://localhost:999/"

func SendVersion(version string, ticker *time.Ticker) {
	client := http.Client{}
	//req, err := http.NewRequest("GET", "http://sokys.ddns.net:999/", nil)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println(err)

	}

	q := req.URL.Query()
	q.Add("version", version)
	q.Add("username", username())
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		count++
		log.Println(err)
		t := time.NewTimer(5 * time.Second)
		<-t.C
		SendVersion(version, ticker)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if version < string(body) {
		updateApp()
		ticker.Stop()
	}
	//if BuildVersion < string(body){
	//	updateApp()
	//}
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
	fmt.Println("let's go pour la mise à jour")
	oss := fileFollowOS()
	if oss != "" {
		err := DownloadFile("/Users/olprog/go/src/search-and-sort-movies/a_trier/search-and-sort-movies-"+oss+"-amd64", URL+"download/search-and-sort-movies-"+oss+"-amd64-temp")
		if err != nil {
			log.Println(err)
		}

	}

}

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

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func username() string {
	user2, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user2.Name
}
