package controllers

import (
	"testing"
)

func Test_folderExist(t *testing.T) {
	type args struct {
		folder    string
		serieName string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			"folderExist", args{
				"/home/sam/go/src/searchAndSort/dlna/Series",
				"macgyver-2016",
			},
			"macgyver",
			true,
		},
		{
			"folderExist", args{
				"/home/sam/go/src/searchAndSort/dlna/Series",
				"young-and-hungry",
			},
			"young-and-hungry",
			true,
		},
		{
			"folderExist", args{
				"/home/sam/go/src/searchAndSort/dlna/Series",
				"star-trek-discovery",
			},
			"star-trek-discovery",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := folderExist(tt.args.folder, tt.args.serieName)
			if got != tt.want {
				t.Errorf("folderExist() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("folderExist() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_checkFolderSerie(t *testing.T) {
	type args struct {
		file      string
		name      string
		serieName string
		season    int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			"checkFolderSerie", args{
				"/home/sam/go/src/searchAndSort/dlna/star-trek-discovery-s01e25.mkv",
				"star-trek-discovery-s01e25.mkv",
				"star-trek-discovery",
				1,
			},
			"/home/sam/go/src/searchAndSort/dlna/star-trek-discovery-s01e25.mkv",
			"/home/sam/go/src/searchAndSort/dlna/Series/star-trek-discovery/season-1/star-trek-discovery-s01e25.mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkFolderSerie(tt.args.file, tt.args.name, tt.args.serieName, tt.args.season)
			if got != tt.want {
				t.Errorf("checkFolderSerie() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkFolderSerie() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_start(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"start", args{
				"GNF2.2016.FRENCH.HDRiP.XViD-STVFRV.avi",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start(tt.args.file)
		})
	}
}

func Test_searchSimilarFolder(t *testing.T) {
	type args struct {
		currentPath string
		newFolder   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// {
		// 	"searchSimilarFolder", args{
		// 		"/home/sam/go/src/searchAndSort/dlna/Series",
		// 		"colony",
		// 	},
		// 	"colony",
		// },
		// {
		// 	"searchSimilarFolder", args{
		// 		"/home/sam/go/src/searchAndSort/dlna/Series",
		// 		"macgyver-2016",
		// 	},
		// 	"macgyver",
		// },
		// {
		// 	"searchSimilarFolder", args{
		// 		"/home/sam/go/src/searchAndSort/dlna/Series",
		// 		"young-&-hungry",
		// 	},
		// 	"young-and-hungry",
		// },
		{
			"searchSimilarFolder", args{
				"C:\\Users\\kamel\\go\\src\\search-and-sort-movies\\dlna",
				"star-trek-discovery",
			},
			"star-trek-discovery",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchSimilarFolder(tt.args.currentPath, tt.args.newFolder); got != tt.want {
				t.Errorf("searchSimilarFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}
