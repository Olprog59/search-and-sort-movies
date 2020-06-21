package myapp

import "testing"

//
//import (
//	"testing"
//)
//
//func Test_checkFolderSerie(t *testing.T) {
//	type args struct {
//		file      string
//		name      string
//		serieName string
//		season    int
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  string
//		want1 string
//	}{
//		{
//			"checkFolderSerie", args{
//				"/home/sam/go/src/searchAndSort/dlna/star-trek-discovery-s01e25.mkv",
//				"star-trek-discovery-s01e25.mkv",
//				"star-trek-discovery",
//				1,
//			},
//			"/home/sam/go/src/searchAndSort/dlna/star-trek-discovery-s01e25.mkv",
//			"/home/sam/go/src/searchAndSort/dlna/Series/star-trek-discovery/season-1/star-trek-discovery-s01e25.mkv",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := checkFolderSerie(tt.args.file, tt.args.name, tt.args.serieName, tt.args.season)
//			if got != tt.want {
//				t.Errorf("checkFolderSerie() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("checkFolderSerie() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_start(t *testing.T) {
//	type args struct {
//		file string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			"start", args{
//				"GNF2.2016.FRENCH.HDRiP.XViD-STVFRV.avi",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			start(tt.args.file)
//		})
//	}
//}
//
//func Test_checkIfSizeIsSame(t *testing.T) {
//	type args struct {
//		oldFile string
//		newFile string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			"checlIfSizeIsSame", args{
//				"\\\\SOKYS\\sam\\Series\\dragon-ball-super\\season-1\\dragon-ball-super-episode-77.mp4",
//				"C:\\Users\\kamel\\go\\src\\search-and-sort-movies\\dlna\\dragon-ball-super-episode-77.mp4",
//			},
//			false,
//		},
//		{
//			"checlIfSizeIsSame", args{
//				"\\\\SOKYS\\sam\\Series\\dragon-ball-super\\season-1\\dragon-ball-super-episode-86.mp4",
//				"C:\\Users\\kamel\\go\\src\\search-and-sort-movies\\dlna\\dragon-ball-super-episode-91.mp4",
//			},
//			true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := checkIfSizeIsSame(tt.args.oldFile, tt.args.newFile); (err != nil) != tt.wantErr {
//				t.Errorf("checkIfSizeIsSame() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func Test_copyFile(t *testing.T) {
//	type args struct {
//		oldFile string
//		newFile string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			"copyFile", args{
//				"C:\\Users\\kamel\\go\\src\\search-and-sort-movies\\dlna\\dragon-ball-super-episode-91.mp4",
//				"\\\\SOKYS\\sam\\Series\\dragon-ball-super\\season-1\\dragon-ball-super-episode-86.mp4",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			copyFile(tt.args.oldFile, tt.args.newFile)
//		})
//	}
//}
//
//func Test_moveOrRenameFile(t *testing.T) {
//	type args struct {
//		filePathOld string
//		filePathNew string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			"moveOrRenameFile", args{
//				"/home/sam/go/src/search-and-sort-movies/dlna/Series STAr trek  Discovery-s01e25.mkv",
//				"/home/sam/go/src/search-and-sort-movies/dlna/Series/star-trek-discovery/season-1/star-trek-discovery-s01e25.mkv",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			moveOrRenameFile(tt.args.filePathOld, tt.args.filePathNew)
//		})
//	}
//}

//func Test_myFile_checkFolderSerie(t *testing.T) {
//	type fields struct {
//		file           string
//		fileWithoutDir string
//		complete       string
//		name           string
//		bingName       string
//		transName      string
//		serieName      string
//		serieNumber    string
//		season         string
//		year           int
//		episode        string
//		count          int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//		want1  string
//	}{
//		{
//			"checkFolderSerie", fields{
//				"/home/sam/go/src/searchAndSort/dlna/star-trek-discovery-s01e25.mkv",
//				"star-trek-discovery-s01e25.mkv",
//				"star-trek-discovery-s01e25",
//				"star-trek-discovery",
//				"",
//				"",
//				"star-trek-discovery",
//				"s01e25",
//				"s01",
//				0,
//				"e25",
//				0,
//			},
//			"star-trek-discovery-s01e25",
//			"/star-trek-discovery/season-01/star-trek-discovery-s01e25",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := &myFile{
//				file:           tt.fields.file,
//				fileWithoutDir: tt.fields.fileWithoutDir,
//				complete:       tt.fields.complete,
//				name:           tt.fields.name,
//				bingName:       tt.fields.bingName,
//				transName:      tt.fields.transName,
//				serieName:      tt.fields.serieName,
//				serieNumber:    tt.fields.serieNumber,
//				season:         tt.fields.season,
//				year:           tt.fields.year,
//				episode:        tt.fields.episode,
//				count:          tt.fields.count,
//			}
//			got, got1 := m.checkFolderSerie()
//			if got != tt.want {
//				t.Errorf("checkFolderSerie() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("checkFolderSerie() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}

func Test_myFile_translateName(t *testing.T) {
	type fields struct {
		file           string
		fileWithoutDir string
		complete       string
		name           string
		bingName       string
		transName      string
		serieName      string
		serieNumber    string
		season         string
		year           int
		episode        string
		count          int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"myFile_translateName", fields{name: "demain"}, "tomorrow",
		},
		{
			"myFile_translateName", fields{name: "imorgon"}, "tomorrow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &myFile{
				file:           tt.fields.file,
				fileWithoutDir: tt.fields.fileWithoutDir,
				complete:       tt.fields.complete,
				name:           tt.fields.name,
				bingName:       tt.fields.bingName,
				transName:      tt.fields.transName,
				serieName:      tt.fields.serieName,
				serieNumber:    tt.fields.serieNumber,
				season:         tt.fields.season,
				year:           tt.fields.year,
				episode:        tt.fields.episode,
				count:          tt.fields.count,
			}
			m.translateName()
			if m.transName != tt.want {
				t.Errorf("translateName() got = %v, want %v", m.transName, tt.want)
			}
		})
	}
}
