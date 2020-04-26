package myapp

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestMyWaTcher(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	go func() {
		for _, v := range listFiles {
			time.Sleep(1 * time.Second)
			err = ioutil.WriteFile(GetEnv("dlna")+"/"+v+".mkv", []byte("coucou"), 0777)
		}
	}()

	if err != nil {
		t.Fatal("Failed to create tmp file")
	}

	go func() {
		fmt.Println("coucou")
		time.Sleep(4 * time.Second)
		_ = os.Mkdir(GetEnv("dlna")+"/dossier1/", 0777)
		time.Sleep(2 * time.Second)
		_ = os.Mkdir(GetEnv("dlna")+"/dossier1/sous-dossier/", 0777)
	}()
	MyWaTcher(GetEnv("dlna"))
}

var listFiles = []string{
	"Casino Royale",
	"Quantum of Solace",
	"Skyfall",
	"8 mile",
	"L'armée des 12 singes",
	"2009 Lost Memories",
	"21 Grams",
	"25th Hour",
	"28 jours plus tards",
	"28 semaines plus tard",
	"300",
	"36 quai des orfèvres",
	"L'age de glace",
	"Alexandre",
	"Ali",
	"Alice aux pays des merveilles",
	"Alien",
	"Alien 2",
	"Alien 3",
	"Alien 4",
	"American Gangster",
	"American History X",
	"American Psycho",
	"Any Given Sunday",
	"Apocalypse Now",
	"Armageddon",
	"L'arme Fatale",
	"L'arme Fatale 2",
	"L'arme Fatale 3",
	"L'arme Fatale 4",
	"Les associés",
	"L'associé du diable",
	"Metro 123",
	"Avatar",
	"Bad Boys",
	"Bad Boys 2",
	"Basic Instinct",
	"Batman",
	"Batman Begins",

	"dossier1/Casino Royale",
	"dossier1/Quantum of Solace",
	"dossier1/Skyfall",
	"dossier1/8 mile",
	"dossier1/L'armée des 12 singes",
	"dossier1/2009 Lost Memories",
	"dossier1/21 Grams",
	"dossier1/25th Hour",
	"dossier1/28 jours plus tards",
	"dossier1/28 semaines plus tard",
	"dossier1/300",
	"dossier1/36 quai des orfèvres",
	"dossier1/L'age de glace",
	"dossier1/Alexandre",
	"dossier1/Ali",
	"dossier1/Alice aux pays des merveilles",
	"dossier1/Alien",
	"dossier1/Alien 2",
	"dossier1/Alien 3",
	"dossier1/Alien 4",
	"dossier1/American Gangster",
	"dossier1/American History X",
	"dossier1/American Psycho",
	"dossier1/Any Given Sunday",
	"dossier1/Apocalypse Now",
	"dossier1/Armageddon",
	"dossier1/L'arme Fatale",
	"dossier1/L'arme Fatale 2",
	"dossier1/L'arme Fatale 3",
	"dossier1/L'arme Fatale 4",
	"dossier1/Les associés",
	"dossier1/L'associé du diable",
	"dossier1/Metro 123",
	"dossier1/Avatar",
	"dossier1/Bad Boys",
	"dossier1/Bad Boys 2",
	"dossier1/Basic Instinct",
	"dossier1/Batman",
	"dossier1/Batman Begins",

	"dossier1/sous-dossier/Casino Royale",
	"dossier1/sous-dossier/Quantum of Solace",
	"dossier1/sous-dossier/Skyfall",
	"dossier1/sous-dossier/8 mile",
	"dossier1/sous-dossier/L'armée des 12 singes",
	"dossier1/sous-dossier/2009 Lost Memories",
	"dossier1/sous-dossier/21 Grams",
	"dossier1/sous-dossier/25th Hour",
	"dossier1/sous-dossier/28 jours plus tards",
	"dossier1/sous-dossier/28 semaines plus tard",
	"dossier1/sous-dossier/300",
	"dossier1/sous-dossier/36 quai des orfèvres",
	"dossier1/sous-dossier/L'age de glace",
	"dossier1/sous-dossier/Alexandre",
	"dossier1/sous-dossier/Ali",
	"dossier1/sous-dossier/Alice aux pays des merveilles",
	"dossier1/sous-dossier/Alien",
	"dossier1/sous-dossier/Alien 2",
	"dossier1/sous-dossier/Alien 3",
	"dossier1/sous-dossier/Alien 4",
	"dossier1/sous-dossier/American Gangster",
	"dossier1/sous-dossier/American History X",
	"dossier1/sous-dossier/American Psycho",
	"dossier1/sous-dossier/Any Given Sunday",
	"dossier1/sous-dossier/Apocalypse Now",
	"dossier1/sous-dossier/Armageddon",
	"dossier1/sous-dossier/L'arme Fatale",
	"dossier1/sous-dossier/L'arme Fatale 2",
	"dossier1/sous-dossier/L'arme Fatale 3",
	"dossier1/sous-dossier/L'arme Fatale 4",
	"dossier1/sous-dossier/Les associés",
	"dossier1/sous-dossier/L'associé du diable",
	"dossier1/sous-dossier/Metro 123",
	"dossier1/sous-dossier/Avatar",
	"dossier1/sous-dossier/Bad Boys",
	"dossier1/sous-dossier/Bad Boys 2",
	"dossier1/sous-dossier/Basic Instinct",
	"dossier1/sous-dossier/Batman",
	"dossier1/sous-dossier/Batman Begins",
}
