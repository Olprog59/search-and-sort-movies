package myapp

import (
	"github.com/jinzhu/gorm"
)

// MoviesExcept :
type MoviesExcept struct {
	gorm.Model
	Name string `json:"name"`
}

// GetEnv :
func GetMoviesExceptFile(value string) bool {
	var movieExcept MoviesExcept
	testDb(func(db *gorm.DB) {
		db.Where("name = ?", value).First(&movieExcept)
	})
	return movieExcept.Name != ""
}

// GetEnv :
func GetAllExcept() []MoviesExcept {
	var movieExcept []MoviesExcept
	testDb(func(db *gorm.DB) {
		db.Find(&movieExcept)
	})
	return movieExcept
}

//func RemoveMoviesExceptFile(value string) {
//	jsonType := readFile()
//	var newJson []MoviesExcept
//	for _, v := range jsonType {
//		if v.Name != value {
//			newJson = append(newJson, v)
//		}
//	}
//	j, err := json.MarshalIndent(newJson, "", " ")
//	if err != nil {
//		log.Println(err)
//	}
//	writeFile(j)
//}

// SetEnv :
func SetMoviesExceptFile(value string) {
	//checkIfMovieExceptFileIsExist()
	//jsonType := readFile()
	//
	//if !GetMoviesExceptFile(value) {
	//	jsonType = append(jsonType, MoviesExcept{Name: value})
	//	j, err := json.MarshalIndent(jsonType, "", " ")
	//	if err != nil {
	//		log.Println(err)
	//	}
	//
	//	writeFile(j)
	//}
}

//func readFile() []MoviesExcept {
//	f, err := ioutil.ReadFile(MoviesExceptFile)
//
//	if os.IsNotExist(err) {
//		return nil
//	}
//
//	if err != nil {
//		log.Println(err)
//	}
//	var jsonType []MoviesExcept
//
//	json.Unmarshal(f, &jsonType)
//
//	return jsonType
//}
//
//func writeFile(jsonByte []byte) {
//	err := ioutil.WriteFile(MoviesExceptFile, jsonByte, 0644)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func checkIfMovieExceptFileIsExist() {
//	// detect if file exists
//	var _, err = os.Stat(MoviesExceptFile)
//
//	// create file if not exists
//	if os.IsNotExist(err) {
//		newJSON := &MoviesExcept{
//			Name: ""}
//		j, err := json.MarshalIndent(newJSON, "", " ")
//		if err != nil {
//			log.Println(err)
//		}
//		writeFile(j)
//	}
//}
