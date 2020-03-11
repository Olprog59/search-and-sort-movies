package model

import "time"

type Video struct {
	Movie []File  `json:"files"`
	Serie []Serie `json:"series"`
}

type Serie struct {
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
	Taille int64     `json:"taille"`
}
