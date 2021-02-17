package myapp

import (
	"github.com/Machiel/slugify"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func loopGetSearchEngine(name string) string {
	slug := slugify.New(slugify.Configuration{
		ReplaceCharacter: '+',
	})
	name = slug.Slugify(name)
	var proposition = getSearchEngine(name)

	if proposition != "" {
		name = proposition
	}
	return name
}

// Ok Test
func getSearchEngine(name string) string {
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
	log.Println(req)
	resp, err := client.Do(req)

	if resp == nil || err != nil {
		time.Sleep(1 * time.Minute)
		return getSearchEngine(name)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		log.Printf("response status code was %d\n", resp.StatusCode)
		return ""
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Printf("response content type was %s not text/html\n", ctype)
		return ""
	}
	var re = regexp.MustCompile(`function\(\){var q='(?P<newName>[\w\s\d]+)';\(`)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
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

	if len(matches) >= lastIndex {
		prop := strings.ToLower(matches[lastIndex])

		if len(proposition) > 0 {
			distance := ComputeDistance(proposition, prop)
			log.Printf("Distance Levenshtein : %d", distance)
			if distance < 10 {
				proposition = strings.ToLower(prop)
			}
		} else {
			proposition = strings.ToLower(prop)
		}
	}

	tab := strings.Split(proposition, " (tv series)")
	if len(tab) > 1 {
		proposition = tab[0]
	}
	log.Println(proposition)

	//doc := html.NewTokenizer(resp.Body)
	//out, _ := os.Create("./bing.txt")
	//defer out.Close()
	//_, _ = io.Copy(out, resp.Body)

	//for {
	//	tokenType := doc.Next()
	//
	//	//if it's an error token, we either reached
	//	//the end of the file, or the HTML was malformed
	//	if tokenType == html.ErrorToken {
	//		err := doc.Err()
	//		if err == io.EOF {
	//			//end of the file, break out of the loop
	//			break
	//		}
	//		log.Println(logger.Fata("error tokenizing HTML: %v", doc.Err()))
	//	}
	//	if tokenType == html.StartTagToken {
	//		//get the token
	//		token := doc.Token()
	//		//if the name of the element is "div"
	//		if atom.A == token.DataAtom {
	//			for _, v := range token.Attr {
	//				if v.Key == "class" && v.Val == "result-title" {
	//					tokenType = doc.Next()
	//					if tokenType == html.TextToken {
	//						//tokenType = doc.Next()
	//						token = doc.Token()
	//						if html.TextToken == tokenType {
	//							proposition = fmt.Sprintf("%v", token)
	//						}
	//					}
	//					token = doc.Token()
	//					break
	//				}
	//			}
	//
	//		}
	//		if proposition != "" {
	//			break
	//		}
	//	}
	//}
	return proposition
}
