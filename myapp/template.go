package myapp

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type page struct {
	Title        string
	Navbar       string
	List         []string
	Config       *Config
	FlashMessage string
	Exception    []MoviesExcept
	Pwd          string
	Log          []string
	Movie        Movie
	Serie        Serie
}

var store = sessions.NewCookieStore([]byte("samsam"))

func StartServerWeb() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods(http.MethodGet)
	r.HandleFunc("/", indexPost).Methods(http.MethodPost)
	r.HandleFunc("/movies", allMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies/{path}", removeMovies).Methods(http.MethodGet)
	r.HandleFunc("/series", allSeries).Methods(http.MethodGet)
	r.HandleFunc("/series/{path}", removeSeries).Methods(http.MethodGet)
	r.HandleFunc("/log", logFile).Methods(http.MethodGet)
	r.HandleFunc("/except", exceptFile).Methods(http.MethodGet)
	r.HandleFunc("/except", exceptFilePost).Methods(http.MethodPost)
	r.HandleFunc("/except/delete", exceptFileDelete).Methods(http.MethodPost)
	r.HandleFunc("/config", configApp).Methods(http.MethodGet)
	r.HandleFunc("/config", configAppPost).Methods(http.MethodPost)
	r.HandleFunc("/mail", configAppMailPost).Methods(http.MethodPost)
	go http.ListenAndServe(":1515", r)
}
func removeSeries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := GetEnv("series") + string(os.PathSeparator) + strings.Replace(vars["path"], "_", "/", -1)
	log.Println(path)
	for _, serie := range ReadAllSeries().Series {
		for _, season := range serie.Seasons {
			for _, file := range season.Files {
				if file.Path == path {
					err := os.RemoveAll(file.Path)
					if err != nil {
						log.Println(err)
					} else {
						SaveAllSeries()
					}
				}
			}
		}
	}
	http.Redirect(w, r, "/series", 301)
}
func removeMovies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := GetEnv("movies") + string(os.PathSeparator) + strings.Replace(vars["path"], "_", "/", -1)
	log.Println(path)
	for _, files := range ReadAllMovies().Files {
		if files.Path == path {
			err := os.RemoveAll(files.Path)
			if err != nil {
				log.Println(err)
			} else {
				SaveAllMovies()
			}
		}

	}
	http.Redirect(w, r, "/movies", 301)
}

func index(w http.ResponseWriter, r *http.Request) {
	var err error
	//t, err := template.ParseFiles("templates/index.html") //parse the html file homepage.html
	//if err != nil { // if there is an error
	//	log.Print("template parsing error: ", err) // log it
	//}
	t := template.New("index")
	t.Parse(header + pageIndex + pageFooter)

	err = t.Execute(w, page{Title: "A trier", Navbar: "index", List: ReadAllFiles()}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                                                                   // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func allMovies(w http.ResponseWriter, r *http.Request) {
	var err error
	//t, err := template.New("movies.html").Funcs(template.FuncMap{
	t := template.New("movies.html").Funcs(template.FuncMap{
		"title": func(str string) string {
			return strings.Title(strings.Join(strings.Split(str, "-"), " "))
		},
		"rename": func(str string) string {
			str = str[len(GetEnv("movies"))+1:]
			return strings.Replace(str, "/", "_", -1)
		},
	})
	//}).ParseFiles("templates/movies.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	//t := template.New("index")
	t.Parse(header + pageMovies + pageFooter)

	err = t.Execute(w, page{Title: "Movies", Navbar: "movies", Movie: ReadAllMovies()}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                                                                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	SaveAllMovies()
}

func allSeries(w http.ResponseWriter, r *http.Request) {
	var err error
	//t, err := template.New("series.html").Funcs(template.FuncMap{
	t := template.New("series.html").Funcs(template.FuncMap{
		"title": func(str string) string {
			return strings.Title(strings.Join(strings.Split(str, "-"), " "))
		},
		"rename": func(str string) string {
			str = str[len(GetEnv("series"))+1:]
			return strings.Replace(str, "/", "_", -1)
		},
	})
	//}).ParseFiles("templates/series.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	//t := template.New("index")
	t.Parse(header + pageSeries + pageFooter)

	err = t.Execute(w, page{Title: "Series", Navbar: "series", Serie: ReadAllSeries()}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                                                                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	SaveAllSeries()
}

func indexPost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("ajax_data")
	oldName := r.FormValue("oldName")

	dlna := GetEnv("dlna")
	os.Rename(dlna+string(os.PathSeparator)+oldName, dlna+string(os.PathSeparator)+name)
}

func logFile(w http.ResponseWriter, r *http.Request) {
	t := template.New("logFile")
	t.Parse(header + pageLog + pageFooter)

	err := t.Execute(w, page{Title: "Log", Navbar: "log", Log: ReadFileLog()}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                                                            // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func exceptFile(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("templates/except.html") //parse the html file homepage.html
	//if err != nil { // if there is an error
	//	log.Print("template parsing error: ", err) // log it
	//}
	t := template.New("exceptFile")
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

	t.Parse(header + pageExcept + pageFooter)

	err = t.Execute(w, page{Title: "Exception", Navbar: "except", Exception: readFile(), FlashMessage: message}) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                                                                                              // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func exceptFilePost(w http.ResponseWriter, r *http.Request) {
	except := r.FormValue("except")

	SetMoviesExceptFile(except)
	session, err := store.Get(r, "flash-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.AddFlash("Les modifications ont bien été prises en compte", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/except", 301)
}

func exceptFileDelete(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("ajax_data")
	RemoveMoviesExceptFile(name)
}

func configApp(w http.ResponseWriter, r *http.Request) {
	t := template.New("configApp")
	//t, err := template.ParseFiles("templates/config.html")
	//if err != nil { // if there is an error
	//	log.Print("template parsing error: ", err) // log it
	//}
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

	t.Parse(header + pageConfig + pageFooter)

	err = t.Execute(w, page{
		Title:  "Configuration",
		Navbar: "config",
		Config: &Config{
			Dlna:   GetEnv("dlna"),
			Movies: GetEnv("movies"),
			Series: GetEnv("series"),
			//Mail:   GetEnv("mail")
		},
		FlashMessage: message,
		Pwd:          pwd("", false)}) //execute the template and pass it the HomePageVars struct to fill in the gaps
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

func configAppMailPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mail := r.Form.Get("mail")

	SetEnv("mail", mail)

	session, err := store.Get(r, "flash-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.AddFlash("Les modifications ont bien été prises en compte", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/config", 301)
}

const header = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css"
          integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
	<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.3.1/css/all.css" integrity="sha384-mzrmE5qonljUremFsqc01SB46JvROS7bZs3IO2EmfFsd15uHvIt+Y8vEf7N7fWAU" crossorigin="anonymous">

    <style>
        .input-group-append{
            cursor: pointer;
        }
		.container-fluid{
			padding-top: 70px;
		}
		.page-log{
			height: 90vh;
			min-height: 500px;
		}
		.page-log > #scroll {
			overflow: scroll;
			background-color: #343A40;
			color: #fff;
		}
    </style>
</head>
<body>
<nav class="navbar navbar-expand-lg navbar-dark bg-dark fixed-top">
    <a class="navbar-brand" href="/">Search and sort Movies</a>
    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent"
            aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
    </button>

    <div class="collapse navbar-collapse" id="navbarSupportedContent">
        <ul class="navbar-nav mr-auto">
            <li class="nav-item {{ if eq .Navbar "index"}}active{{end}}">
                <a class="nav-link" href="/">Liste à trier<span class="sr-only">(current)</span></a>
            </li>
 			<li class="nav-item {{ if eq .Navbar "movies"}}active{{end}}">
                <a class="nav-link" href="/movies">Movies</a>
            </li>
            <li class="nav-item {{ if eq .Navbar "series"}}active{{end}}">
                <a class="nav-link" href="/series">Series</a>
            </li>
            <li class="nav-item {{ if eq .Navbar "except"}}active{{end}}">
                <a class="nav-link" href="/except">Exception</a>
            </li>
            <li class="nav-item {{ if eq .Navbar "config"}}active{{end}}">
                <a class="nav-link" href="/config">Configuration</a>
            </li>
            <li class="nav-item {{ if eq .Navbar "log"}}active{{end}}">
                <a class="nav-link" href="/log">Log</a>
            </li>
        </ul>
    </div>
</nav>

<div class="container-fluid">
{{ $message := .FlashMessage }}
{{ if ne $message "" }}
    <div class="alert alert-success alert-dismissible fade show" role="alert">
    {{ $message }}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
{{ end }}
`

const pageIndex = `
    <ul class="list-group list-group-flush">
    {{ range .List }}
        <li class="list-group-item list-group-item-action">
            <div class="input-group">
                <div class="input-group-prepend">
                    <span class="input-group-text text" id="inputGroup-sizing-default">{{ . }}</span>
                </div>
                <input type="text"
                       placeholder="Entrer ici le nouveau nom de fichier et valider avec le bouton de droite"
                       class="form-control newName"
                       aria-label="Default"
                       aria-describedby="inputGroup-sizing-default">
                    <div class="input-group-append">
                        <span class="input-group-text">&#10004;</span>
                    </div>
            </div>
        </li>
    {{ end }}
    </ul>
`

const pageMovies = `
	<div class="row">
    {{ range $v := .Movie.Files }}
        <div class="col-sm-6 col-lg-4 col-xl-3 mb-3 files">
            <div class="card text-white bg-dark mb-3">
                <div class="card-body">
                
					<div class="row">
						<div class="col-10">{{ $v.Name | title }}</div> 
						<div class="col-2">
							<a href="/movies/{{ $v.Path | rename }}" class="delete">
								<i class="fas fa-trash-alt"></i>
							</a>
						</div>
					</div>
                </div>
            </div>
        </div>
    {{ end }}
    </div>
`

const pageSeries = `
 	<div class="row">
    {{ range $s := .Serie.Series }}
        <div class="col-sm-6 col-lg-4 col-xl-3 mb-3">
            <div class="card">
                <div class="card-body">
                    <h4 class="card-title text-center">
                    {{ $s.Name | title }}
                    </h4>
                    <div class="card-text">
                        <div class="accordion" id="parent-{{ $s.Name }}">
                        {{ range $season := $s.Seasons }}
                            <div class="card text-white bg-dark">
                                <div class="card-header" id="heading-{{ $s.Name }}">
                                    <h5 class="mb-0">
                                        <button class="btn btn-scondary col-sm-12" type="button" data-toggle="collapse"
                                                data-target="#{{ $s.Name }}-{{ $season.Name }}" aria-expanded="false"
                                                aria-controls="{{ $s.Name }}-{{ $season.Name }}">
                                        {{ $season.Name }}
                                        </button>
                                    </h5>
                                </div>

                                <div id="{{ $s.Name }}-{{ $season.Name }}" class="collapse"
                                     aria-labelledby="heading-{{ $s.Name }}"
                                     data-parent="#parent-{{ $s.Name }}">
                                    <div class="card-body">
                                        <ul class="list-group">
                                    {{ range $file := $season.Files }}
                                            <li class="list-group-item list-group-item-light files">
												<div class="row">
													<div class="col-10">{{ $file.Name }}</div> 
													<div class="col-2">
														<a href="/series/{{ $file.Path | rename }}" class="delete">
															<i class="fas fa-trash-alt"></i>
														</a>
													</div>
												</div>
											</li>
                                    {{ end }}
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        {{ end }}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {{ end }}
    </div>
`

const pageLog = `
<div class="row justify-content-center page-log">
<div class="col-10" id="scroll">
{{ range .Log }}
        <samp>{{ . }}</samp><br>
    {{ end }}
</div>
</div>
`

const pageExcept = `
 <form action="/except" method="post">
        <div class="form-group">
            <label for="except"></label>
            <input type="text" name="except" class="form-control" id="except" aria-describedby="exceptHelp"
                   placeholder="Entrer une exception de titre de film ou série" required>
            <small id="dlnaHelp" class="form-text text-muted">Entrer une exception si le film ou la série ne se trie pas
                bien ex: the-100
            </small>
        </div>
        <button type="submit" class="btn btn-primary">Sauvegarder</button>
    </form>
    <hr>
    <ul class="list-group">
    {{ range .Exception }}
        <li class="list-group-item">
        {{ .Name }}
            <button type="button" class="close removeExcept" data-name="{{ .Name }}" aria-label="Close">
                <span aria-hidden="true">&times;</span>
            </button>
        </li>
    {{ end }}
    </ul>
`

const pageConfig = `
<div class="alert alert-info alert-dismissible fade show" role="alert">
        {{ .Pwd }}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
 <form action="/config" method="post">
        <div class="form-group">
            <label for="dlna">Dossier à Trier</label>
            <input type="text" name="dlna" class="form-control" id="dlna" aria-describedby="dlnaHelp"
                   placeholder="Entrer le chemin où seront trier les films et séries" required value="{{.Config.Dlna}}">
            <small id="dlnaHelp" class="form-text text-muted">Entrer un chemin de type : C:\users\dlna ou
                /home/user/dlna
            </small>
        </div>
        <div class="form-group">
            <label for="movies">Dossier de films</label>
            <input type="text" name="movies" class="form-control" id="movies" aria-describedby="moviesHelp"
                   placeholder="Entrer le chemin où seront stocker les films" required value="{{.Config.Movies}}">
            <small id="moviesHelp" class="form-text text-muted">Entrer un chemin de type : C:\users\movies ou
                /home/user/movies
            </small>
        </div>
        <div class="form-group">
            <label for="series">Dossier de séries</label>
            <input type="text" name="series" class="form-control" id="series" aria-describedby="seriesHelp"
                   placeholder="Entrer le chemin où seront stocker les séries" required value="{{.Config.Series}}">
            <small id="seriesHelp" class="form-text text-muted">Entrer un chemin de type : C:\users\series ou
                /home/user/series
            </small>
        </div>
        <button type="submit" class="btn btn-primary">Sauvegarder</button>
    </form>
 	<!--<form action="/mail" method="post">
        <div class="form-group">
            <label for="mail">Email afin d'envoyer si un problème survient'</label>
            <input type="email" name="mail" class="form-control" id="mail" aria-describedby="mailHelp"
                   placeholder="Entrer le chemin où seront trier les films et séries" required value="{{.Config.Mail}}">
            <small id="mailHelp" class="form-text text-muted">Entrer un email</small>
        </div>
        <button type="submit" class="btn btn-primary">Sauvegarder</button>
    </form>-->
`

const pageFooter = `
</div>
<script src="https://code.jquery.com/jquery-3.3.1.min.js"
        integrity="sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"
        integrity="sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49"
        crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js"
        integrity="sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy"
        crossorigin="anonymous"></script>

<script>
    $(function () {
		$(function() {
  			var wtf    = $('#scroll');
  			var height = wtf[0].scrollHeight;
  			wtf.scrollTop(height);
		});
        $(".removeExcept").click(function () {
            const that = $(this);
            const name = that.data("name");
            $.ajax({
                url: "/except/delete",
                type: "post",
                dataType: 'html',
                data: {ajax_data: name},
            }).done(function () {
                that.parent().remove()
            }).fail(function () {
                alert("un problème est survenu")
            })
        })

        setTimeout(function () {
            $(".alert").alert('close')
        }, 3000)

		$('.input-group-append').on("click", function () {
            const that = $(this);
            const parent = that.parent();
            const val = parent.children(".input-group-prepend").children(".text").text();
            const input = parent.children(".newName").val();
            if (input !== "") {
                $.ajax({
                    url: "/",
                    type: "post",
                    dataType: 'html',
                    data: {ajax_data: input, oldName: val},
                }).done(function () {
                    parent.parent().parent().prepend('<div class="alert alert-info alert-dismissible fade show" role="alert">\n' +
                            '                        Rechargement en cours de la page\n' +
                            '                            <button type="button" class="close" data-dismiss="alert" aria-label="Close">\n' +
                            '                                <span aria-hidden="true">&times;</span>\n' +
                            '                            </button>\n' +
                            '                        </div>');
                    setTimeout(function () {
                        location.reload();
                    }, 2000)
                }).fail(function () {
                    alert("un problème est survenu")
                })
            }
        })

		$(".delete").click(function(e){
			that = $(this);
			e.preventDefault();
			let href = that.attr('href');
			$.ajax({
			    url: href,
			    type: 'GET',
			    success: function(result) {
					//console.log("success");
    		    	that.closest(".files").slideUp();
			    }
			});
		});

    })
</script>

</body>
</html>
`
