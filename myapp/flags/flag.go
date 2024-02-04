package flags

import (
	"flag"
	"media-organizer/myapp"
)

func Flags() {
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	flag.Parse()

	if *scan {
		myapp.StartScan()
	}
}
