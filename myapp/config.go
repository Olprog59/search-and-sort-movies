package myapp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config :
type Config struct {
	Dlna             string `json:"dlna"`
	Series           string `json:"series"`
	Movies           string `json:"movies"`
	CheckUpdate      string `json:"check_update"`
	RetryCheckUpdate string `json:"retry_check_update"`
	User             string `json:"user"`
	Group            string `json:"group"`
}

// GetEnv :
func GetEnv(key string) string {
	checkIfConfigFileIsExist()

	jsonType := readJSONFile(ConfigFile)

	if jsonType[key] == nil {
		SetEnv(key, "")
	}

	jsonType = readJSONFile(ConfigFile)

	return jsonType[key].(string)
}

// SetEnv :
func SetEnv(key, value string) {
	checkIfConfigFileIsExist()
	// open file using READ & WRITE permission
	jsonType := readJSONFile(ConfigFile)

	jsonType[key] = value

	j, err := json.MarshalIndent(jsonType, "", " ")
	if err != nil {
		log.Println(err)
	}

	writeJSONFile(ConfigFile, j)
}

// CheckIfConfigFileIsExist : Create file is not exist
func checkIfConfigFileIsExist() {
	if _, err := os.Stat(FolderConfig); os.IsNotExist(err) {
		os.Mkdir(FolderConfig, os.ModeSticky|0755)
	}

	// detect if file exists
	var _, err = os.Stat(ConfigFile)

	// create file if not exists
	if os.IsNotExist(err) {
		newJSON := &Config{
			Dlna:             "", //pwd("dlna", true),
			Series:           "", //pwd("dlna/Series", true),
			Movies:           "", //pwd("dlna/Movies", true),
			CheckUpdate:      "24h0m0s",
			RetryCheckUpdate: "12h0m0s",
			// Music:  pwd("dlna/Music", true),
		}
		j, err := json.MarshalIndent(newJSON, "", " ")
		if err != nil {
			log.Println(err)
		}

		writeJSONFile(ConfigFile, j)
	}
}

// ReadJSONFile :
func readJSONFile(file string) map[string]interface{} {
	f, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println(err)
	}
	var jsonType map[string]interface{}

	json.Unmarshal(f, &jsonType)

	return jsonType
}

// ReadJSONFile :
func readTextFile(file string) []byte {
	f, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println(err)
	}

	return f
}

func writeJSONFile(file string, jsonByte []byte) {
	err := ioutil.WriteFile(file, jsonByte, 0644)
	// file, err := os.Create(ConfigFile)
	if err != nil {
		log.Println(err)
	}
}
