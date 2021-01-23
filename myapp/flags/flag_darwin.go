package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"search-and-sort-movies/myapp"
	"search-and-sort-movies/myapp/logger"
	"search-and-sort-movies/myapp/model"
)

func Flags(BuildName, BuildVersion, BuildHash, BuildDate, BuildClean string) {
	vers := flag.Bool("v", false, "Indique la version de l'application")
	scan := flag.Bool("scan", false, "Lancer le scan au démarrage de l'application")
	jsonFormat := flag.Bool("j", false, "Retour json")
	flag.Parse()

	buildInfo := model.BuildInfo{
		BuildName:    BuildName,
		BuildVersion: BuildVersion,
		BuildDate:    BuildDate,
		BuildClean:   BuildClean,
		BuildHash:    BuildHash,
		Os:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	if *vers {
		if *jsonFormat {
			prettyJson, err := json.MarshalIndent(&buildInfo, "", " ")
			if err != nil {
				fmt.Println(logger.Warn(err))
			}
			fmt.Printf("%s\n", string(prettyJson))
			os.Exit(1)
		}
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
		myapp.StartScan()
	}

}
