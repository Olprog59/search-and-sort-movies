package myapp

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func MyWatcher(location string) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watch.Events:
				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					folder, file := filepath.Split(event.Name)
					re := regexp.MustCompile(regexFile)
					if re.MatchString(filepath.Ext(file)) {
						folder = filepath.Clean(folder)
						// Quand c'est en local
						//fmt.Println(folder)
						//_, end := filepath.Split(GetEnv("dlna"))
						//fmt.Println(end)
						//if end != folder {
						if GetEnv("dlna") != folder {
							if err := watch.Remove(folder); err != nil {
								log.Println(err)
							}
						}
					}
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					_, file := filepath.Split(event.Name)
					f, _ := os.Stat(event.Name)
					if f.IsDir() {
						log.Println(event.Name)
						err = filepath.Walk(event.Name, func(path string, info os.FileInfo, err error) error {
							//files = append(files, path)
							re := regexp.MustCompile(regexFile)
							if re.MatchString(filepath.Ext(path)) {
								println("c'est parti pour " + filepath.Dir(path))
								err = watch.Add(filepath.Dir(path))
								if err != nil {
									print(err)
								}
							}
							return nil
						})
					}
					re := regexp.MustCompile(regexFile)
					if re.MatchString(filepath.Ext(file)) {
						log.Println("Détection de : ", file)
						Process(event.Name)
					}
				}
			case err := <-watch.Errors:
				log.Println("error:", err)
			}
		}
	}()

	if err := watch.Add(location); err != nil {
		log.Fatal(err)
	}

	<-done
}
