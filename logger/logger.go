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

//Black   = color("\033[1;30m%s\033[0m")
//Success   = color("\033[1;32m%s\033[0m")
//White   = color("\033[1;37m%s\033[0m")

var (
	err     = color("\033[1;31mErr\t: %s\033[0m")
	warn    = color("\033[1;33mWarn\t: %s\033[0m")
	debug   = color("\033[1;34mDebug\t: %s\033[0m")
	info    = color("\033[1;35mInfo\t: %s\033[0m")
	notice  = color("\033[1;36mNotice\t: %s\033[0m")
	success = color("\033[1;32mSuccess\t: %s\033[0m")
)

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		log.Println("Failed to get the caller information")
	}
	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
}

var logBuffer = NewCircularBuffer(1000) // Initialise le buffer avec une taille de 1000

func color(colorString string) func(message string, args ...any) string {
	return func(message string, args ...any) string {
		// Construction directe du message avec les arguments, sans double formatage
		formattedMessage := fmt.Sprintf(message, args...)
		timestamp := time.Now().Format("02-01-2006 15:04:05")
		fileLine := getFileAndLine()
		fullMessage := fmt.Sprintf("%s - %s\t: %s", timestamp, fileLine, formattedMessage)
		msg := fmt.Sprintf(colorString, fullMessage)

		logBuffer.Append(msg) // Ajoute à l'historique
		broadcastLog(msg)     // Envoie aux clients SSE
		return msg
	}
}

func l(colorFn func(message string, args ...any) string, message string, params ...any) {
	fmt.Println(colorFn(message, params...))
}

func Err(message string, params ...any) {
	l(err, message, params...)
}
func Warn(message string, params ...any) {
	l(warn, message, params...)
}
func Debug(message string, params ...any) {
	l(debug, message, params...)
}
func Info(message string, params ...any) {
	l(info, message, params...)
}
func Notice(message string, params ...any) {
	l(notice, message, params...)
}
func Success(message string, params ...any) {
	l(success, message, params...)
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

	logs := logBuffer.GetAll()

	for _, l := range logs {
		if _, err := fmt.Fprintf(w, "data: %s\n\n", l); err != nil {
			return
		}
		flusher.Flush()
	}

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
