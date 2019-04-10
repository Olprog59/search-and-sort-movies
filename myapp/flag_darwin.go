package myapp

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func Flags(BuildName, BuildVersion, BuildHash, BuildDate, BuildClean string) {
	vers := flag.Bool("v", false, "Indique la version de l'application")
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	sendmail := flag.Bool("email", false, "Envoie email test")
	flag.Parse()

	if *vers {
		// flag.PrintDefaults()
		fmt.Printf("Name: %s\n", BuildName)
		fmt.Printf("Version: %s\n", BuildVersion)
		fmt.Printf("Git Commit Hash: %s\n", BuildHash)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Built from clean source tree: %s\n", BuildClean)
		fmt.Printf("OS: %s\n", runtime.GOOS)
		fmt.Printf("Architecture: %s\n", runtime.GOARCH)
		os.Exit(1)
	}

	if *scan {
		startScan()
	}

	if *sendmail{
		SendMail("", "")
		os.Exit(1)
	}

}
