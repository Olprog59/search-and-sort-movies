package lib

import "sync"

type ObservableSlice struct {
	slice []sliceFile
	lock  sync.Mutex
	// Channel pour notifier les changements
	notifyChan chan []sliceFile
}

type numObs int

type sliceFile struct {
	file     string
	working  bool
	duration string
}

func NewObservableSlice() *ObservableSlice {
	return &ObservableSlice{
		slice:      make([]sliceFile, 10),
		notifyChan: make(chan []sliceFile),
	}
}

const (
	NotRunning = iota
	Running
)

func (o *ObservableSlice) Add(item sliceFile) {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.slice = append(o.slice, item)
	// Notifie les observateurs du nouveau contenu de la slice
	o.notifyChan <- o.slice
}

func (o *ObservableSlice) Remove(item string) {
	o.lock.Lock()
	defer o.lock.Unlock()

	for k, v := range o.slice {
		if v.file == item {
			o.slice = append(o.slice[:k], o.slice[k+1:]...)
			// Notifie les observateurs du nouveau contenu de la slice
			o.notifyChan <- o.slice
			return
		}
	}
}

func (o *ObservableSlice) Get() []sliceFile {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.slice
}

func (o *ObservableSlice) SameItem(item string) bool {
	o.lock.Lock()
	defer o.lock.Unlock()

	for _, v := range o.slice {
		if v.file == item {
			return true
		}
	}
	return false
}

func (o *ObservableSlice) Watch() <-chan []sliceFile {
	return o.notifyChan
}
