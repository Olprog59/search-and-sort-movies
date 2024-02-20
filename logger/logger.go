package logger

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
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
	Green   = Color("\033[1;32m%s\033[0m")
)

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		log.Println("Failed to get the caller information")
	}
	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		msg := fmt.Sprintf(colorString, fmt.Sprint(time.Now().Format("02-01-2006 15:04:05")+": "+getFileAndLine()+": "+fmt.Sprint(args...))+"\n")
		//msg := fmt.Sprintf(colorString, fmt.Sprint(time.Now().Format("02-01-2006 15:04:05")+": "+fmt.Sprint(args...)))
		logMessage(msg)   // Ajoute à l'historique
		broadcastLog(msg) // Envoie aux clients SSE
		return msg
	}
	return sprint
}

func L(color func(...interface{}) string, message string, param ...interface{}) {
	fmt.Printf(color(fmt.Sprintf(message, param...)))
}

var logBuffer = NewCircularBuffer(1000) // Initialise le buffer avec une taille de 1000

func logMessage(message string) {
	logBuffer.Append(message) // Ajoute un nouveau log
}

func ServeLogs(w http.ResponseWriter, r *http.Request) {
	client := sseClient{
		id: uuid.NewString(),
		ch: make(chan string),
	}
	register <- client
	defer func() { unregister <- client }()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	//logs := logBuffer.GetAll()
	//for _, log := range logs {
	//	if _, err := fmt.Fprintf(w, "data: %s\n\n", log); err != nil {
	//		return
	//	}
	//	flusher.Flush()
	//}

	// Envoie un commentaire toutes les 30 secondes pour garder la connexion ouverte
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	ctx := r.Context()

	for {
		select {
		case logg := <-client.ch:
			if _, err := fmt.Fprintf(w, "data: %s\n\n", logg); err != nil {
				return
			}
			flusher.Flush()
		case <-ticker.C:
			if _, err := fmt.Fprintf(w, ": keep-alive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case <-ctx.Done():
			return // Fermer la connexion lorsque le client se déconnecte

		}
	}
}
