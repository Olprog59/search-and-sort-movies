package myapp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Config :
type Config struct {
	Dlna   string `json:"dlna"`
	Series string `json:"series"`
	Movies string `json:"movies"`
	// Music  string `json:"music"`
}

const (
	// ConfigFile :
	ConfigFile = ".config.json"
)

// GetEnv :
func GetEnv(key string) string {
	checkIfConfigFileIsExist()

	jsonType := readJSONFile()

	return jsonType[key].(string)
}

// SetEnv :
func SetEnv(key, value string) {
	checkIfConfigFileIsExist()
	// open file using READ & WRITE permission
	jsonType := readJSONFile()

	jsonType[key] = value

	j, err := json.MarshalIndent(jsonType, "", " ")
	if err != nil {
		log.Println(err)
	}

	writeJSONFile(j)
}

// CheckIfConfigFileIsExist : Create file is not exist
func checkIfConfigFileIsExist() {
	// detect if file exists
	var _, err = os.Stat(ConfigFile)

	// create file if not exists
	if os.IsNotExist(err) {
		newJSON := &Config{
			Dlna:   "", //pwd("dlna", true),
			Series: "", //pwd("dlna/Series", true),
			Movies: "", //pwd("dlna/Movies", true),
			// Music:  pwd("dlna/Music", true),
		}
		j, err := json.MarshalIndent(newJSON, "", " ")
		if err != nil {
			log.Println(err)
		}

		writeJSONFile(j)
	}
}

// ReadJSONFile :
func readJSONFile() map[string]interface{} {
	f, err := ioutil.ReadFile(ConfigFile)

	if err != nil {
		log.Println(err)
	}
	var jsonType map[string]interface{}

	json.Unmarshal(f, &jsonType)

	return jsonType
}

func writeJSONFile(jsonByte []byte) {
	err := ioutil.WriteFile(ConfigFile, jsonByte, 0644)
	// file, err := os.Create(ConfigFile)
	if err != nil {
		log.Println(err)
	}
}

func pwd(name string, endPathSeparator bool) string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	if endPathSeparator {
		return filepath.Clean(pwd+string(os.PathSeparator)+name) + string(os.PathSeparator)
	}
	return filepath.Clean(pwd + string(os.PathSeparator) + name)
}

func StartScan(auto bool) {
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

func FirstConnect() bool {
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

func FirstConfig() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Hello, bienvenue sur l'application de tri des vidéos.")
		fmt.Println("Ceci est ta première connexion donc il faut configurer des petites choses.")
		pwd, _ := os.Getwd()
		fmt.Println("A savoir, que tu te trouves dans le répertoire : " + pwd)
		fmt.Println("Commençons par indiqué l'emplacement des fichiers à trier : (ex : /home/user/a_trier ou windows : C:\\Users\\Dupont\\Documents\\a_trier)")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		SetEnv("dlna", path.Clean(strings.TrimSpace(text)))
		fmt.Println("Ensuite, il faut indiqué le dossier ou tu veux mettre tes films : (ex : /mnt/dlna/movies ou windows : F:\\films)")
		text, _ = reader.ReadString('\n')
		fmt.Println(text)
		SetEnv("movies", path.Clean(strings.TrimSpace(text)))
		fmt.Println("Pour finir, il faut indiqué le dossier ou tu veux mettre tes séries : (ex : /mnt/dlna/series ou windows : F:\\series)")
		text, _ = reader.ReadString('\n')
		fmt.Println(text)
		SetEnv("series", path.Clean(strings.TrimSpace(text)))
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

func CheckFolderExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}
}

func fileInFolder() (int, []os.FileInfo) {
	files, err := ioutil.ReadDir(GetEnv("dlna"))
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
			Process(f.Name())
		}
	}
	log.Println("Tri terminé !")
}
