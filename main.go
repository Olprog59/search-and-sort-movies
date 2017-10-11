package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"search-and-sort-movies/myapp"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

var (
	buildVersion string
	buildHash    string
	buildDate    string
	buildClean   string
	buildName    = "search-and-sort-movies"
)

func main() {

	vers := flag.Bool("v", false, "Indique la version de l'application")
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	windows := flag.Bool("windows", false, "Lancer l'application sans l'invite de commandes")
	flag.Parse()

	if *vers {
		// flag.PrintDefaults()
		fmt.Printf("Name: %s\n", buildName)
		fmt.Printf("Version: %s\n", buildVersion)
		fmt.Printf("Git Commit Hash: %s\n", buildHash)
		fmt.Printf("Build Date: %s\n", buildDate)
		fmt.Printf("Built from clean source tree: %s\n", buildClean)
		fmt.Printf("OS: %s\n", runtime.GOOS)
		fmt.Printf("Architecture: %s\n", runtime.GOARCH)
		os.Exit(1)
	}

	if *scan {
		startScan(true)
	}

	if *windows {
		myapp.HiddenWindow()
	}

	// Write log to file : log_SearchAndSort
	f, err := os.OpenFile("log_SearchAndSort", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	checkFolderExists(myapp.GetEnv("dlna"))
	checkFolderExists(myapp.GetEnv("movies"))
	checkFolderExists(myapp.GetEnv("series"))

	fmt.Println("Start :-D")

	// startScan(false)

	fmt.Println("Ecoute sur le dossier : " + myapp.GetEnv("dlna"))
	myapp.Watcher(myapp.GetEnv("dlna"))

}

func startScan(auto bool) {
	if count, file := fileInFolder(); count > 0 {
		if auto {
			fmt.Println("Scan automatique")
			go boucleFiles(file)
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Je vois qu'il y a des fichiers vidéos actuellement dans ton dossier source.")
			fmt.Println("Veux tu faire le tri? (O/n)")
			text, _ := reader.ReadString('\n')
			fmt.Println(text)
			if strings.TrimSpace(text) == "n" || strings.TrimSpace(text) == "N" {
				return
			}

			go boucleFiles(file)
		}
	}
}

func firstConnect() bool {
	_, err := os.Stat(".config.json")

	if os.IsNotExist(err) {
		log.Println(err)
		return true
	}
	return false
}

func readJSONFileConsole() {
	f, err := ioutil.ReadFile(".config.json")

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
		fmt.Println("Pour finir, il faut indiqué le dossier ou tu veux mettre tes séries : (ex : /mnt/dlna/series ou windows : F:\\series)")
		text, _ = reader.ReadString('\n')
		fmt.Println(text)
		myapp.SetEnv("series", path.Clean(strings.TrimSpace(text)))
		fmt.Println("Pour la musique, il faut attendre les prochaines versions. :-(  ")

		fmt.Println("Super. Voilà tout est configuré. On va vérifier le fichier : ")
		fmt.Println('\n')
		readJSONFileConsole()
		fmt.Println("Est-ce que cela est correct? (O/n)")
		text, _ = reader.ReadString('\n')
		if strings.TrimSpace(text) == "n" || strings.TrimSpace(text) == "N" {
			continue
		} else {
			break
		}

	}

	fmt.Println("Cool!!! C'est parti. Enjoy")
}

func checkFolderExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}
}

func fileInFolder() (int, []os.FileInfo) {
	files, err := ioutil.ReadDir(myapp.GetEnv("dlna"))
	if err != nil {
		log.Fatal(err)
	}

	var count int
	for _, f := range files {
		if !f.IsDir() {
			re := regexp.MustCompile(`(.mkv|.mp4|.avi|.flv)`)
			if re.MatchString(filepath.Ext(f.Name())) {
				count++
			}
		}
	}
	return count, files
}

func boucleFiles(files []os.FileInfo) {
	log.Println("Démarrage du tri !")
	for _, f := range files {
		if !f.IsDir() {
			log.Println("Movies : " + f.Name())
			myapp.Process(f.Name())
		}
	}
	log.Println("Tri terminé !")
}
