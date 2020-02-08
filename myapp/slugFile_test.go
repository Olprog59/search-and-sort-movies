package myapp

import (
	"testing"
)

func Test_slugFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name                  string
		args                  args
		wantName              string
		wantSerieName         string
		wantSerieNumberReturn string
		wantYear              int
	}{

		{"slugFile", args{
			"Pirates.of.the.Caribbean.Dead.Men.Tell.No.Tales.2017.MULTi.1080p.WEB.H264-SPACED.WwW.Zone-Telechargement.Ws.mkv",
		},
			"pirates-of-the-caribbean-dead-men-tell-no-tales.mkv",
			"",
			"",
			2017,
		},
		{"slugFile", args{
			"modern.family.s1e09.french.dvdrip.xvid-jmt-Zone-telechargement.Ws.avi",
		},
			"modern-family-s01e09.avi",
			"modern-family",
			"s01e09",
			0,
		},
		{"slugFile", args{
			"Snatched.2016.MULTi.1080p.BluRay.x264-VENUE.WwW.Zone-Telechargement.Ws.mkv",
		},
			"snatched.mkv",
			"",
			"",
			2016,
		},
		{"slugFile", args{
			"Demain.Nous.Appartient.S01E4.FRENCH.HDTV.XviD-ZT.WwW.Zone-Telechargement.Ws.avi",
		},
			"demain-nous-appartient-s01e04.avi",
			"demain-nous-appartient",
			"s01e04",
			0,
		},
		{"slugFile", args{
			"MacGyver.2016.S2E1.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement.mkv",
		},
			"macgyver-s02e01.mkv",
			"macgyver-2016",
			"s02e01",
			2016,
		},
		{"slugFile", args{
			"The.Flash.2014.S04E01.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement.Ws.mkv",
		},
			"the-flash-s04e01.mkv",
			"the-flash-2014",
			"s04e01",
			2014,
		},
		{"slugFile", args{
			"Les.Tuche.French.DVDRIP-zone-telechargement.ws.avi",
		},
			"les-tuche.avi",
			"",
			"",
			0,
		},
		{"slugFile", args{
			"Kingsman.The.Golden.Circle.2017.FRENCH.BDRip.XviD-GZR.WwW.Zone-Telechargement.Ws.avi",
		},
			"kingsman-the-golden-circle.avi",
			"",
			"",
			2017,
		},
		{"slugFile", args{
			"new.girl.episode.101.trois.gars.une.fille-Zone-Telechargement.Ws.avi",
		},
			"new-girl-s00e101.avi",
			"new-girl",
			"s00e101",
			0,
		},
		{"slugFile", args{
			"Thor le monde des tenebres 720p-Shanks@Zone-Telechargement.Ws.mkv",
		},
			"thor-le-monde-des-tenebres.mkv",
			"",
			"",
			0,
		},
		{"slugFile", args{
			"Pitch.Perfect.3.2017.MULTI.1080p.BluRay.x264-VENUE.Zone-Telechargement1.Com.mkv",
		},
			"pitch-perfect-3.mkv",
			"",
			"",
			2017,
		},
		{"slugFile", args{
			"acts of violence 2018-Ww2.zone-telechargement1.com.mkv",
		},
			"acts-of-violence.mkv",
			"",
			"",
			2018,
		},
		{"slugFile", args{
			"Major 2nd - Episode 23 VOSTFR (720p)-Zone-Telechargement1 org.mkv",
		},
			"major-2nd-s00e23.mkv",
			"major-2nd",
			"s00e23",
			0,
		},
		{"slugFile", args{
			"Inaz Ares episode 21 VOSTFR-Zone-Telechargement1 org.mkv",
		},
			"inaz-ares-s00e21.mkv",
			"inaz-ares",
			"s00e21",
			0,
		},
		{"slugFile", args{
			"Inaz Ares episode 22 VOSTFR-Zone-Telechargement1 org.mkv",
		},
			"inaz-ares-s00e22.mkv",
			"inaz-ares",
			"s00e22",
			0,
		},
		{"slugFile", args{
			"The.100.S05E13.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement1.ORG.mkv",
		},
			"the-100-s05e13.mkv",
			"the-100",
			"s05e13",
			0,
		},
		{"slugFile", args{
			"la-nonne.mkv",
		},
			"la-nonne.mkv",
			"",
			"",
			0,
		},
		{"slugFile", args{
			"The.Nun.2018.MULTi.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
		},
			"the-nun.mkv",
			"",
			"",
			2018,
		},
		{"slugFile", args{
			"The.Nun.VOSTFR.2018.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
		},
			"the-nun.mkv",
			"",
			"",
			2018,
		},
		{"slugFile", args{
			"The.Nun.MULTI.2018.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
		},
			"the-nun.mkv",
			"",
			"",
			2018,
		},
		{"slugFile", args{
			"The_Vampire Diaries S04E23 MULTI Ici ou Ailleurs  BluRay720p  2013.mkv",
		},
			"the-vampire-diaries-s04e23.mkv",
			"the-vampire-diaries",
			"s04e23",
			0,
		},
		{"slugFile", args{
			"Project.Blue.Book.S01E01.FRENCH.720p.HDTV.x264-HuSSLe.WwW.Zone-Telechargement.NET.mkv",
		},
			"project-blue-book-s01e01.mkv",
			"project-blue-book",
			"s01e01",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotSerieName, gotSerieNumberReturn, gotYear := slugFile(tt.args.file)
			if gotName != tt.wantName {
				t.Errorf("slugFile() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotSerieName != tt.wantSerieName {
				t.Errorf("slugFile() gotSerieName = %v, want %v", gotSerieName, tt.wantSerieName)
			}
			if gotSerieNumberReturn != tt.wantSerieNumberReturn {
				t.Errorf("slugFile() gotSerieNumberReturn = %v, want %v", gotSerieNumberReturn, tt.wantSerieNumberReturn)
			}
			if gotYear != tt.wantYear {
				t.Errorf("slugFile() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
		})
	}
}

func Test_slugSerieSeasonEpisode(t *testing.T) {
	type args struct {
		serieNumber string
	}
	tests := []struct {
		name                 string
		args                 args
		wantSeasonAndEpisode string
		wantSeason           int
		wantEpisode          int
	}{
		{"slugSerieSeasonEpisode", args{
			"s01e09",
		},
			"s01e09",
			1,
			9,
		},
		{"slugSerieSeasonEpisode", args{
			"109",
		},
			"s00e109",
			0,
			109,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeasonAndEpisode, gotSeason, gotEpisode := slugSerieSeasonEpisode(tt.args.serieNumber)
			if gotSeasonAndEpisode != tt.wantSeasonAndEpisode {
				t.Errorf("slugSerieSeasonEpisode() gotSeasonAndEpisode = %v, want %v", gotSeasonAndEpisode, tt.wantSeasonAndEpisode)
			}
			if gotSeason != tt.wantSeason {
				t.Errorf("slugSerieSeasonEpisode() gotSeason = %v, want %v", gotSeason, tt.wantSeason)
			}
			if gotEpisode != tt.wantEpisode {
				t.Errorf("slugSerieSeasonEpisode() gotEpisode = %v, want %v", gotEpisode, tt.wantEpisode)
			}
		})
	}
}
