package controllers

import "testing"

func Test_checkMovieDB(t *testing.T) {
	type args struct {
		tv   bool
		lang bool
		name string
		date []string
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
				[]string{"2017"},
			},
			"https://api.themoviedb.org/3/search/movie?api_key=ea8779638f078f25daa3913e80fe46eb&query=transformers-the-last-knight&primary_release_year=2017",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkMovieDB(tt.args.tv, tt.args.lang, tt.args.name, tt.args.date...); got != tt.want {
				t.Errorf("checkMovieDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
