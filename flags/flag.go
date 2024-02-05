package flags

import (
	"flag"
	"media-organizer/lib"
)

func Flags() {
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	flag.Parse()

	if *scan {
		lib.StartScan()
	}
}
