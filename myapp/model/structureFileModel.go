package model

import "time"

type AllFiles struct {
	Movie Movie `json:"movie"`
	Serie Serie `json:"serie"`
}

type Movie struct {
	Files []File `json:"files"`
}

type Serie struct {
	Folder []Folder `json:"folder"`
}

type Folder struct {
	Name    string
	Seasons []Season `json:"seasons"`
}

type Season struct {
	Name  string
	Files []File `json:"files"`
}

type File struct {
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Taille string    `json:"taille"`
}
