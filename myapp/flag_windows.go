package myapp

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	buildVersion string
	buildHash    string
	buildDate    string
	buildClean   string
	buildName    = "search-and-sort-movies"
)

func Flags() {
	vers := flag.Bool("v", false, "Indique la version de l'application")
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	windows := flag.Bool("windows", false, "Lancer l'application sans l'invite de commandes")
	flag.Parse()

	if *vers {
		// flag.PrintDefaults()
		fmt.Printf("Name: %s\n", buildName)
		fmt.Printf("Version: %s\n", buildVersion)
		fmt.Printf("Git Commit Hash: %s\n", buildHash)
		fmt.Printf("Build Date: %s\n", buildDate)
		fmt.Printf("Built from clean source tree: %s\n", buildClean)
		fmt.Printf("OS: %s\n", runtime.GOOS)
		fmt.Printf("Architecture: %s\n", runtime.GOARCH)
		os.Exit(1)
	}

	if *scan {
		startScan(true)
	}

	if *windows {
		HiddenWindow()
	}

}
