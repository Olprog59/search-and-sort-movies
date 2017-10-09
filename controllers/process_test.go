package controllers

import "testing"

func Test_checkFolderSerie(t *testing.T) {
	type args struct {
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
			"checFolderSerie",
			args{
				"macgyver-s02e01.mkv",
				"macgyver-2016",
				02,
			},
			"/macgyver-s02e01.mkv",
			"/macgyver-2016/season-2/macgyver-s02e01.mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkFolderSerie(tt.args.name, tt.args.serieName, tt.args.season)
			if got != tt.want {
				t.Errorf("checkFolderSerie() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkFolderSerie() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
