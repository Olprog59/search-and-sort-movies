package myapp

import (
	"fmt"
	"io/ioutil"
	"log"
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

	//for _, v := range listFiles {
	//	os.Remove("/Users/olprog/go/src/search-and-sort-movies/a_trier/"+v+".mkv")
	//}
	//defer os.RemoveAll("/Users/olprog/go/src/search-and-sort-movies/a_trier/Coucou_file.mkv")

	go func() {
		fmt.Println("coucou")
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
}
