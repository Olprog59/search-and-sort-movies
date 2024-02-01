package flags

import (
	"flag"
	"search-and-sort-movies/myapp"
)

func Flags() {
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	flag.Parse()

	if *scan {
		myapp.StartScan()
	}
}
