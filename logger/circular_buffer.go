package logger

import (
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

func (cb *CircularBuffer) Append(l string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.logs[cb.index] = l               // Insère le log à l'index actuel
	cb.index = (cb.index + 1) % cb.size // Incrémente l'index et le remet à 0 si on dépasse la taille
}

func (cb *CircularBuffer) GetAll() []string {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	var logs []string

	if cb.logs[cb.index] == "" {
		logs = cb.logs[:cb.index]
	} else {
		logs = append(cb.logs[cb.index:], cb.logs[:cb.index]...)
	}

	return logs
}
