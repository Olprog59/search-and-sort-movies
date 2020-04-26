package myapp

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"log"
	"net/http"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"strings"

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

func getBingName(name string) string {

	var proposition = ""

	resp, _ := http.Get("https://www.bing.com/search?q=" + name)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code was %d\n", resp.StatusCode)
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}

	defer resp.Body.Close()
	doc := html.NewTokenizer(resp.Body)

	for {
		tokenType := doc.Next()

		//if it's an error token, we either reached
		//the end of the file, or the HTML was malformed
		if tokenType == html.ErrorToken {
			err := doc.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			log.Fatalf("error tokenizing HTML: %v", doc.Err())
		}
		if tokenType == html.StartTagToken {
			//get the token
			token := doc.Token()
			//if the name of the element is "div"
			if atom.Div == token.DataAtom {
				for _, v := range token.Attr {
					if v.Key == "id" && v.Val == "sp_requery" {
						tokenType = doc.Next()
						if tokenType == html.TextToken {
							tokenType = doc.Next()
						}
						token = doc.Token()
						if atom.A == token.DataAtom {
							tokenType = doc.Next()
							token = doc.Token()
							log.Println(html.TextToken == tokenType, atom.Strong == token.DataAtom)
							if html.TextToken == tokenType {
								proposition = fmt.Sprintf("%v", token)
							}

							if atom.Strong == token.DataAtom {
								tokenType = doc.Next()
								token = doc.Token()
								proposition = fmt.Sprintf("%v", token)
							}

							tokenType = doc.Next()
							token = doc.Token()
							if atom.Strong == token.DataAtom {
								tokenType = doc.Next()
								token = doc.Token()
								if tokenType == html.TextToken {
									proposition += fmt.Sprintf("%v", token)
								}
							}
							break
						}
						break
					}
				}

			}
			if proposition != "" {
				break
			}
		}
	}
	log.Println(proposition)
	return proposition
}

func checkMovieDB(tv, lang bool, name, originalName string) string {

	var language string

	if lang {
		language = "&language=fr-FR"
	}

	if tv {
		name = slugRemoveYearSerieForSearchMovieDB(name)
	}

	slugify.New(slugify.Configuration{
		ReplaceCharacter: '+',
	})
	slugify.Slugify(originalName)
	slugify.Slugify(name)
	if newName := getBingName(originalName); newName != "" {
		originalName = slugify.Slugify(newName)
	}
	if newName := getBingName(name); newName != "" {
		name = slugify.Slugify(newName)
	}
	var url string
	if len(originalName) > 0 {
		url = "https://api.themoviedb.org/3/search/multi?api_key=" + constants.ApiV3 + language + "&query=" + originalName
	} else {
		url = "https://api.themoviedb.org/3/search/multi?api_key=" + constants.ApiV3 + language + "&query=" + name
	}
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

func dbSeries(lang bool, name, original string) (MoviesDb, error) {
	url := checkMovieDB(true, lang, name, original)
	return readJSONFromUrlTV(url)
}

func dbMovies(lang bool, name, original string) (MoviesDb, error) {
	url := checkMovieDB(false, lang, name, original)
	return readJSONFromUrlMovie(url)
}

func readJSONFromUrlTV(url string) (MoviesDb, error) {
	var movie MoviesDb

	resp, err := http.Get(url)
	if err != nil {
		return movie, err
	}

	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&movie)

	return movie, nil
}

func readJSONFromUrlMovie(url string) (MoviesDb, error) {
	var movie MoviesDb

	resp, err := http.Get(url)
	if err != nil {
		return movie, err
	}

	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&movie)

	return movie, nil
}
