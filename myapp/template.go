package myapp

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
)

type page struct {
	Title        string
	Navbar       string
	List         []string
	Config       *Config
	FlashMessage string
}

var store = sessions.NewCookieStore([]byte("samsam"))

func StartServerWeb() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods(http.MethodGet)
	r.HandleFunc("/except", exceptFile).Methods(http.MethodGet)
	r.HandleFunc("/config", configApp).Methods(http.MethodGet)
	r.HandleFunc("/config", configAppPost).Methods(http.MethodPost)
	go http.ListenAndServe(":1111", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html") //parse the html file homepage.html
	if err != nil {                                       // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, page{Title: "A trier", Navbar: "index", List: ReadAllFiles()}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                                                                   // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
func exceptFile(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/except.html") //parse the html file homepage.html
	if err != nil {                                        // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, page{Title: "Exception", Navbar: "except"}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                                                // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
func configApp(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/config.html") //parse the html file homepage.html
	if err != nil {                                        // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	session, err := store.Get(r, "flash-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	mess := session.Flashes("message")
	var message string
	if len(mess) > 0 {
		message = fmt.Sprint(mess[0])
	}
	session.Save(r, w)

	err = t.Execute(w, page{
		Title:  "Configuration",
		Navbar: "config",
		Config: &Config{
			Dlna:   GetEnv("dlna"),
			Movies: GetEnv("movies"),
			Series: GetEnv("series")},
		FlashMessage: message}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
func configAppPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dlna := r.Form.Get("dlna")
	movies := r.Form.Get("movies")
	series := r.Form.Get("series")

	SetEnv("dlna", dlna)
	SetEnv("movies", movies)
	SetEnv("series", series)

	session, err := store.Get(r, "flash-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.AddFlash("Les modifications ont bien été prises en compte", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/config", 301)
}
