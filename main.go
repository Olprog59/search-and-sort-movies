package main

import (
	"embed"
	"github.com/sam-docker/media-organizer/constants"
	"github.com/sam-docker/media-organizer/flags"
	"github.com/sam-docker/media-organizer/lib"
	"github.com/sam-docker/media-organizer/logger"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
)

//go:embed all:html
var content embed.FS

//go:embed all:statics
var statics embed.FS

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.SetFlags(0)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	checkIfFolderExistAndCreate(constants.BE_SORTED)
	checkIfFolderExistAndCreate(constants.MOVIES)
	checkIfFolderExistAndCreate(constants.SERIES)
}

func main() {

	// Activation du serveur de profilage
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go logger.ManageClients()

	logger.L(logger.Magenta, "Start :-D")

	// Nettoyer le dossier de be_sorted
	lib.CleanFolder(constants.BE_SORTED)

	go lib.MyWatcher(constants.BE_SORTED, constants.ObsSlice)

	go func() {
		for slice := range constants.ObsSlice.Watch() {
			lib.StartProcessing(slice, constants.ObsSlice)
			logger.L(logger.Magenta, "%#+v", slice)
		}
	}()

	// scan
	flags.Flags()

	// Start server
	mux := http.NewServeMux()
	mid := logMiddleware(mux)
	sub, err := fs.Sub(statics, "statics")
	if err != nil {
		return
	}

	mux.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.FS(sub))))
	mux.HandleFunc("GET /{$}", indexHandler)
	mux.HandleFunc("GET /scan", scanHandler)
	mux.HandleFunc("POST /change", changeHandler)
	mux.HandleFunc("POST /remove", removeHandler)
	mux.HandleFunc("POST /force", forceHandler)
	// Route pour les logs SSE
	mux.HandleFunc("GET /logs", logger.ServeLogs)

	logger.L(logger.Magenta, "Start server on localhost:8080")
	err = http.ListenAndServe("localhost:8080", mid)
	logger.L(logger.Yellow, "L'application a été arrêtée: %s", err)
}

func scanHandler(w http.ResponseWriter, _ *http.Request) {
	listFiles(w, constants.BE_SORTED)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFS(content, "html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
}

func logMiddleware(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// force UTF-8 for all requests
		if !strings.HasPrefix(r.URL.Path, "/statics/") {
			logger.L(logger.Teal, "%s %s %s", r.RemoteAddr, r.Method, r.URL)
		}
		mux.ServeHTTP(w, r)
	})
}

func checkIfFolderExistAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0766)
	}
}
