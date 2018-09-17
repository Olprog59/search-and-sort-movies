package myapp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// MoviesExcept :
type MoviesExcept struct {
	Name string `json:"name"`
}

const (
	// MoviesExceptFile :
	MoviesExceptFile = ".movies_except.json"
)

// GetEnv :
func GetMoviesExceptFile(value string) bool {
	checkIfMovieExceptFileIsExist()
	jsonType := readFile()

	for _, v := range jsonType {
		if v.Name == value {
			return true
		}
	}

	return false
}

func RemoveMoviesExceptFile(value string) {
	jsonType := readFile()
	var newJson []MoviesExcept
	for _, v := range jsonType {
		if v.Name != value {
			newJson = append(newJson, v)
		}
	}
	j, err := json.MarshalIndent(newJson, "", " ")
	if err != nil {
		log.Println(err)
	}
	writeFile(j)
}

// SetEnv :
func SetMoviesExceptFile(value string) {
	checkIfMovieExceptFileIsExist()
	jsonType := readFile()

	if !GetMoviesExceptFile(value) {
		jsonType = append(jsonType, MoviesExcept{Name: value})
		j, err := json.MarshalIndent(jsonType, "", " ")
		if err != nil {
			log.Println(err)
		}

		writeFile(j)
	}
}

func readFile() []MoviesExcept {
	checkIfMovieExceptFileIsExist()
	f, err := ioutil.ReadFile(MoviesExceptFile)

	if err != nil {
		log.Println(err)
	}
	var jsonType []MoviesExcept

	json.Unmarshal(f, &jsonType)

	return jsonType
}

func writeFile(jsonByte []byte) {
	checkIfMovieExceptFileIsExist()
	err := ioutil.WriteFile(MoviesExceptFile, jsonByte, 0644)
	if err != nil {
		log.Println(err)
	}
}

func checkIfMovieExceptFileIsExist() {
	// detect if file exists
	var _, err = os.Stat(MoviesExceptFile)

	// create file if not exists
	if os.IsNotExist(err) {
		newJSON := &MoviesExcept{
			Name: "",
		}
		j, err := json.MarshalIndent(newJSON, "", " ")
		if err != nil {
			log.Println(err)
		}

		writeJSONFile(j)
	}
}
