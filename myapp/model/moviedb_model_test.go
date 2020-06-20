package model

import "testing"

func TestGetTrailer(t *testing.T) {
	type args struct {
		id    int64
		serie bool
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			"GetTrailer", args{
				347201,
				false,
			},
			"Qyonn5Vbg7s",
			"YouTube",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetTrailer(tt.args.id, tt.args.serie)
			if got != tt.want {
				t.Errorf("GetTrailer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetTrailer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetImage(t *testing.T) {
	type args struct {
		movie string
		serie bool
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int64
	}{
		{
			"GetImage", args{
				"naruto",
				false,
			},
			"https://image.tmdb.org/t/p/w500/1k6iwC4KaPvTBt1JuaqXy3noZRY.jpg",
			347201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetImage(tt.args.movie, tt.args.serie)
			if got != tt.want {
				t.Errorf("GetImage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetImage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
