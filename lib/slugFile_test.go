package lib

import (
	"testing"
)

func Test_myFile_slugFile(t *testing.T) {
	type fields struct {
		complete string
	}
	tests := []struct {
		name            string
		fields          fields
		wantComplete    string
		wantName        string
		wantSerieName   string
		wantSerieNumber string
		wantSeason      int
		wantEpisode     int
		wantYear        int
		wantLanguage    string
	}{
		{
			"slugFile", fields{
				"Charmed.2018.S02E11.FRENCH.720p.AMZN.WEB-DL.DD5.1.H264-FRATERNiTY.mkv",
			},
			"charmed-s02e11-720p (2018).mkv",
			"charmed-s02e11-720p (2018)",
			"charmed",
			"s02e11",
			02,
			11,
			2018,
			"french",
		},
		{
			"slugFile", fields{
				"Motherland.Fort.Salem.S01E09.FRENCH.720p.WEB.DDP5.1.H264-FRATERNiTY.mkv",
			},
			"motherland-fort-salem-s01e09-720p.mkv",
			"motherland-fort-salem-s01e09-720p",
			"motherland-fort-salem",
			"s01e09",
			1,
			9,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"Salvation.S02E04.FRENCH.720p.AMZN.WEB-DL.DD5.1.H264-FRATERNiTY-ZT.WwW.Zone-Telechargement.NET.mkv",
			},
			"salvation-s02e04-720p.mkv",
			"salvation-s02e04-720p",
			"salvation",
			"s02e04",
			2,
			4,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"The.100.S05E13.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement1.ORG.mkv",
			},
			"the-100-s05e13-720p.mkv",
			"the-100-s05e13-720p",
			"the-100",
			"s05e13",
			05,
			13,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"Sekai-Raws radiant 2 - 09 VOSTFR CR 720p-Zone-Telechargement.NET.mp4",
			},
			"sekai-raws-radiant-s02e09-720p.mp4",
			"sekai-raws-radiant-s02e09-720p",
			"sekai-raws-radiant",
			"s02e09",
			2,
			9,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"Sekai-Raws radiant 2 - 12 VOSTFR CR 720p-Zone-Telechargement.NET .mp4",
			},
			"sekai-raws-radiant-s02e12-720p.mp4",
			"sekai-raws-radiant-s02e12-720p",
			"sekai-raws-radiant",
			"s02e12",
			02,
			12,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"Sekai-Raws radiant 2 - 14 VOSTFR CR 720p-www.zone-annuaire.com.mp4",
			},
			"sekai-raws-radiant-s02e14-720p.mp4",
			"sekai-raws-radiant-s02e14-720p",
			"sekai-raws-radiant",
			"s02e14",
			02,
			14,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"radiant - Saison 2 Episode 8 VOSTFR-Zone-Telechargement.NET.mp4",
			},
			"radiant-s02e08.mp4",
			"radiant-s02e08",
			"radiant",
			"s02e08",
			2,
			8,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"The.100.S05E13.FASTSUB.VOSTFR.720p.HDTV.x264-ZT.WwW.Zone-Telechargement1.ORG.mkv",
			},
			"the-100-s05e13-720p.mkv",
			"the-100-s05e13-720p",
			"the-100",
			"s05e13",
			05,
			13,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"9-1-1.S03E01.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
			},
			"9-1-1-s03e01-720p.mkv",
			"9-1-1-s03e01-720p",
			"9-1-1",
			"s03e01",
			3,
			1,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"9-1-1.S03E02.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
			},
			"9-1-1-s03e02-720p.mkv",
			"9-1-1-s03e02-720p",
			"9-1-1",
			"s03e02",
			3,
			2,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"9-1-1.S03E03.FRENCH.720p.WEB.H264-CiELOS.mkv",
			},
			"9-1-1-s03e03-720p.mkv",
			"9-1-1-s03e03-720p",
			"9-1-1",
			"s03e03",
			3,
			3,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"9-1-1.S03E04.FRENCH.720p.WEB.H264-CiELOS.WwW.Zone-Annuaire.COM.mkv",
			},
			"9-1-1-s03e04-720p.mkv",
			"9-1-1-s03e04-720p",
			"9-1-1",
			"s03e04",
			3,
			4,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"9-1-1.S03E05.FRENCH.720p.WEB.H264-CiELOS.WwW.Zone-Annuaire.COM.mkv",
			},
			"9-1-1-s03e05-720p.mkv",
			"9-1-1-s03e05-720p",
			"9-1-1",
			"s03e05",
			3,
			5,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"9-1-1.S03E07.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
			},
			"9-1-1-s03e07-720p.mkv",
			"9-1-1-s03e07-720p",
			"9-1-1",
			"s03e07",
			3,
			7,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"9-1-1.S03E08.FRENCH.720p.HDTV.x264-SH0W.WwW.Zone-Annuaire.COM.mkv",
			},
			"9-1-1-s03e08-720p.mkv",
			"9-1-1-s03e08-720p",
			"9-1-1",
			"s03e08",
			3,
			8,
			0,
			"french",
		},
		{
			"slugFile", fields{
				"Ahiru.no.Sora.E28.VOSTFR.x264--ZONE-ANNUAIRE.COM.mp4",
			},
			"ahiru-no-sora-s00e28.mp4",
			"ahiru-no-sora-s00e28",
			"ahiru-no-sora",
			"s00e28",
			0,
			28,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"Boruto e81 Vostfr.mp4",
			},
			"boruto-s00e81.mp4",
			"boruto-s00e81",
			"boruto",
			"s00e81",
			0,
			81,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"Ducobu.3.2020.FRENCH.1080p.WEB.H264-PREUMS_wWw.Extreme-Down.Ninja.mkv",
			},
			"ducobu-3-1080p (2020).mkv",
			"ducobu-3-1080p (2020)",
			"",
			"",
			0,
			0,
			2020,
			"french",
		},
		{
			"slugFile", fields{
				"DC.League.of.Super.Pets.2022.4K.MULTi.2160p.HDR.WEB.EAC3.x265-EXTREME_wWw.Extreme-Down.io.mkv",
			},
			"dc-league-of-super-pets-4k (2022).mkv",
			"dc-league-of-super-pets-4k (2022)",
			"",
			"",
			0,
			0,
			2022,
			"multi",
		},
		{
			"slugFile", fields{
				"dc-league-of-super-pets-2022.mkv",
			},
			"dc-league-of-super-pets (2022).mkv",
			"dc-league-of-super-pets (2022)",
			"",
			"",
			0,
			0,
			2022,
			"",
		},
		{
			"slugFile", fields{
				"Fairy Tail episode 001 MULTI ''Fairy Tail'' ... BluRay1080p ! {Chris44} (Saison 01) ''Arc Macao'' 2009.mkv",
			},
			"fairy-tail-s00e01-1080p (2009).mkv",
			"fairy-tail-s00e01-1080p (2009)",
			"fairy-tail",
			"s00e01",
			0,
			01,
			2009,
			"multi",
		},
		{
			"slugFile", fields{
				"La.Brea.S01E010.FiNAL.MULTi.1080p.AMZN.WEB-DL.H264-TiNA.mkv",
			},
			"la-brea-s01e10-1080p.mkv",
			"la-brea-s01e10-1080p",
			"la-brea",
			"s01e10",
			01,
			10,
			0,
			"multi",
		},
		{
			"slugFile", fields{
				"One Piece 1000 MULTI ''Puissance hors du commun ! L'équipage au Chapeau de paille au complet'' ... WebDl1080p ! 2021.mkv",
			},
			"one-piece-1000-1080p (2021).mkv",
			"one-piece-1000-1080p (2021)",
			"",
			"",
			0,
			0,
			2021,
			"multi",
		},
		{
			"slugFile", fields{
				"One.Piece.S01E1059.SUBFRENCH.1080p.WEB.x264.AAC-Tsundere-Raws.mkv",
			},
			"one-piece-s01e1059-1080p.mkv",
			"one-piece-s01e1059-1080p",
			"one-piece",
			"s01e1059",
			01,
			1059,
			0,
			"subfrench",
		},
		{
			"slugFile", fields{
				"Jujutsu Kaisen s02e17.mp4",
			},
			"jujutsu-kaisen-s02e17.mp4",
			"jujutsu-kaisen-s02e17",
			"jujutsu-kaisen",
			"s02e17",
			02,
			17,
			0,
			"",
		},
		{
			"slugFile", fields{
				"Jujutsu Kaisen s02e18.mp4",
			},
			"jujutsu-kaisen-s02e18.mp4",
			"jujutsu-kaisen-s02e18",
			"jujutsu-kaisen",
			"s02e18",
			02,
			18,
			0,
			"",
		},
		{
			"slugFile", fields{
				"World.War.Z.2013.MULTi.1080p.AMZN.WEB.DDP5.1.H265-TFA.mkv",
			},
			"world-war-z-1080p (2013).mkv",
			"world-war-z-1080p (2013)",
			"",
			"",
			0,
			0,
			2013,
			"multi",
		},
		{
			"slugFile", fields{
				"[Kaerizaki-Fansub] One Piece s21e1092 VOSTFR FHD (1920x1080) .mp4",
			},
			"one-piece-s21e1092-1080p.mp4",
			"one-piece-s21e1092-1080p",
			"one-piece",
			"s21e1092",
			21,
			1092,
			0,
			"vostfr",
		},
		{
			"slugFile", fields{
				"Stargate Universe - 2X02 - Aftermath.mkv",
			},
			"stargate-universe-s02e02.mkv",
			"stargate-universe-s02e02",
			"stargate-universe",
			"s02e02",
			2,
			2,
			0,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &myFile{
				complete: tt.fields.complete,
			}
			err := m.slugFile()
			if err != nil {
				return
			}
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
			if tt.wantYear != m.year {
				t.Errorf("wantYear : %v, want : %v", m.year, tt.wantYear)
			}
			if tt.wantLanguage != m.language {
				t.Errorf("wantLanguage : %v, want : %v", m.language, tt.wantLanguage)
			}
		})
	}
}

func Test_myFile_extractResolution(t *testing.T) {
	type fields struct {
		file           string
		ext            string
		resolution     string
		fileWithoutDir string
		complete       string
		completeSlug   string
		name           string
		serieName      string
		serieNumber    string
		season         int
		year           int
		episode        int
		episodeRaw     string
		count          int
		language       string
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &myFile{
				file:           tt.fields.file,
				ext:            tt.fields.ext,
				resolution:     tt.fields.resolution,
				fileWithoutDir: tt.fields.fileWithoutDir,
				complete:       tt.fields.complete,
				completeSlug:   tt.fields.completeSlug,
				name:           tt.fields.name,
				serieName:      tt.fields.serieName,
				serieNumber:    tt.fields.serieNumber,
				season:         tt.fields.season,
				year:           tt.fields.year,
				episode:        tt.fields.episode,
				episodeRaw:     tt.fields.episodeRaw,
				count:          tt.fields.count,
				language:       tt.fields.language,
			}
			m.extractResolution()
		})
	}
}

func Test_myFile_extractYear(t *testing.T) {
	type fields struct {
		file           string
		ext            string
		resolution     string
		fileWithoutDir string
		complete       string
		completeSlug   string
		name           string
		serieName      string
		serieNumber    string
		season         int
		year           int
		episode        int
		episodeRaw     string
		count          int
		language       string
	}
	type args struct {
		str string
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &myFile{
				file:           tt.fields.file,
				ext:            tt.fields.ext,
				resolution:     tt.fields.resolution,
				fileWithoutDir: tt.fields.fileWithoutDir,
				complete:       tt.fields.complete,
				completeSlug:   tt.fields.completeSlug,
				name:           tt.fields.name,
				serieName:      tt.fields.serieName,
				serieNumber:    tt.fields.serieNumber,
				season:         tt.fields.season,
				year:           tt.fields.year,
				episode:        tt.fields.episode,
				episodeRaw:     tt.fields.episodeRaw,
				count:          tt.fields.count,
				language:       tt.fields.language,
			}
			m.extractYear(tt.args.str)
		})
	}
}

func Test_myFile_formatageFinal(t *testing.T) {
	type fields struct {
		file           string
		ext            string
		resolution     string
		fileWithoutDir string
		complete       string
		completeSlug   string
		name           string
		serieName      string
		serieNumber    string
		season         int
		year           int
		episode        int
		episodeRaw     string
		count          int
		language       string
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &myFile{
				file:           tt.fields.file,
				ext:            tt.fields.ext,
				resolution:     tt.fields.resolution,
				fileWithoutDir: tt.fields.fileWithoutDir,
				complete:       tt.fields.complete,
				completeSlug:   tt.fields.completeSlug,
				name:           tt.fields.name,
				serieName:      tt.fields.serieName,
				serieNumber:    tt.fields.serieNumber,
				season:         tt.fields.season,
				year:           tt.fields.year,
				episode:        tt.fields.episode,
				episodeRaw:     tt.fields.episodeRaw,
				count:          tt.fields.count,
				language:       tt.fields.language,
			}
			err := m.formatageFinal()
			if err != nil {
				return
			}
		})
	}
}

func Test_myFile_slugFile1(t *testing.T) {
	type fields struct {
		file           string
		ext            string
		resolution     string
		fileWithoutDir string
		complete       string
		completeSlug   string
		name           string
		serieName      string
		serieNumber    string
		season         int
		year           int
		episode        int
		episodeRaw     string
		count          int
		language       string
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &myFile{
				file:           tt.fields.file,
				ext:            tt.fields.ext,
				resolution:     tt.fields.resolution,
				fileWithoutDir: tt.fields.fileWithoutDir,
				complete:       tt.fields.complete,
				completeSlug:   tt.fields.completeSlug,
				name:           tt.fields.name,
				serieName:      tt.fields.serieName,
				serieNumber:    tt.fields.serieNumber,
				season:         tt.fields.season,
				year:           tt.fields.year,
				episode:        tt.fields.episode,
				episodeRaw:     tt.fields.episodeRaw,
				count:          tt.fields.count,
				language:       tt.fields.language,
			}
			err := m.slugFile()
			if err != nil {
				return
			}
		})
	}
}
