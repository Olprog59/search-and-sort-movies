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
				"Charmed.2018.S02E11.FRENCH.720p.AMZN.WEB-DL.DD5.1.H264-FRATERNiTY.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"charmed-s02e11-720p-2018.mkv",
			"charmed-s02e11",
			"charmed",
			"s02e11",
			"s02",
			"e11",
		},
		{
			"slugFile", fields{
				"",
				"Motherland.Fort.Salem.S01E09.FRENCH.720p.WEB.DDP5.1.H264-FRATERNiTY.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"motherland-fort-salem-s01e09-720p.mkv",
			"motherland-fort-salem-s01e09",
			"motherland-fort-salem",
			"s01e09",
			"s01",
			"e09",
		},
		{
			"slugFile", fields{
				"",
				"Salvation.S02E04.FRENCH.720p.AMZN.WEB-DL.DD5.1.H264-FRATERNiTY-ZT.WwW.Zone-Telechargement.NET.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"salvation-s02e04-720p.mkv",
			"salvation-s02e04",
			"salvation",
			"s02e04",
			"s02",
			"e04",
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
			"the-100-s05e13-720p.mkv",
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
			"sekai-raws-radiant-s02e09-720p.mp4",
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
			"sekai-raws-radiant-s02e12-720p.mp4",
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
			"sekai-raws-radiant-s02e14-720p.mp4",
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
			"the-100-s05e13-720p.mkv",
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
			"9-1-1-s03e01-720p.mkv",
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
			"9-1-1-s03e02-720p.mkv",
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
			"9-1-1-s03e03-720p.mkv",
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
			"9-1-1-s03e04-720p.mkv",
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
			"9-1-1-s03e05-720p.mkv",
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
			"9-1-1-s03e07-720p.mkv",
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
			"9-1-1-s03e08-720p.mkv",
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
			"ducobu-3-1080p-2020.mkv",
			"ducobu-3",
			"",
			"",
			"",
			"",
		},
		{
			"slugFile", fields{
				"",
				"DC.League.of.Super.Pets.2022.4K.MULTi.2160p.HDR.WEB.EAC3.x265-EXTREME_wWw.Extreme-Down.io.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"dc-league-of-super-pets-4k-2022.mkv",
			"dc-league-of-super-pets",
			"",
			"",
			"",
			"",
		},
		{
			"slugFile", fields{
				"",
				"dc-league-of-super-pets-2022.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"dc-league-of-super-pets-2022.mkv",
			"dc-league-of-super-pets",
			"",
			"",
			"",
			"",
		},
		{
			"slugFile", fields{
				"",
				"Fairy Tail episode 001 MULTI ''Fairy Tail'' ... BluRay1080p ! {Chris44} (Saison 01) ''Arc Macao'' 2009.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"fairy-tail-s00e001-1080p-2009.mkv",
			"fairy-tail-s00e001",
			"fairy-tail",
			"s00e001",
			"s00",
			"e001",
		},
		{
			"slugFile", fields{
				"",
				"La.Brea.S01E010.FiNAL.MULTi.1080p.AMZN.WEB-DL.H264-TiNA.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"la-brea-s01e10-1080p.mkv",
			"la-brea-s01e10",
			"la-brea",
			"s01e10",
			"s01",
			"e10",
		},
		{
			"slugFile", fields{
				"",
				"One Piece 1000 MULTI ''Puissance hors du commun ! L'équipage au Chapeau de paille au complet'' ... WebDl1080p ! 2021.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"one-piece-1000-1080p-2021.mkv",
			"one-piece-1000",
			"",
			"",
			"",
			"",
		},
		{
			"slugFile", fields{
				"",
				"One.Piece.S01E1059.SUBFRENCH.1080p.WEB.x264.AAC-Tsundere-Raws.mkv",
				"",
				"",
				"",
				"",
				"",
			},
			"one-piece-s01e1059-1080p.mkv",
			"one-piece-s01e1059",
			"one-piece",
			"s01e1059",
			"s01",
			"e1059",
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
