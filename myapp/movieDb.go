package myapp

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"strings"
	"time"

	"github.com/Machiel/slugify"
)

//func UnmarshalMovies(data []byte) (MoviesDb, error) {
//	var r MoviesDb
//	err := json.Unmarshal(data, &r)
//	return r, err
//}

func (r *MoviesDb) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MoviesDb struct {
	Page         *int64   `json:"page,omitempty"`
	Results      []Result `json:"results"`
	TotalResults *int64   `json:"total_results,omitempty"`
	TotalPages   *int64   `json:"total_pages,omitempty"`
}

type Result struct {
	PosterPath       *string           `json:"poster_path"`
	Popularity       *float64          `json:"popularity,omitempty"`
	ID               *int64            `json:"id,omitempty"`
	Overview         *string           `json:"overview,omitempty"`
	BackdropPath     *string           `json:"backdrop_path"`
	VoteAverage      *float64          `json:"vote_average,omitempty"`
	MediaType        *MediaType        `json:"media_type,omitempty"`
	FirstAirDate     *string           `json:"first_air_date,omitempty"`
	OriginCountry    []string          `json:"origin_country"`
	GenreIDS         []int64           `json:"genre_ids"`
	OriginalLanguage *OriginalLanguage `json:"original_language,omitempty"`
	VoteCount        *int64            `json:"vote_count,omitempty"`
	Name             *string           `json:"name,omitempty"`
	OriginalName     *string           `json:"original_name,omitempty"`
	Adult            *bool             `json:"adult,omitempty"`
	ReleaseDate      *string           `json:"release_date,omitempty"`
	OriginalTitle    *string           `json:"original_title,omitempty"`
	Title            *string           `json:"title,omitempty"`
	Video            *bool             `json:"video,omitempty"`
	ProfilePath      *string           `json:"profile_path"`
	KnownFor         []Result          `json:"known_for"`
}

type MediaType string

//const (
//	Movie  MediaType = "movie"
//	Person MediaType = "person"
//	Tv     MediaType = "tv"
//)

type OriginalLanguage string

//const (
//	En OriginalLanguage = "en"
//	It OriginalLanguage = "it"
//)

func (m *myFile) checkMovieDB(tv, lang bool, name string) string {

	var language string

	if lang {
		language = "&language=fr-FR"
	}

	if tv {
		name = slugRemoveYearSerieForSearchMovieDB(name)
	}

	slug := slugify.New(slugify.Configuration{
		ReplaceCharacter: '+',
	})

	var url string

	name = slug.Slugify(name)
	if newName := loopGetBingName(name); newName != "" {
		name = slug.Slugify(newName)
	}
	url = "https://api.themoviedb.org/3/search/multi?api_key=" + constants.ApiV3 + language + "&query=" + name
	//}
	m.bingName = name
	log.Println(url)
	return url
}

/**
@Deprecated
*/
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

func (m *myFile) dbSeries(lang bool, name string) (MoviesDb, error) {
	url := m.checkMovieDB(true, lang, name)
	return readJSONFromUrlTV(url)
}

func (m *myFile) dbMovies(lang bool, name string) (MoviesDb, error) {
	url := m.checkMovieDB(false, lang, name)
	return readJSONFromUrlMovie(url)
}

// Ok Test
// ex : https://api.themoviedb.org/3/search/multi?api_key=ea8779638f078f25daa3913e80fe46eb&query=naruto
func readJSONFromUrlTV(url string) (MoviesDb, error) {
	var movie MoviesDb

	resp, err := http.Get(url)

	if resp == nil || err != nil {
		time.Sleep(1 * time.Minute)
		return readJSONFromUrlTV(url)
	}

	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&movie)

	return movie, nil
}

// Ok Test
func readJSONFromUrlMovie(url string) (MoviesDb, error) {
	var movie MoviesDb

	resp, err := http.Get(url)

	if resp == nil || err != nil {
		time.Sleep(1 * time.Minute)
		return readJSONFromUrlMovie(url)
	}

	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&movie)

	return movie, nil
}
