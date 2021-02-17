package model

import (
	"encoding/json"
	"fmt"
	"github.com/Machiel/slugify"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
	"strings"
	"time"
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

func UnmarshalTrailer(data []byte) (Trailer, error) {
	var r Trailer
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Trailer) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Trailer struct {
	ID      int64           `json:"id"`
	Results []ResultTrailer `json:"results"`
}

type ResultTrailer struct {
	ID       string `json:"id"`
	Iso6391  string `json:"iso_639_1"`
	Iso31661 string `json:"iso_3166_1"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Site     string `json:"site"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

// Ok Test
func GetTrailer(id int64, serie bool) (string, string) {
	if id == 0 {
		return "", ""
	}
	var videos string
	if serie {
		videos = "tv"
	} else {
		videos = "movie"
	}
	url := fmt.Sprintf("https://api.themoviedb.org/3/%s/%d/videos?api_key=%s&language=en-US", videos, id, constants.ApiV3)
	resp, err := http.Get(url)
	var bodyBytes []byte
	if err != nil || resp == nil {
		logger.L(logger.Green, "Pas possible d'accéder à https://api.themoviedb.org/")
		time.Sleep(1 * time.Minute)
		return GetTrailer(id, serie)
	} else {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	defer resp.Body.Close()

	var trailer Trailer
	trailer, err = UnmarshalTrailer(bodyBytes)
	if err != nil {
		logger.L(logger.Red, "", err)
	}
	if len(trailer.Results) > 0 {
		return trailer.Results[0].Key, trailer.Results[0].Site
	}
	return "", ""
}

// Ok Test
func GetImage(movie string, serie bool) (string, int64) {
	var url string

	movie, year := splitVideos(movie, serie)
	if serie {
		url = fmt.Sprintf("https://api.themoviedb.org/3/search/tv?api_key=%s&query=%s&page=1&include_adult=true", constants.ApiV3, movie)
		if year != "" {
			url += "&first_air_date_year=" + year
		}
	} else {
		url = fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s&page=1&include_adult=true", constants.ApiV3, movie)
		if year != "" {
			url += "&year=" + year
		}
	}

	resp, err := http.Get(url)
	var bodyBytes []byte
	if err != nil || resp == nil {
		logger.L(logger.Green, "Pas possible d'accéder à https://api.themoviedb.org/")
		time.Sleep(1 * time.Minute)
		return GetImage(movie, serie)
	} else {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	defer resp.Body.Close()

	var moviedb MovieDBModel
	moviedb, err = UnmarshalMovieDBModel(bodyBytes)
	if err != nil {
		logger.L(logger.Red, "", err)
	}
	if len(moviedb.Results) > 0 {
		if moviedb.Results[0].PosterPath != "" {
			return "https://image.tmdb.org/t/p/w500" + moviedb.Results[0].PosterPath, moviedb.Results[0].ID
		}
	}
	return "", 0
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
