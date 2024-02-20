package logger

import "sync"

type sseClient struct {
	id string
	ch chan string
}

var (
	clients      = make(map[string]sseClient) // Garder une trace des clients SSE
	clientsMutex = sync.Mutex{}               // Assure la sécurité des goroutines pour la manipulation des clients
	register     = make(chan sseClient)       // Canal pour enregistrer de nouveaux clients
	unregister   = make(chan sseClient)       // Canal pour désenregistrer les clients
)

// Fonction pour écouter les canaux register et unregister
func ManageClients() {
	for {
		select {
		case client := <-register:
			clientsMutex.Lock()
			clients[client.id] = client
			clientsMutex.Unlock()
		case client := <-unregister:
			clientsMutex.Lock()
			delete(clients, client.id)
			close(client.ch)
			clientsMutex.Unlock()
		}
	}
}

// Fonction pour envoyer un log à tous les clients SSE
func broadcastLog(log string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for _, client := range clients {
		client.ch <- log
	}
}
