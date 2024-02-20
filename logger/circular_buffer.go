package logger

import (
	"log"
	"runtime"
	"sync"
)

type CircularBuffer struct {
	logs  []string   // Slice pour stocker les logs
	size  int        // Taille maximale du buffer
	index int        // Index actuel pour l'insertion
	mu    sync.Mutex // Mutex pour rendre le buffer sûr pour les goroutines
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		logs:  make([]string, size),
		size:  size,
		index: 0,
	}
}

func (cb *CircularBuffer) Append(logs string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Gestionnaire de panique
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in GetAll:", r)
			// Imprimer la pile d'exécution pour le débogage
			buf := make([]byte, 1024)
			runtime.Stack(buf, false)
			log.Println(string(buf))
		}
	}()

	cb.logs[cb.index] = logs            // Insère le log à l'index actuel
	cb.index = (cb.index + 1) % cb.size // Incrémente l'index et le remet à 0 si on dépasse la taille
}

func (cb *CircularBuffer) GetAll() (logs []string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Gestionnaire de panique
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in GetAll:", r)
			// Imprimer la pile d'exécution pour le débogage
			buf := make([]byte, 1024)
			runtime.Stack(buf, false)
			log.Println(string(buf))
		}
	}()

	L(Magenta, "GetAll - Debut")

	if cb.logs[cb.index] == "" {
		logs = cb.logs[:cb.index]
	} else {
		logs = append(cb.logs[cb.index:], cb.logs[:cb.index]...)
	}

	L(Magenta, "GetAll - Fin")

	return logs
}
