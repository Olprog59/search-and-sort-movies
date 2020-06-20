package myapp

import (
	"reflect"
	"testing"
)

//func Test_checkMovieDB(t *testing.T) {
//	type args struct {
//		tv       bool
//		lang     bool
//		name     string
//		original string
//	}
//	tests := []struct {
//		name string
//		args args
//		want string
//	}{
//		{
//			"checkMovieDB",
//			args{
//				false,
//				false,
//				"transformers-the-last-knight",
//				"transformers-the-last-knight",
//			},
//			"https://api.themoviedb.org/3/search/movie?api_key=ea8779638f078f25daa3913e80fe46eb&query=transformers-the-last-knight&year=2017",
//		},
//		{
//			"checkMovieDB",
//			args{
//				true,
//				false,
//				"the-flash-2014",
//				"the-flash-2014",
//			},
//			"https://api.themoviedb.org/3/search/tv?api_key=ea8779638f078f25daa3913e80fe46eb&query=the-flash&first_air_date_year=2014",
//		},
//		{
//			"checkMovieDB",
//			args{
//				true,
//				false,
//				"brooklyn-nine-nine",
//				"brooklyn-nine-nine",
//			},
//			"https://api.themoviedb.org/3/search/tv?api_key=ea8779638f078f25daa3913e80fe46eb&query=brooklyn-nine-nine",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := m.checkMovieDB(tt.args.tv, tt.args.lang, tt.args.name, tt.args.original); got != tt.want {
//				t.Errorf("checkMovieDB() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

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

// Test Ok juste je n'ai pas rempli moviedb
func Test_readJSONFromUrlTV(t *testing.T) {
	var moviedb MoviesDb
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    MoviesDb
		wantErr bool
	}{
		{
			"readJSONFromUrlTV", args{
				"https://api.themoviedb.org/3/search/multi?api_key=ea8779638f078f25daa3913e80fe46eb&query=naruto",
			},
			moviedb,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readJSONFromUrlTV(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("readJSONFromUrlTV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readJSONFromUrlTV() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readJSONFromUrlMovie(t *testing.T) {
	var moviedb MoviesDb
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    MoviesDb
		wantErr bool
	}{
		{
			"readJSONFromUrlMovie", args{
				"https://api.themoviedb.org/3/search/multi?api_key=ea8779638f078f25daa3913e80fe46eb&query=demain+ne+meurt+jamais",
			},
			moviedb,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readJSONFromUrlMovie(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("readJSONFromUrlMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readJSONFromUrlMovie() got = %v, want %v", got, tt.want)
			}
		})
	}
}
