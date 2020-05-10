package myapp

import (
	"testing"
)

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
		{
			"slugFile", fields{
				"",
				"Boruto e81 Vostfr.mp4",
				"",
				"",
				"",
				"",
				"",
			},
			"boruto-s00e81.mp4",
			"boruto-s00e81",
			"boruto",
			"s00e81",
			"s00",
			"e81",
		},
		{
			"slugFile", fields{
				"",
				"Ducobu.3.2020.FRENCH.1080p.WEB.H264-PREUMS_wWw.Extreme-Down.Ninja.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"ducobu-3.mkv",
			"ducobu-3",
			"",
			"",
			"",
			"",
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
