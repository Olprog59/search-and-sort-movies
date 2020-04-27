package myapp

import (
	"testing"
)

//func Test_slugFile2(t *testing.T) {
//	type args struct {
//		file string
//	}
//	tests := []struct {
//		name                  string
//		args                  args
//		wantName              string
//		wantSerieName         string
//		wantSerieNumberReturn string
//		wantYear              int
//	}{
//
//		{"slugFile", args{
//			"Pirates.of.the.Caribbean.Dead.Men.Tell.No.Tales.2037.MULTi.1080p.WEB.H264-SPACED.WwW.Zone-Telechargement.Ws.mkv",
//		},
//			"pirates-of-the-caribbean-dead-men-tell-no-tales.mkv",
//			"",
//			"",
//			2017,
//		},
//		{"slugFile", args{
//			"modern.family.s1e09.french.dvdrip.xvid-jmt-Zone-telechargement.Ws.avi",
//		},
//			"modern-family-s01e09.avi",
//			"modern-family",
//			"s01e09",
//			0,
//		},
//		{"slugFile", args{
//			"Snatched.2016.MULTi.1080p.BluRay.x264-VENUE.WwW.Zone-Telechargement.Ws.mkv",
//		},
//			"snatched.mkv",
//			"",
//			"",
//			2016,
//		},
//		{"slugFile", args{
//			"Demain.Nous.Appartient.S01E4.FRENCH.HDTV.XviD-ZT.WwW.Zone-Telechargement.Ws.avi",
//		},
//			"demain-nous-appartient-s01e04.avi",
//			"demain-nous-appartient",
//			"s01e04",
//			0,
//		},
//		{"slugFile", args{
//			"MacGyver.2016.S2E1.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement.mkv",
//		},
//			"macgyver-s02e01.mkv",
//			"macgyver-2016",
//			"s02e01",
//			2016,
//		},
//		{"slugFile", args{
//			"The.Flash.2014.S04E01.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement.Ws.mkv",
//		},
//			"the-flash-s04e01.mkv",
//			"the-flash-2014",
//			"s04e01",
//			2014,
//		},
//		{"slugFile", args{
//			"Les.Tuche.French.DVDRIP-zone-telechargement.ws.avi",
//		},
//			"les-tuche.avi",
//			"",
//			"",
//			0,
//		},
//		{"slugFile", args{
//			"Kingsman.The.Golden.Circle.2017.FRENCH.BDRip.XviD-GZR.WwW.Zone-Telechargement.Ws.avi",
//		},
//			"kingsman-the-golden-circle.avi",
//			"",
//			"",
//			2017,
//		},
//		{"slugFile", args{
//			"new.girl.episode.101.trois.gars.une.fille-Zone-Telechargement.Ws.avi",
//		},
//			"new-girl-s00e101.avi",
//			"new-girl",
//			"s00e101",
//			0,
//		},
//		{"slugFile", args{
//			"Thor le monde des tenebres 720p-Shanks@Zone-Telechargement.Ws.mkv",
//		},
//			"thor-le-monde-des-tenebres.mkv",
//			"",
//			"",
//			0,
//		},
//		{"slugFile", args{
//			"Pitch.Perfect.3.2017.MULTI.1080p.BluRay.x264-VENUE.Zone-Telechargement1.Com.mkv",
//		},
//			"pitch-perfect-3.mkv",
//			"",
//			"",
//			2017,
//		},
//		{"slugFile", args{
//			"acts of violence 2018-Ww2.zone-telechargement1.com.mkv",
//		},
//			"acts-of-violence.mkv",
//			"",
//			"",
//			2018,
//		},
//		{"slugFile", args{
//			"Major 2nd - Episode 23 VOSTFR (720p)-Zone-Telechargement1 org.mkv",
//		},
//			"major-2nd-s00e23.mkv",
//			"major-2nd",
//			"s00e23",
//			0,
//		},
//		{"slugFile", args{
//			"Inaz Ares episode 21 VOSTFR-Zone-Telechargement1 org.mkv",
//		},
//			"inaz-ares-s00e21.mkv",
//			"inaz-ares",
//			"s00e21",
//			0,
//		},
//		{"slugFile", args{
//			"Inaz Ares episode 22 VOSTFR-Zone-Telechargement1 org.mkv",
//		},
//			"inaz-ares-s00e22.mkv",
//			"inaz-ares",
//			"s00e22",
//			0,
//		},
//		{"slugFile", args{
//			"The.100.S05E13.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement1.ORG.mkv",
//		},
//			"the-100-s05e13.mkv",
//			"the-100",
//			"s05e13",
//			0,
//		},
//		{"slugFile", args{
//			"la-nonne.mkv",
//		},
//			"la-nonne.mkv",
//			"",
//			"",
//			0,
//		},
//		{"slugFile", args{
//			"The.Nun.2018.MULTi.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
//		},
//			"the-nun.mkv",
//			"",
//			"",
//			2018,
//		},
//		{"slugFile", args{
//			"The.Nun.VOSTFR.2018.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
//		},
//			"the-nun.mkv",
//			"",
//			"",
//			2018,
//		},
//		{"slugFile", args{
//			"The.Nun.MULTI.2018.TRUEFRENCH.1080p.HDLight.x264-RDH.WwW.Annuaire-Telechargement.CoM.mkv",
//		},
//			"the-nun.mkv",
//			"",
//			"",
//			2018,
//		},
//		{"slugFile", args{
//			"The_Vampire Diaries S04E23 MULTI Ici ou Ailleurs  BluRay720p  2013.mkv",
//		},
//			"the-vampire-diaries-s04e23.mkv",
//			"the-vampire-diaries",
//			"s04e23",
//			0,
//		},
//		{"slugFile", args{
//			"Project.Blue.Book.S01E01.FRENCH.720p.HDTV.x264-HuSSLe.WwW.Zone-Telechargement.NET.mkv",
//		},
//			"project-blue-book-s01e01.mkv",
//			"project-blue-book",
//			"s01e01",
//			0,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var m *myFile
//			m.slugFile()
//			if gotName != tt.wantName {
//				t.Errorf("slugFile() gotName = %v, want %v", gotName, tt.wantName)
//			}
//			if gotSerieName != tt.wantSerieName {
//				t.Errorf("slugFile() gotSerieName = %v, want %v", gotSerieName, tt.wantSerieName)
//			}
//			if gotSerieNumberReturn != tt.wantSerieNumberReturn {
//				t.Errorf("slugFile() gotSerieNumberReturn = %v, want %v", gotSerieNumberReturn, tt.wantSerieNumberReturn)
//			}
//			if gotYear != tt.wantYear {
//				t.Errorf("slugFile() gotYear = %v, want %v", gotYear, tt.wantYear)
//			}
//		})
//	}
//}

//func Test_slugSerieSeasonEpisode(t *testing.T) {
//	var m *myFile
//	type args struct {
//		serieNumber string
//	}
//	tests := []struct {
//		name                 string
//		args                 args
//		wantSeasonAndEpisode string
//		wantSeason           int
//		wantEpisode          int
//	}{
//		{"slugSerieSeasonEpisode", args{
//			"s01e09",
//		},
//			"s01e09",
//			1,
//			9,
//		},
//		{"slugSerieSeasonEpisode", args{
//			"109",
//		},
//			"s00e109",
//			0,
//			109,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m.slugSerieSeasonEpisode()
//			if m.serieNumber != tt.wantSeasonAndEpisode {
//				t.Errorf("slugSerieSeasonEpisode() gotSeasonAndEpisode = %v, want %v", m.serieNumber, tt.wantSeasonAndEpisode)
//			}
//			if m.season != tt.wantSeason {
//				t.Errorf("slugSerieSeasonEpisode() gotSeason = %v, want %v", m.season, tt.wantSeason)
//			}
//			if m.episode != tt.wantEpisode {
//				t.Errorf("slugSerieSeasonEpisode() gotEpisode = %v, want %v", m.episode, tt.wantEpisode)
//			}
//		})
//	}
//}

func Test_myFile_slugFile(t *testing.T) {
	type fields struct {
		file        string
		complete    string
		name        string
		serieName   string
		serieNumber string
		season      string
		episode     string
	}
	tests := []struct {
		name            string
		fields          fields
		wantComplete    string
		wantName        string
		wantSerieName   string
		wantSerieNumber string
		wantSeason      string
		wantEpisode     string
	}{
		{
			"slugFile", fields{
				"",
				"The.100.S05E13.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement1.ORG.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"the-100-s05e13.mkv",
			"the-100-s05e13",
			"the-100",
			"s05e13",
			"s05",
			"e13",
		},
		{
			"slugFile", fields{
				"",
				"Sekai-Raws radiant 2 - 09 VOSTFR CR 720p-Zone-Telechargement.NET.mp4",
				"",
				"",
				"",
				"",
				"",
			},
			"sekai-raws-radiant-s02e09.mp4",
			"sekai-raws-radiant-s02e09",
			"sekai-raws-radiant",
			"s02e09",
			"s02",
			"e09",
		},
		{
			"slugFile", fields{
				"",
				"Sekai-Raws radiant 2 - 12 VOSTFR CR 720p-Zone-Telechargement.NET .mp4",
				"",
				"",
				"",
				"",
				"",
			},
			"sekai-raws-radiant-s02e12.mp4",
			"sekai-raws-radiant-s02e12",
			"sekai-raws-radiant",
			"s02e12",
			"s02",
			"e12",
		},
		{
			"slugFile", fields{
				"",
				"Sekai-Raws radiant 2 - 14 VOSTFR CR 720p-www.zone-annuaire.com.mp4",
				"",
				"",
				"",
				"",
				"",
			},
			"sekai-raws-radiant-s02e14.mp4",
			"sekai-raws-radiant-s02e14",
			"sekai-raws-radiant",
			"s02e14",
			"s02",
			"e14",
		},
		{
			"slugFile", fields{
				"",
				"radiant - Saison 2 Episode 8 VOSTFR-Zone-Telechargement.NET.mp4",
				"",
				"",
				"",
				"",
				"",
			},
			"radiant-s02e08.mp4",
			"radiant-s02e08",
			"radiant",
			"s02e08",
			"s02",
			"e08",
		},
		{
			"slugFile", fields{
				"",
				"The.100.S05E13.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement1.ORG.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"the-100-s05e13.mkv",
			"the-100-s05e13",
			"the-100",
			"s05e13",
			"s05",
			"e13",
		},
		{
			"slugFile", fields{
				"",
				"9-1-1.S03E01.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"9-1-1-s03e01.mkv",
			"9-1-1-s03e01",
			"9-1-1",
			"s03e01",
			"s03",
			"e01",
		},
		{
			"slugFile", fields{
				"",
				"9-1-1.S03E02.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"9-1-1-s03e02.mkv",
			"9-1-1-s03e02",
			"9-1-1",
			"s03e02",
			"s03",
			"e02",
		},
		{
			"slugFile", fields{
				"",
				"9-1-1.S03E03.FRENCH.720p.WEB.H264-CiELOS.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"9-1-1-s03e03.mkv",
			"9-1-1-s03e03",
			"9-1-1",
			"s03e03",
			"s03",
			"e03",
		},
		{
			"slugFile", fields{
				"",
				"9-1-1.S03E04.FRENCH.720p.WEB.H264-CiELOS.WwW.Zone-Annuaire.COM.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"9-1-1-s03e04.mkv",
			"9-1-1-s03e04",
			"9-1-1",
			"s03e04",
			"s03",
			"e04",
		},
		{
			"slugFile", fields{
				"",
				"9-1-1.S03E05.FRENCH.720p.WEB.H264-CiELOS.WwW.Zone-Annuaire.COM.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"9-1-1-s03e05.mkv",
			"9-1-1-s03e05",
			"9-1-1",
			"s03e05",
			"s03",
			"e05",
		},
		{
			"slugFile", fields{
				"",
				"9-1-1.S03E07.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"9-1-1-s03e07.mkv",
			"9-1-1-s03e07",
			"9-1-1",
			"s03e07",
			"s03",
			"e07",
		},
		{
			"slugFile", fields{
				"",
				"9-1-1.S03E08.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"9-1-1-s03e08.mkv",
			"9-1-1-s03e08",
			"9-1-1",
			"s03e08",
			"s03",
			"e08",
		},
		{
			"slugFile", fields{
				"",
				"Ahiru.no.Sora.E28.VOSTFR.x264--ZONE-ANNUAIRE.COM.mp4",
				"",
				"",
				"",
				"",
				"",
			},
			"ahiru-no-sora-s00e28.mp4",
			"ahiru-no-sora-s00e28",
			"ahiru-no-sora",
			"s00e28",
			"s00",
			"e28",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &myFile{
				file:        tt.fields.file,
				complete:    tt.fields.complete,
				name:        tt.fields.name,
				serieName:   tt.fields.serieName,
				serieNumber: tt.fields.serieNumber,
				season:      tt.fields.season,
				episode:     tt.fields.episode,
			}
			m.slugFile()
			if tt.wantComplete != m.complete {
				t.Errorf("wantComplete : %v, want : %v", m.complete, tt.wantComplete)
			}
			if tt.wantName != m.name {
				t.Errorf("wantName : %v, want : %v", m.name, tt.wantName)
			}
			if tt.wantSerieName != m.serieName {
				t.Errorf("wantSerieName : %v, want : %v", m.serieName, tt.wantSerieName)
			}
			if tt.wantSerieNumber != m.serieNumber {
				t.Errorf("wantSerieNumber : %v, want : %v", m.serieNumber, tt.wantSerieNumber)
			}
			if tt.wantSeason != m.season {
				t.Errorf("wantSeason : %v, want : %v", m.season, tt.wantSeason)
			}
			if tt.wantEpisode != m.episode {
				t.Errorf("wantEpisode : %v, want : %v", m.episode, tt.wantEpisode)
			}
		})
	}
}
