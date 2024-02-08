package model

import "sync"

type ObservableSlice struct {
	Slice []SliceFile
	Lock  sync.Mutex
	// Channel pour notifier les changements
	notifyChan chan []SliceFile
}

type numObs int

type SliceFile struct {
	File      string
	Working   bool
	Duration  string
	TypeMedia string
	Force     bool
}

func NewObservableSlice() *ObservableSlice {
	return &ObservableSlice{
		Slice:      make([]SliceFile, 0),
		notifyChan: make(chan []SliceFile),
	}
}

const (
	NotRunning = iota
	Running
)

func (o *ObservableSlice) Add(item SliceFile) {
	o.Lock.Lock()
	defer o.Lock.Unlock()

	o.Slice = append(o.Slice, item)
	// Notifie les observateurs du nouveau contenu de la slice
	o.notifyChan <- o.Slice
}

func (o *ObservableSlice) Remove(item string) {
	o.Lock.Lock()
	defer o.Lock.Unlock()

	for k, v := range o.Slice {
		if v.File == item {
			o.Slice = append(o.Slice[:k], o.Slice[k+1:]...)
			// Notifie les observateurs du nouveau contenu de la slice
			o.notifyChan <- o.Slice
			return
		}
	}
}

func (o *ObservableSlice) Get() []SliceFile {
	o.Lock.Lock()
	defer o.Lock.Unlock()
	return o.Slice
}

func (o *ObservableSlice) GetByName(name string) *SliceFile {
	o.Lock.Lock()
	defer o.Lock.Unlock()

	for _, v := range o.Slice {
		if v.File == name {
			return &v
		}
	}
	return nil
}

func (o *ObservableSlice) SameItem(item string) bool {
	o.Lock.Lock()
	defer o.Lock.Unlock()

	for _, v := range o.Slice {
		if v.File == item {
			return true
		}
	}
	return false
}

func (o *ObservableSlice) Watch() <-chan []SliceFile {
	return o.notifyChan
}

func (o *ObservableSlice) SetTypeMedia(file string, s string) {
	o.Lock.Lock()
	defer o.Lock.Unlock()

	for k, v := range o.Slice {
		if v.File == file {
			o.Slice[k].TypeMedia = s
			return
		}
	}
}

func (o *ObservableSlice) SetForce(path string, b bool) *SliceFile {
	o.Lock.Lock()
	defer o.Lock.Unlock()

	for k, v := range o.Slice {
		if v.File == path {
			o.Slice[k].Force = b
			return &o.Slice[k]
		}
	}
	return nil
}
