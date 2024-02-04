package logger

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

//Black   = Color("\033[1;30m%s\033[0m")
//Green   = Color("\033[1;32m%s\033[0m")
//White   = Color("\033[1;37m%s\033[0m")

var (
	Red     = Color("\033[1;31m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
)

//func getFileAndLine() string {
//	_, file, line, ok := runtime.Caller(3)
//	if !ok {
//		panic("Could not get context info for logger!")
//	}
//	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
//}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		//msg := fmt.Sprintf(colorString, fmt.Sprint(time.Now().Format("02-01-2006 15:04:05")+": "+getFileAndLine()+": "+fmt.Sprint(args...)))
		msg := fmt.Sprintf(colorString, fmt.Sprint(time.Now().Format("02-01-2006 15:04:05")+": "+fmt.Sprint(args...)))
		LogMessage(msg)   // Ajoute à l'historique
		broadcastLog(msg) // Envoie aux clients SSE
		return msg
	}
	return sprint
}

func L(color func(...interface{}) string, message string, param ...interface{}) {
	fmt.Println(color(fmt.Sprintf(message, param...)))
}

var logBuffer = NewCircularBuffer(1000) // Initialise le buffer avec une taille de 1000

func LogMessage(message string) {
	logBuffer.Append(message) // Ajoute un nouveau log
}

func ServeLogs(w http.ResponseWriter, _ *http.Request) {
	// Initialiser le client SSE
	client := sseClient{
		id: uuid.NewString(), // Générer un nouvel ID de client, nécessite "github.com/google/uuid"
		ch: make(chan string),
	}
	register <- client
	defer func() { unregister <- client }()

	// Configurer les en-têtes SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Envoyer l'historique des logs
	logs := logBuffer.GetAll()
	for _, log := range logs {
		_, err := fmt.Fprintf(w, "data: %s\n\n", log)
		if err != nil {
			return
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	// Boucle pour envoyer les nouveaux logs
	for log := range client.ch {
		_, err := fmt.Fprintf(w, "data: %s\n\n", log)
		if err != nil {
			return
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}
