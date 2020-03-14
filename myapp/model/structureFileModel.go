package model

import (
	"time"
)

type Video struct {
	Movie []File  `json:"files"`
	Serie []Serie `json:"series"`
}

type Serie struct {
	Name        string
	Image       string   `json:"image"`
	Seasons     []Season `json:"seasons"`
	TrailerKey  string   `json:"trailer_key"`
	TrailerSite string   `json:"trailer_site"`
}

type Season struct {
	Name    string
	Files   []File `json:"files"`
	Trailer string `json:"trailer"`
}

type File struct {
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Date        time.Time `json:"date"`
	Taille      int64     `json:"taille"`
	TrailerKey  string    `json:"trailer_key"`
	TrailerSite string    `json:"trailer_site"`
}
