package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/sam-docker/media-organizer/constants"
	"github.com/sam-docker/media-organizer/lib"
	"github.com/sam-docker/media-organizer/logger"
	"github.com/sam-docker/media-organizer/model"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type fileInfo struct {
	Name, Path, UniqueID, ErrorMessage, Action, Message string
	Disabled                                            bool
	IsDir                                               bool
}

func listFiles(w http.ResponseWriter, dir string) {
	movieFiles, otherFiles := classifyFiles(dir)
	logger.Info("movieFiles: %d, otherFiles: %d", len(movieFiles), len(otherFiles))
	m, err := generateHTML(movieFiles, "/change", fileTemplate, false)
	if err != nil {
		logger.Err("Error generating change HTML: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	o, err := generateHTML(otherFiles, "/remove", fileTemplate, true)
	if err != nil {
		logger.Err("Error generating remove HTML: %s", err)
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
	re := regexp.MustCompile(constants.RegexFileExtension)

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
		logger.Warn("Tu as renommé : %s en %s", origin, data.Path)
		if err := os.Rename(origin, data.Path); err != nil {
			data.ErrorMessage = fmt.Sprintf("Un problème est survenu lors du renommage (%s)", err.Error())
		} else {
			data.ErrorMessage = "Le fichier a bien été renommé. Il va être traité d'ici quelques instants."
			constants.ObsSlice.Remove(origin)
			duration, err := lib.GetMediaDuration(data.Path)
			if err != nil {
				logger.Err("Error checking media file: %s", err)
			}
			constants.ObsSlice.Add(model.SliceFile{File: data.Path, Working: false, Duration: duration})
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
		return
	}
	return
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
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
			return
		}
		http.Error(w, "Failed to remove file: "+origin, http.StatusInternalServerError)
		return
	}

	// suppression du dossier parent si vide
	if filepath.Dir(origin) != constants.BE_SORTED {
		dir, err := os.ReadDir(filepath.Dir(origin))
		if err != nil {
			logger.Err("Error reading directory: %s", err)
		}
		if len(dir) == 0 {
			err := os.Remove(filepath.Dir(origin))
			if err != nil {
				logger.Err("Error removing directory: %s", err)
			}
		}
	}

	// Envoie une réponse vide si la suppression a réussi
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(""))
	if err != nil {
		return
	}
}

func forceHandler(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(writer, "Failed to parse form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	path := request.FormValue("path")
	typeMedia := request.FormValue("type_media")

	logger.Info("Force media type for %s to %s", path, typeMedia)
	newFileSlice := constants.ObsSlice.SetForce(path, true)
	newFileSlice.Working = false
	newFileSlice.TypeMedia = typeMedia
	if newFileSlice != nil {
		constants.ObsSlice.Remove(path)
		constants.ObsSlice.Add(*newFileSlice)
	}
	// Envoie une réponse vide si la suppression a réussi
	http.StatusText(http.StatusOK)
	return
}

func generateHTML(files []fileInfo, action, temp string, disable bool) (string, error) {
	// Prépare les données pour le template
	for i := range files {
		files[i].UniqueID = uuid.New().String()
		files[i].Action = action
		files[i].Disabled = disable
		if constants.ObsSlice.GetByName(files[i].Path) != nil {
			files[i].Message = constants.ObsSlice.GetByName(files[i].Path).TypeMedia
		}
	}

	// Parse et exécute le template
	tmpl, err := template.New("fileTemplate").Parse(temp)
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
<div class="file {{if .Disabled}}remove{{end}}" id="file-{{.UniqueID}}">
    <form hx-post="{{.Action}}" hx-swap="outerHTML" hx-target="#file-{{.UniqueID}}" hx-indicator="#loading-{{.UniqueID}}">
        <input type="text" name="filename" value="{{.Name}}" {{if .Disabled}}disabled{{end}}/>
        <input type="hidden" name="origin" value="{{.Path}}"/>
        <input type="hidden" name="uuid" value="{{.UniqueID}}"/>
        <button type="submit">{{if .Disabled}}&#x2717;{{else}}&#x2714;{{end}}</button>
        <span id="loading-{{.UniqueID}}" class="htmx-indicator"></span>
    </form>
	<div id='error-message-{{.UniqueID}}'>{{.ErrorMessage}}</div>
	{{if .Message}}
	<span id="message-{{.UniqueID}}">
		Ce fichier est identifié comme {{if eq .Message "movie"}}un <b>film</b>{{else}}une <b>série</b>{{end}}. Si cette classification vous convient, vous avez la possibilité de renommer le fichier pour qu'il soit clairement reconnu comme tel. Si vous préférez, vous pouvez également spécifier explicitement qu'il s'agit {{if eq .Message "movie"}}d'une <b>série</b>{{else}}d'un <b>film</b>{{end}} en cliquant sur ce bouton.
		<button hx-post="/force" hx-target="#message-{{.UniqueID}}" hx-swap="outerHTML" hx-vals='{"path": "{{.Path}}", "type_media": "{{if eq .Message "movie"}}serie{{else}}film{{end}}"}'>Forcer en tant que {{if eq .Message "movie"}}<b>série</b>{{else}}<b>film</b>{{end}}</button>
	</span>
	{{end}}
</div>
{{end}}
`

const changeTemplate = `
<div class="file" id="file-{{.UniqueID}}">
    <form hx-post="/change" hx-swap="outerHTML" hx-target="#file-{{.UniqueID}}" hx-indicator="#loading-{{.UniqueID}}">
        <input type="text" name="filename" value="{{.Name}}"/>
        <input type="hidden" name="origin" value="{{.Path}}"/>
        <input type="hidden" name="uuid" value="{{.UniqueID}}"/>
        <button type="submit">&#x2714;</button>
        <span id="loading-{{.UniqueID}}" class="htmx-indicator"></span>
    </form>
    <div id='error-message-{{.UniqueID}}'>{{.ErrorMessage}}</div>
	{{if .Message}}
	<span id="message-{{.UniqueID}}">
		Ce fichier est identifié comme {{if eq .Message "movie"}}un <b>film</b>{{else}}une <b>série</b>{{end}}. Si cette classification vous convient, vous avez la possibilité de renommer le fichier pour qu'il soit clairement reconnu comme tel. Si vous préférez, vous pouvez également spécifier explicitement qu'il s'agit {{if eq .Message "movie"}}d'une <b>série</b>{{else}}d'un <b>film</b>{{end}} en cliquant sur ce bouton.
		<button hx-post="/force" hx-target="#message-{{.UniqueID}}" hx-swap="outerHTML" hx-vals='{"path": "{{.Path}}", "type_media": "{{if eq .Message "movie"}}serie{{else}}film{{end}}"}'>Forcer en tant que {{if eq .Message "movie"}}<b>série</b>{{else}}<b>film</b>{{end}}</button>
	</span>
	{{end}}
</div>
`

const removeTemplate = `
<div class="file remove" id="file-{{.UniqueID}}">
    <form hx-post="/remove" hx-swap="outerHTML" hx-target="#file-{{.UniqueID}}" hx-indicator="#loading-{{.UniqueID}}">
        <input type="text" name="filename" value="{{.Name}}" disabled/>
        <input type="hidden" name="origin" value="{{.Path}}"/>
        <input type="hidden" name="uuid" value="{{.UniqueID}}"/>
        <button type="submit">&#x2717;</button>
        <span id="loading-{{.UniqueID}}" class="htmx-indicator"></span>
    </form>
    <div id='error-message-{{.UniqueID}}'>{{.ErrorMessage}}</div>
</div>
`
