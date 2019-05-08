package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"search-and-sort-movies/myapp"
	"strings"
	"time"
)

var (
	BuildVersion string
	BuildHash    string
	BuildDate    string
	BuildClean   string
	BuildName    = "search-and-sort-movies"
	count        int
	ticker       *time.Ticker
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	myapp.Flags(BuildName, BuildVersion, BuildHash, BuildDate, BuildClean)
	// Write log to file : log_SearchAndSort
	f, err := os.OpenFile(myapp.LOGFILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Check if it's the first connection
	if firstConnect() {
		firstConfig()
	} else {
		for {
			if myapp.GetEnv("dlna") == "" || myapp.GetEnv("movies") == "" || myapp.GetEnv("series") == "" {
				firstConfig()
			} else {
				break
			}
		}
	}

	if myapp.GetEnv("check_update") == "" {
		myapp.SetEnv("check_update", (24 * time.Hour).String())
	}

	if myapp.GetEnv("retry_check_update") == "" {
		myapp.SetEnv("retry_check_update", (12 * time.Hour).String())
	}

	ticker = time.NewTicker(myapp.TimeDurationParse(myapp.GetEnv("retry_check_update")))

	go func() {
		for range ticker.C {
			myapp.SendVersion(BuildVersion, ticker, myapp.TimeDurationParse(myapp.GetEnv("retry_check_update")))
		}
	}()

	checkFolderExists(myapp.GetEnv("dlna"))
	checkFolderExists(myapp.GetEnv("movies"))
	checkFolderExists(myapp.GetEnv("series"))

	log.Println("Start :-D")
	log.Println("Ecoute sur le dossier : " + myapp.GetEnv("dlna"))
	fmt.Println("Start :-D")
	fmt.Println("Ecoute sur le dossier : " + myapp.GetEnv("dlna"))

	myapp.MyWatcher(myapp.GetEnv("dlna"))

}

func checkFolderExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}
}

func firstConnect() bool {
	_, err := os.Stat(myapp.ConfigFile)

	if os.IsNotExist(err) {
		log.Println(err)
		return true
	}
	return false
}

func readJSONFileConsole() {
	f, err := ioutil.ReadFile(myapp.ConfigFile)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%s\n", string(f))
}

func firstConfig() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Hello, bienvenue sur l'application de tri des vidéos.")
		fmt.Println("Ceci est ta première connexion donc il faut configurer des petites choses.")
		pwd, _ := os.Getwd()
		fmt.Println("A savoir, que tu te trouves dans le répertoire : " + pwd)
		fmt.Println("Commençons par indiqué l'emplacement des fichiers à trier : (ex : /home/user/a_trier ou windows : C:\\Users\\Dupont\\Documents\\a_trier)")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		myapp.SetEnv("dlna", path.Clean(strings.TrimSpace(text)))
		fmt.Println("Ensuite, il faut indiqué le dossier ou tu veux mettre tes films : (ex : /mnt/dlna/movies ou windows : F:\\films)")
		text, _ = reader.ReadString('\n')
		fmt.Println(text)
		myapp.SetEnv("movies", path.Clean(strings.TrimSpace(text)))
		fmt.Println("Ensuite, il faut indiqué le dossier ou tu veux mettre tes séries : (ex : /mnt/dlna/series ou windows : F:\\series)")
		text, _ = reader.ReadString('\n')
		fmt.Println(text)
		myapp.SetEnv("series", path.Clean(strings.TrimSpace(text)))
		//fmt.Println("Et pour finir un email. Celui-ci permettra d'envoyer un email si un problème est rencontré.")
		//text, _ = reader.ReadString('\n')
		//fmt.Println(text)
		//myapp.SetEnv("email", path.Clean(strings.TrimSpace(text)))

		fmt.Println("Super. Voilà tout est configuré. On va vérifier le fichier : ")
		fmt.Println('\n')
		readJSONFileConsole()
		fmt.Println("Est-ce que cela est correct? (O/n)")
		text, _ = reader.ReadString('\n')
		if strings.TrimSpace(text) == "n" || strings.TrimSpace(text) == "N" {
			firstConfig()
		} else {
			fmt.Println("Juste une petite chose. Si tu veux changer de répertoire, tu peux le faire dans le fichier caché : .config.json dans le répertoire searchMoviesConfig")
			break
		}

	}

	fmt.Println("Cool!!! C'est parti. Enjoy")
}
