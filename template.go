package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
)

type fileInfo struct {
	Name, Path, UniqueID, ErrorMessage, Action string
	Disabled                                   bool
}

func listFiles(w http.ResponseWriter, dir string) {
	movieFiles, otherFiles := classifyFiles(dir)
	m, err := generateHTML(movieFiles, "/change", false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	o, err := generateHTML(otherFiles, "/remove", true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write([]byte(fmt.Sprintf(`<div class="files" id="files">%s</div><div class="files other">%s</div>`, m, o)))
	if err != nil {
		return
	}
}

func classifyFiles(dir string) ([]fileInfo, []fileInfo) {
	var movieFiles, otherFiles []fileInfo
	re := regexp.MustCompile(constants.RegexFile)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		file := fileInfo{Name: filepath.Base(path), Path: path}
		if re.MatchString(filepath.Ext(path)) {
			movieFiles = append(movieFiles, file)
		} else {
			otherFiles = append(otherFiles, file)
		}
		return nil
	})
	if err != nil {
		return nil, nil
	}

	return movieFiles, otherFiles
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	origin := r.FormValue("origin")
	newValue := r.FormValue("filename")
	uniqueID := r.FormValue("uuid")

	// Préparation des données pour le template
	data := fileInfo{
		UniqueID: uniqueID,
		Action:   "/change",
		Name:     newValue,
		Path:     filepath.Dir(origin) + string(filepath.Separator) + newValue,
	}

	// Vérifie si le nom du fichier a changé
	if filepath.Base(origin) == newValue {
		data.ErrorMessage = "Aucun changement détecté"
	} else {
		logger.L(logger.Magenta, "Rename file %s to %s", origin, data.Path)
		if err := os.Rename(origin, data.Path); err != nil {
			data.ErrorMessage = fmt.Sprintf("Un problème est survenu lors du renommage (%s)", err.Error())
		} else {
			data.ErrorMessage = "Le fichier a bien été renommé. Il va être traité d'ici quelques instants."
		}
	}

	// Génère la réponse HTML en utilisant le template
	tmpl, err := template.New("change").Parse(changeTemplate)
	if err != nil {
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	origin := r.FormValue("origin")
	uniqueID := uuid.New().String() // Génère un ID unique pour cette opération

	if err := os.Remove(origin); err != nil {
		// Préparation des données pour le template en cas d'erreur
		data := fileInfo{
			UniqueID:     uniqueID,
			Name:         filepath.Base(origin),
			Path:         origin,
			ErrorMessage: fmt.Sprintf("Un problème est survenue lors de la suppression (%s)", err.Error()),
		}

		tmpl, err := template.New("remove").Parse(removeTemplate)
		if err != nil {
			http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Envoie une réponse vide si la suppression a réussi
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(""))
	if err != nil {
		return
	}
}

func generateHTML(files []fileInfo, action string, disabled bool) (string, error) {
	// Prépare les données pour le template
	for i := range files {
		files[i].UniqueID = uuid.New().String()
		files[i].Disabled = disabled
		files[i].Action = action
	}

	// Parse et exécute le template
	tmpl, err := template.New("fileTemplate").Parse(fileTemplate)
	if err != nil {
		return "", err
	}

	var html bytes.Buffer
	err = tmpl.Execute(&html, map[string]interface{}{
		"Files": files,
	})
	if err != nil {
		return "", err
	}

	return html.String(), nil
}

const fileTemplate = `
{{range .Files}}
<div class="file" id="file-{{.UniqueID}}">
    <form hx-post="{{.Action}}" hx-swap="outerHTML" hx-target="#file-{{.UniqueID}}" hx-indicator="#loading-{{.UniqueID}}">
        <input type="text" name="filename" value="{{.Name}}" {{if .Disabled}}disabled{{end}}/>
        <input type="hidden" name="origin" value="{{.Path}}"/>
        <input type="hidden" name="uuid" value="{{.UniqueID}}"/>
        <button type="submit">&#x2714;</button>
        <span id="loading-{{.UniqueID}}" class="htmx-indicator"></span>
    </form>
	<div id='error-message-{{.UniqueID}}'>{{.ErrorMessage}}</div>
</div>
{{end}}
`

const changeTemplate = `
<div class="file" id="file-{{.UniqueID}}">
    <form hx-post="{{.Action}}" hx-swap="outerHTML" hx-target="#file-{{.UniqueID}}" hx-indicator="#loading-{{.UniqueID}}">
        <input type="text" name="filename" value="{{.Name}}" {{if .Disabled}}disabled{{end}}/>
        <input type="hidden" name="origin" value="{{.Path}}"/>
        <input type="hidden" name="uuid" value="{{.UniqueID}}"/>
        <button type="submit">&#x2714;</button>
        <span id="loading-{{.UniqueID}}" class="htmx-indicator"></span>
    </form>
    <div id='error-message-{{.UniqueID}}'>{{.ErrorMessage}}</div>
</div>
`

const removeTemplate = `
<div class="file remove" id="file-{{.UniqueID}}">
    <form hx-post="/remove" hx-swap="outerHTML" hx-target="#file-{{.UniqueID}}" hx-indicator="#loading-{{.UniqueID}}">
        <input type="text" name="filename" value="{{.Name}}" disabled/>
        <input type="hidden" name="origin" value="{{.Path}}"/>
        <input type="hidden" name="uuid" value="{{.UniqueID}}"/>
        <button type="submit">&#x2714;</button>
        <span id="loading-{{.UniqueID}}" class="htmx-indicator"></span>
    </form>
    <div id='error-message-{{.UniqueID}}'>{{.ErrorMessage}}</div>
</div>
`
