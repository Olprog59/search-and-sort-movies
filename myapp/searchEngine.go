package myapp

import (
	"fmt"
	"github.com/Machiel/slugify"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"net/http"
	"search-and-sort-movies/myapp/logger"
	"strings"
	"time"
)

func loopGetSearchEngine(name string) string {
	slug := slugify.New(slugify.Configuration{
		ReplaceCharacter: '+',
	})
	name = slug.Slugify(name)
	var proposition = getSearchEngine(name)

	for proposition != "" {
		name = slug.Slugify(proposition)
		proposition = getSearchEngine(name)
	}
	return name
}

// Ok Test
func getSearchEngine(name string) string {
	var proposition string

	req, err := http.NewRequest("GET", "https://www.ecosia.org/search", nil)

	if req != nil {
		param := req.URL.Query()
		param.Add("q", name)
		req.URL.RawQuery = param.Encode()
	}

	client := new(http.Client)

	resp, err := client.Do(req)

	if resp == nil || err != nil {
		time.Sleep(1 * time.Minute)
		return getSearchEngine(name)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("response status code was %d\n", resp.StatusCode)
		return name
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		fmt.Printf("response content type was %s not text/html\n", ctype)
		return name
	}

	defer resp.Body.Close()
	doc := html.NewTokenizer(resp.Body)
	//out, _ := os.Create("./bing.txt")
	//defer out.Close()
	//_, _ = io.Copy(out, resp.Body)

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
			fmt.Println(logger.Fata("error tokenizing HTML: %v", doc.Err()))
		}
		if tokenType == html.StartTagToken {
			//get the token
			token := doc.Token()
			//if the name of the element is "div"
			if atom.A == token.DataAtom {
				for _, v := range token.Attr {
					if v.Key == "class" && v.Val == "result-title" {
						tokenType = doc.Next()
						if tokenType == html.TextToken {
							//tokenType = doc.Next()
							token = doc.Token()
							if html.TextToken == tokenType {
								proposition = fmt.Sprintf("%v", token)
							}
						}
						token = doc.Token()
						break
					}
				}

			}
			if proposition != "" {
				break
			}
		}
	}
	return proposition
}
