package myapp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"search-and-sort-movies/myapp/constants"
)

// Config :
type Config struct {
	Dlna   string `json:"dlna"`
	Series string `json:"series"`
	Movies string `json:"movies"`
	Port   string `json:"port"`
}

// GetEnv :
func GetEnv(key string) string {
	checkIfConfigFileIsExist()

	jsonType := readJSONFile(constants.ConfigFile)

	if jsonType[key] == nil {
		SetEnv(key, "")
	}

	jsonType = readJSONFile(constants.ConfigFile)

	return jsonType[key].(string)
}

// SetEnv :
func SetEnv(key, value string) {
	checkIfConfigFileIsExist()
	// open file using READ & WRITE permission
	jsonType := readJSONFile(constants.ConfigFile)

	jsonType[key] = value

	j, err := json.MarshalIndent(jsonType, "", " ")
	if err != nil {
		log.Println(err)
	}

	writeJSONFile(constants.ConfigFile, j)
}

// CheckIfConfigFileIsExist : Create file is not exist
func checkIfConfigFileIsExist() {
	if _, err := os.Stat(constants.FolderConfig); os.IsNotExist(err) {
		_ = os.Mkdir(constants.FolderConfig, os.ModeSticky|0755)
	}

	// detect if file exists
	var _, err = os.Stat(constants.ConfigFile)

	// create file if not exists
	if os.IsNotExist(err) {
		newJSON := &Config{
			Dlna:   "", //pwd("dlna", true),
			Series: "", //pwd("dlna/Series", true),
			Movies: "", //pwd("dlna/Movies", true),
		}
		j, err := json.MarshalIndent(newJSON, "", " ")
		if err != nil {
			log.Println(err)
		}

		writeJSONFile(constants.ConfigFile, j)
	}
}

// ReadJSONFile :
func readJSONFile(file string) map[string]interface{} {
	f, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println(err)
	}
	var jsonType map[string]interface{}

	_ = json.Unmarshal(f, &jsonType)

	return jsonType
}

func writeJSONFile(file string, jsonByte []byte) {
	err := ioutil.WriteFile(file, jsonByte, 0644)
	// file, err := os.Create(ConfigFile)
	if err != nil {
		log.Println(err)
	}
}
