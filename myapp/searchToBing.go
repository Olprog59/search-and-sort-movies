package myapp

import (
	"fmt"
	"github.com/Machiel/slugify"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"log"
	"net/http"
	"strings"
)

func loopGetBingName(name string) string {
	slug := slugify.New(slugify.Configuration{
		ReplaceCharacter: '+',
	})
	name = slug.Slugify(name)
	var proposition = getBingName(name)

	for proposition != "" {
		name = slug.Slugify(proposition)
		proposition = getBingName(name)
	}
	return name
}

func getBingName(name string) string {
	var proposition string

	resp, _ := http.Get("https://www.bing.com/search?q=" + name)
	log.Println("https://www.bing.com/search?q=" + name)

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
							//log.Println(html.TextToken == tokenType, atom.Strong == token.DataAtom)
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
	return proposition
}
