package myapp

import (
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
