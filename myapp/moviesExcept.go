package myapp

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	jsonType := readFile()

	for _, v := range jsonType {
		if v.Name == value {
			return true
		}
	}

	return false
}

// SetEnv :
func SetMoviesExceptFile(value string) {
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
	f, err := ioutil.ReadFile(MoviesExceptFile)

	if err != nil {
		log.Println(err)
	}
	var jsonType []MoviesExcept

	json.Unmarshal(f, &jsonType)

	return jsonType
}

func writeFile(jsonByte []byte) {
	err := ioutil.WriteFile(MoviesExceptFile, jsonByte, 0644)
	if err != nil {
		log.Println(err)
	}
}
