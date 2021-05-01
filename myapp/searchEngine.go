package myapp

import (
	"github.com/Machiel/slugify"
	"io/ioutil"
	"net/http"
	"regexp"
	"search-and-sort-movies/myapp/logger"
	"strings"
	"time"
)

func loopGetSearchEngine(name string) string {
	slug := slugify.New(slugify.Configuration{
		ReplaceCharacter: '+',
	})
	name = slug.Slugify(name)
	proposition, distance := getSearchEngine(name)

	if proposition != "" && distance < 7 {
		name = proposition
	}
	return name
}

// Ok Test
func getSearchEngine(name string) (string, int) {
	var proposition string

	req, err := http.NewRequest("GET", "https://www.google.com/search", nil)

	if req != nil {
		param := req.URL.Query()
		param.Add("q", name)
		param.Add("lr", "lang_en")
		param.Add("hl", "en")
		req.URL.RawQuery = param.Encode()
	}

	client := new(http.Client)
	resp, err := client.Do(req)

	if resp == nil || err != nil {
		time.Sleep(1 * time.Minute)
		return getSearchEngine(name)
	}

	if resp.StatusCode != http.StatusOK {
		logger.L(logger.Red, "response status code was %d", resp.StatusCode)
		return "", 0
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		logger.L(logger.Red, "response content type was %s not text/html", ctype)
		return "", 0
	}
	var re = regexp.MustCompile(`function\(\){var q='(?P<newName>[\w\s\d]+)';\(`)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.L(logger.Red, "%s", err)
	}
	matches := re.FindStringSubmatch(string(body))
	lastIndex := re.SubexpIndex("newName")
	defer resp.Body.Close()

	if len(matches) >= lastIndex {
		proposition = strings.ToLower(matches[lastIndex])
	}

	re = regexp.MustCompile(`(?i)>(?P<newName>[\w\d:\-_\s()]{1,60})(\s&#8212;|\s-)\swikip`)
	matches = re.FindStringSubmatch(string(body))
	lastIndex = re.SubexpIndex("newName")
	defer resp.Body.Close()

	var distance int

	if len(matches) >= lastIndex {
		prop := strings.ToLower(matches[lastIndex])
		distance = ComputeDistance(proposition, prop)
		proposition = strings.ToLower(prop)

		logger.L(logger.Green, "Distance Levenshtein : %d", distance)

	}

	tab := strings.Split(proposition, " (")
	if len(tab) > 1 {
		proposition = tab[0]
	}
	logger.L(logger.Green, "Proposition par l'algo de recherche google : "+proposition)

	return proposition, distance
}
