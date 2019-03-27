package myapp

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

//type page struct {
//	Config    *Config
//	Exception []MoviesExcept
//	Pwd       string
//}

func RestStartServerWeb(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/movies", restAllMovies).Methods(http.MethodGet)
	//r.HandleFunc("/api/movies/{path}", restRemoveMovies).Methods(http.MethodGet)
	r.HandleFunc("/api/series", restAllSeries).Methods(http.MethodGet)
	//r.HandleFunc("/api/series/{path}", restRemoveSeries).Methods(http.MethodGet)
	r.HandleFunc("/api/log", restLogFile).Methods(http.MethodGet)
	r.HandleFunc("/api/except", restExceptFile).Methods(http.MethodGet)
	r.HandleFunc("/api/except", restExceptFilePost).Methods(http.MethodPost)
	//r.HandleFunc("/api/except/delete", restExceptFileDelete).Methods(http.MethodPost)
	r.HandleFunc("/api/config", restConfigApp).Methods(http.MethodGet)
	r.HandleFunc("/api/config", restConfigAppPost).Methods(http.MethodPost)
	//r.HandleFunc("/api/refresh", restRefresh).Methods(http.MethodPut)
	//r.HandleFunc("/mail", configAppMailPost).Methods(http.MethodPost)
	return r
}

//func restRemoveSeries(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	path := GetEnv("series") + string(os.PathSeparator) + strings.Replace(vars["path"], "_", "/", -1)
//	log.Println(path)
//	for _, serie := range ReadAllSeries().Serie {
//		for _, season := range serie.Seasons {
//			for _, file := range season.File {
//				if file.Path == path {
//					err := os.RemoveAll(file.Path)
//					if err != nil {
//						log.Println(err)
//					} else {
//						SaveAllSeries()
//					}
//				}
//			}
//		}
//	}
//	http.Redirect(w, r, "/series", 301)
//}
//func restRemoveMovies(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	path := GetEnv("movies") + string(os.PathSeparator) + strings.Replace(vars["path"], "_", "/", -1)
//	log.Println(path)
//	for _, files := range ReadAllMovies().Files {
//		if files.Path == path {
//			err := os.RemoveAll(files.Path)
//			if err != nil {
//				log.Println(err)
//			} else {
//				SaveAllMovies()
//			}
//		}
//
//	}
//	http.Redirect(w, r, "/movies", 301)
//}

//func restRefresh(w http.ResponseWriter, r *http.Request) {
//	serie := make(chan bool)
//	movie := make(chan bool)
//	go func() {
//		serie <- SaveAllSeries()
//		movie <- SaveAllMovies()
//	}()
//
//	if <-serie && <-movie {
//		w.Write([]byte("true"))
//	}
//}

//func restIndex(w http.ResponseWriter, r *http.Request) {
//	var err error
//	//t, err := template.ParseFiles("templates/index.html") //parse the html file homepage.html
//	//if err != nil { // if there is an error
//	//	log.Print("template parsing error: ", err) // log it
//	//}
//	t := template.New("index")
//	t.Parse(header + pageIndex + pageFooter)
//
//	err = t.Execute(w, page{Title: "A trier", Navbar: "index", List: ReadAllFiles()}) //execute the template and pass it the HomePageVars struct to fill in the gaps
//	if err != nil {                                                                   // if there is an error
//		log.Print("template executing error: ", err) //log it
//	}
//}

func originProblem(w http.ResponseWriter, r *http.Request) {
	OriginList := r.Header["Origin"]
	Origin := ""
	if len(OriginList) > 0 {
		Origin = OriginList[0]
	}
	w.Header().Add("Access-Control-Allow-Origin", Origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

func restAllMovies(w http.ResponseWriter, r *http.Request) {
	originProblem(w, r)
	json.NewEncoder(w).Encode(ReadAllMovies())
}

func restAllSeries(w http.ResponseWriter, r *http.Request) {
	originProblem(w, r)
	json.NewEncoder(w).Encode(ReadAllSeries())
}

//func restIndexPost(w http.ResponseWriter, r *http.Request) {
//	name := r.FormValue("ajax_data")
//	oldName := r.FormValue("oldName")
//
//	dlna := GetEnv("dlna")
//	os.Rename(dlna+string(os.PathSeparator)+oldName, dlna+string(os.PathSeparator)+name)
//}

func reverse(ss []mapLog) []mapLog {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
	return ss
}

func restLogFile(w http.ResponseWriter, r *http.Request) {
	originProblem(w, r)
	json.NewEncoder(w).Encode(reverse(ReadFileLog()))
}

func restExceptFile(w http.ResponseWriter, r *http.Request) {
	originProblem(w, r)
	json.NewEncoder(w).Encode(GetAllExcept())
}

func restExceptFilePost(w http.ResponseWriter, r *http.Request) {
	except := r.FormValue("except")

	SetMoviesExceptFile(except)
	//session, err := store.Get(r, "flash-session")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
	//session.AddFlash("Les modifications ont bien été prises en compte", "message")
	//session.Save(r, w)

	http.Redirect(w, r, "/except", 301)
}

//func restExceptFileDelete(w http.ResponseWriter, r *http.Request) {
//	name := r.FormValue("ajax_data")
//	RemoveMoviesExceptFile(name)
//}

func restConfigApp(w http.ResponseWriter, r *http.Request) {
	originProblem(w, r)
	var p = page{
		Config: &Config{
			Dlna:   GetEnv("dlna"),
			Movies: GetEnv("movies"),
			Series: GetEnv("series"),
		},
		Pwd: pwd("", false),
	}
	json.NewEncoder(w).Encode(p)
}

func restConfigAppPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dlna := r.Form.Get("dlna")
	movies := r.Form.Get("movies")
	series := r.Form.Get("series")

	SetEnv("dlna", dlna)
	SetEnv("movies", movies)
	SetEnv("series", series)

	//session, err := store.Get(r, "flash-session")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
	//session.AddFlash("Les modifications ont bien été prises en compte", "message")
	//session.Save(r, w)

	http.Redirect(w, r, "/config", 301)
}
