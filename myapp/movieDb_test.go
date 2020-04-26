package myapp

import (
	"testing"
)

func Test_checkMovieDB(t *testing.T) {
	type args struct {
		tv       bool
		lang     bool
		name     string
		original string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"checkMovieDB",
			args{
				false,
				false,
				"transformers-the-last-knight",
				"transformers-the-last-knight",
			},
			"https://api.themoviedb.org/3/search/movie?api_key=ea8779638f078f25daa3913e80fe46eb&query=transformers-the-last-knight&year=2017",
		},
		{
			"checkMovieDB",
			args{
				true,
				false,
				"the-flash-2014",
				"the-flash-2014",
			},
			"https://api.themoviedb.org/3/search/tv?api_key=ea8779638f078f25daa3913e80fe46eb&query=the-flash&first_air_date_year=2014",
		},
		{
			"checkMovieDB",
			args{
				true,
				false,
				"brooklyn-nine-nine",
				"brooklyn-nine-nine",
			},
			"https://api.themoviedb.org/3/search/tv?api_key=ea8779638f078f25daa3913e80fe46eb&query=brooklyn-nine-nine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkMovieDB(tt.args.tv, tt.args.lang, tt.args.name, tt.args.original); got != tt.want {
				t.Errorf("checkMovieDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slugRemoveYearSerieForSearchMovieDB(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantNew string
	}{
		{
			"slugRemoveYearSerieForSearchMovieDB", args{
				"the-flash-2014",
			},
			"the-flash",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNew := slugRemoveYearSerieForSearchMovieDB(tt.args.name); gotNew != tt.wantNew {
				t.Errorf("slugRemoveYearSerieForSearchMovieDB() = %v, want %v", gotNew, tt.wantNew)
			}
		})
	}
}

func Test_translateName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name          string
		args          args
		wantTranslate string
		wantOriginal  string
	}{
		{
			"translateName", args{
				"28+jours+plus+tards",
			},
			"28 days later",
			"28 jours plus tards",
		},
		{
			"translateName", args{
				"28+semaines+plus+tard",
			},
			"28 weeks later",
			"28 semaines plus tard",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTranslate, gotOriginal := translateName(tt.args.name)
			if gotTranslate != tt.wantTranslate {
				t.Errorf("translateName() gotTranslate = %v, want %v", gotTranslate, tt.wantTranslate)
			}
			if gotOriginal != tt.wantOriginal {
				t.Errorf("translateName() gotOriginal = %v, want %v", gotOriginal, tt.wantOriginal)
			}
		})
	}
}
