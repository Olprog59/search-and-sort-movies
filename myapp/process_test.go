package myapp
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
