package myapp

//
//import (
//	"fmt"
//	"io/ioutil"
//	"log"
//	"testing"
//	"time"
//)
//
//func TestMyWatcher(t *testing.T) {
//
//	fmt.SetFlags(fmt.LstdFlags | fmt.Lshortfile)
//	var err error
//	go func() {
//		time.Sleep(2 * time.Second)
//		for _, v := range listFiles {
//			err = ioutil.WriteFile(GetEnv("dlna")+"/"+v, []byte("coucou"), 0777)
//		}
//	}()
//
//	if err != nil {
//		t.Fatal("Failed to create tmp file")
//	}
//
//	go func() {
//		fmt.Println("coucou")
//	}()
//	MyWatcher(GetEnv("dlna"))
//}
//
//var listFiles = []string{
//	"36 quai des orfèvres 1987.mkv",
//	"l ve es bel.mkv",
//	"Boku no Hero Academia S4 - 03 VOSTFR 1080p-Zone-Telechargement.NET.mkv",
//	"The.Flash.2014.S04E01.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement.Ws.mkv",
//	"Les.Tuche.French.DVDRIP-zone-telechargement.ws.avi",
//	"Kingsman.The.Golden.Circle.2017.FRENCH.BDRip.XviD-GZR.WwW.Zone-Telechargement.Ws.avi",
//	"new.girl.episode.101.trois.gars.une.fille-Zone-Telechargement.Ws.avi",
//	"Thor le monde des tenebres 720p-Shanks@Zone-Telechargement.Ws.mkv",
//	"Pitch.Perfect.3.2017.MULTI.1080p.BluRay.x264-VENUE.Zone-Telechargement1.Com.mkv",
//	"acts of violence 2018-Ww2.zone-telechargement1.com.mkv",
//	"Major 2nd - Episode 23 VOSTFR (720p)-Zone-Telechargement1 org.mkv",
//	"Inazuma Ares episode 21 VOSTFR-Zone-Telechargement1 org.mkv",
//	"Inazuma Ares episode 22 VOSTFR-Zone-Telechargement1 org.mkv",
//	"The.100.S05E13.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement1.ORG.mkv",
//	"la-nonne.mkv",
//	"The.Nun.2018.MULTi.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
//	"The.Nun.VOSTFR.2018.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
//	"The.Nun.MULTI.2018.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
//	"The_Vampire Diaries S04E23 MULTI Ici ou Ailleurs  BluRay720p  2013.mkv",
//	"Project.Blue.Book.S01E01.FRENCH.720p.HDTV.x264-HuSSLe.WwW.Zone-Telechargement.NET.mkv",
//}
