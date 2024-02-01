package main

import (
	"embed"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"search-and-sort-movies/myapp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/flags"
	"search-and-sort-movies/myapp/logger"
	"time"
)

//go:embed all:templates
var content embed.FS

//go:embed all:statics
var statics embed.FS

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.SetFlags(0)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	fmt.Println("path: " + constants.A_TRIER)
	if constants.ALL != "" {
		checkIfFolderExistAndCreate(constants.ALL)
		constants.A_TRIER = constants.ALL + "/be_sorted"
		constants.MOVIES = constants.ALL + "/movies"
		constants.SERIES = constants.ALL + "/series"
	}
	checkIfFolderExistAndCreate(constants.A_TRIER)
	checkIfFolderExistAndCreate(constants.MOVIES)
	checkIfFolderExistAndCreate(constants.SERIES)
}

func main() {
	logger.L(logger.Magenta, "Start :-D")
	// scan
	flags.Flags()

	go myapp.MyWatcher(constants.A_TRIER)

	mux := http.NewServeMux()
	mid := logMiddleware(mux)
	sub, err := fs.Sub(statics, "statics")
	if err != nil {
		return
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/scan", scanHandler)
	mux.HandleFunc("/change", changeHandler)
	mux.HandleFunc("/remove", removeHandler)

	logger.L(logger.Magenta, "Start server on localhost:8080")
	err = http.ListenAndServe(":8080", mid)
	if err != nil {
		return
	}
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	listFiles(w, constants.A_TRIER)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		origin := r.FormValue("origin")
		err = os.Remove(origin)
		if err != nil {
			uniqueID := uuid.New().String()
			otherFilesString := fmt.Sprintf(`<div class="file remove" id="file-%s">
	<form hx-post="/remove" hx-swap="outerHTML"  hx-target="#file-%s">
		<input type="text" name="filename" value="%s" disabled/>
		<input type="hidden" name="origin" value="%s"/>
		<input type="hidden" name="uuid" value="%s"/>
		<button type="submit">&#x2717;</button>
	</form>
<div id='error-message-%s'>Un problème est survenue lors de la suppression (%s)</div>
</div>`, uniqueID, uniqueID, filepath.Base(origin), origin, uniqueID, uniqueID, err.Error())
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(otherFilesString))
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(""))
		return
	} else {
		http.Error(w, "Only POST method", http.StatusInternalServerError)
	}
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		origin := r.FormValue("origin")
		newValue := r.FormValue("filename")
		id := r.FormValue("uuid")

		if filepath.Base(origin) == newValue {
			errorForm := fmt.Sprintf(`
				<div class="file" id="file-%s">
					<form hx-post="/change" hx-swap="outerHTML" hx-target="#file-%s">
						<input type="text" name="filename" value="%s">
						<input type="hidden" name="origin" value="%s">
						<input type="hidden" name="uuid" value="%s"/>
						<button type="submit">✔</button>
					</form>
					<div id='error-message-%s'>Aucun changement détecté</div>
					<script>
						setTimeout(function() {
							var errorDiv = document.getElementById("error-message-%s");
							if (errorDiv) {
								errorDiv.style.display = "none";
							}
						}, 5000);
					</script>
				</div>`, id, id, newValue, filepath.Dir(origin)+string(filepath.Separator)+newValue, id, id, id)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(errorForm))
			return
		}

		err = os.Rename(origin, filepath.Dir(origin)+string(filepath.Separator)+newValue)
		if err != nil {
			errorForm := fmt.Sprintf(`
				<div class="file" id="file-%s">
					<form hx-post="/change" hx-swap="outerHTML" hx-target="#file-%s">
						<input type="text" name="filename" value="%s">
						<input type="hidden" name="origin" value="%s">
						<input type="hidden" name="uuid" value="%s"/>
						<button type="submit">✔</button>
					</form>
					<div id='error-message-%s'>Un problème est apparu lors du renommage</div>
					<script>
						setTimeout(function() {
							var errorDiv = document.getElementById("error-message-%s");
							if (errorDiv) {
								errorDiv.style.display = "none";
							}
						}, 5000);
					</script>
				</div>`, id, id, newValue, filepath.Dir(origin)+string(filepath.Separator)+newValue, id, id, id)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(errorForm))
			return
		}
		// Après renommage réussi, renvoie le formulaire mis à jour :
		updatedForm := fmt.Sprintf(`
				<div class="file" id="file-%s">
					<form hx-post="/change" hx-swap="outerHTML" hx-target="#file-%s">
						<input type="text" name="filename" value="%s">
						<input type="hidden" name="origin" value="%s">
						<input type="hidden" name="uuid" value="%s"/>
						<button type="submit">✔</button>
					</form>
					<div id='error-message-%s blue'>Le fichier a bien été renommé. Attendez quelques secondes afin de vérifier si il a bien été déplacé.</div>
					<script>
						setTimeout(function() {
							var errorDiv = document.getElementById("error-message-%s");
							if (errorDiv) {
								errorDiv.style.display = "none";
							}
						}, 5000);
					</script>
				</div>`, id, id, newValue, filepath.Dir(origin)+string(filepath.Separator)+newValue, id, id, id)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(updatedForm))
		return
	} else {
		http.Error(w, "Only POST method", http.StatusInternalServerError)
	}
}

func logMiddleware(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// force UTF-8 for all requests
		logger.L(logger.Teal, "%s %s %s", r.RemoteAddr, r.Method, r.URL)
		mux.ServeHTTP(w, r)
	})
}

func listFiles(w http.ResponseWriter, dir string) {

	var movieFiles []struct{ Name, Path string }
	var otherFiles []struct{ Name, Path string }

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		re := regexp.MustCompile(constants.RegexFile)
		if re.MatchString(filepath.Ext(path)) {
			movieFiles = append(movieFiles, struct{ Name, Path string }{filepath.Base(path), path})
		} else {
			otherFiles = append(otherFiles, struct{ Name, Path string }{filepath.Base(path), path})
		}
		return nil
	})

	w.Header().Set("Content-Type", "text/html")
	movieFilesString := ""
	for _, file := range movieFiles {
		uniqueID := uuid.New().String()
		movieFilesString += fmt.Sprintf(`<div class="file" id="file-%s">
	<form hx-post="/change" hx-swap="outerHTML" hx-target="#file-%s">
		<input type="text" name="filename" value="%s" />
		<input type="hidden" name="origin" value="%s"/>
		<input type="hidden" name="uuid" value="%s"/>
		<button type="submit">&#x2714;</button>
	</form></div>`, uniqueID, uniqueID, file.Name, file.Path, uniqueID)
	}

	otherFilesString := ""
	for _, file := range otherFiles {
		uniqueID := uuid.New().String()
		otherFilesString += fmt.Sprintf(`<div class="file remove" id="file-%s">
	<form hx-post="/remove" hx-swap="outerHTML"  hx-target="#file-%s">
		<input type="text" name="filename" value="%s" disabled/>
		<input type="hidden" name="origin" value="%s"/>
		<input type="hidden" name="uuid" value="%s"/>
		<button type="submit">&#x2717;</button>
	</form></div>`, uniqueID, uniqueID, file.Name, file.Path, uniqueID)
	}

	_, err := w.Write([]byte(fmt.Sprintf(`<div id="error-container"></div><div class="files"><h3>Renommage de </h3>%s</div><div class="files other">%s</div>`, movieFilesString, otherFilesString)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return

}

func checkIfFolderExistAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0766)
	}
}

func SetFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Name: name, Value: encode(value)}
	http.SetCookie(w, c)
}

func GetFlash(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return nil, err
	}
	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return value, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
