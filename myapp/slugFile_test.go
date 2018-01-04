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
			"modern.family.s01e09.french.dvdrip.xvid-jmt-Zone-telechargement.Ws.avi",
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
			"Demain.Nous.Appartient.S01E04.FRENCH.HDTV.XviD-ZT.WwW.Zone-Telechargement.Ws.avi",
		},
			"demain-nous-appartient-s01e04.avi",
			"demain-nous-appartient",
			"s01e04",
			0,
		},
		{"slugFile", args{
			"MacGyver.2016.S02E01.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement.mkv",
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
		name        string
		args        args
		wantSeason  int
		wantEpisode int
	}{
		{"slugSerieSeasonEpisode", args{
			"s01e09",
		},
			1,
			9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeason, gotEpisode := slugSerieSeasonEpisode(tt.args.serieNumber)
			if gotSeason != tt.wantSeason {
				t.Errorf("slugSerieSeasonEpisode() gotSeason = %v, want %v", gotSeason, tt.wantSeason)
			}
			if gotEpisode != tt.wantEpisode {
				t.Errorf("slugSerieSeasonEpisode() gotEpisode = %v, want %v", gotEpisode, tt.wantEpisode)
			}
		})
	}
}
