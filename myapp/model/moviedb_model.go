package model

import (
	"encoding/json"
	"fmt"
	"github.com/Machiel/slugify"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

func UnmarshalMovieDBModel(data []byte) (MovieDBModel, error) {
	var r MovieDBModel
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MovieDBModel) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MovieDBModel struct {
	Page         int64    `json:"page"`
	TotalResults int64    `json:"total_results"`
	TotalPages   int64    `json:"total_pages"`
	Results      []Result `json:"results"`
}

type Result struct {
	Popularity       float64 `json:"popularity"`
	ID               int64   `json:"id"`
	Video            bool    `json:"video"`
	VoteCount        int64   `json:"vote_count"`
	VoteAverage      float64 `json:"vote_average"`
	Title            string  `json:"title"`
	ReleaseDate      string  `json:"release_date"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIDS         []int64 `json:"genre_ids"`
	BackdropPath     string  `json:"backdrop_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
}

const apiMovieDB = "ea8779638f078f25daa3913e80fe46eb"

func GetImage(movie string, serie bool) string {
	var url string

	movie, year := splitVideos(movie, serie)
	if serie {
		url = fmt.Sprintf("https://api.themoviedb.org/3/search/tv?api_key=%s&query=%s&page=1&include_adult=true", apiMovieDB, movie)
		if year != "" {
			url += "&first_air_date_year=" + year
		}
	} else {
		url = fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s&page=1&include_adult=true", apiMovieDB, movie)
		if year != "" {
			url += "&year=" + year
		}
	}
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Pas possible d'accéder à https://api.themoviedb.org/")
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	var moviedb MovieDBModel
	moviedb, err = UnmarshalMovieDBModel(bodyBytes)
	if err != nil {
		log.Println(err)
	}
	if len(moviedb.Results) > 0 {
		if moviedb.Results[0].PosterPath != "" {
			return "https://image.tmdb.org/t/p/w500" + moviedb.Results[0].PosterPath
		}
	}
	return ""
}

func splitVideos(movie string, serie bool) (string, string) {
	ext := filepath.Ext(movie)
	if ext != "" && !serie {
		movie = movie[0 : len(movie)-len(ext)]
	}
	movie = slugify.Slugify(movie)
	movie = strings.Replace(movie, "-", "+", -1)
	return splitVideosyear(movie)
}

func splitVideosyear(movie string) (string, string) {
	var re = regexp.MustCompile(`(?m)\d{4}`)
	if re.MatchString(movie) {
		year := re.FindString(movie)
		// -1 pour retirer le caractère avant (+ | -)
		movie = movie[0 : len(movie)-len(year)-1]
		return movie, year
	}
	return movie, ""
}
