package myapp

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/Machiel/slugify"
)

type movieDBTv struct {
	Page    int `json:"page"`
	Results []struct {
		BackdropPath     string   `json:"backdrop_path"`
		FirstAirDate     string   `json:"first_air_date"`
		GenreIds         []int    `json:"genre_ids"`
		ID               int      `json:"id"`
		Name             string   `json:"name"`
		OriginCountry    []string `json:"origin_country"`
		OriginalLanguage string   `json:"original_language"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity"`
		PosterPath       string   `json:"poster_path"`
		VoteAverage      int      `json:"vote_average"`
		VoteCount        int      `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type movieDBMovie struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		ID               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

const (
	apiV4 = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJlYTg3Nzk2MzhmMDc4ZjI1ZGFhMzkxM2U4MGZlNDZlYiIsInN1YiI6IjU5Y2Y3NjdiYzNhMzY4MWViMTAxOThjNyIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.mAxfQbzn4WIft74XAooGGiw7PhHxMNTu8TtTvPwhh1c"
	apiV3 = "ea8779638f078f25daa3913e80fe46eb"
)

func checkMovieDB(tv, lang bool, name string, date ...string) string {

	var language string
	var tvOrMovie = "movie"

	if lang {
		language = "&language=fr-FR"
	}

	if tv {
		tvOrMovie = "tv"
		name = slugRemoveYearSerieForSearchMovieDB(name)
	}

	var year string
	if len(date) > 0 {
		if tvOrMovie == "movie" {
			year = "&year=" + date[0]
		} else {
			year = "&first_air_date_year=" + date[0]
		}
	}

	url := "https://api.themoviedb.org/3/search/" + tvOrMovie + "?api_key=" + apiV3 + language + "&query=" + name + year

	return url

}

func slugRemoveYearSerieForSearchMovieDB(name string) (new string) {
	year := regexp.MustCompile(`^[0-9]{4}$`)
	for _, v := range strings.Split(name, "-") {
		if year.MatchString(v) {
			break
		}
		new += v + " "
	}
	return slugify.Slugify(new)
}

func dbSeries(lang bool, name, date string) (movieDBTv, error) {
	url := checkMovieDB(true, lang, name, date)
	return readJSONFromUrl_TV(url)
}

func dbMovies(lang bool, name, date string) (movieDBMovie, error) {
	url := checkMovieDB(false, lang, name, date)
	return readJSONFromUrl_Movie(url)
}

func readJSONFromUrl_TV(url string) (movieDBTv, error) {
	var movie movieDBTv

	resp, err := http.Get(url)
	if err != nil {
		return movie, err
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&movie)

	return movie, nil
}

func readJSONFromUrl_Movie(url string) (movieDBMovie, error) {
	var movie movieDBMovie

	resp, err := http.Get(url)
	if err != nil {
		return movie, err
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&movie)

	return movie, nil
}
