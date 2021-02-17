package myapp

import "testing"

func Test_loopGetSearchEngine(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"loopGetSearchEngine", args{
				"Boruto 81 Vostfr",
			},
			"boruto+81+vostfr",
		},
		{
			"loopGetSearchEngine", args{
				"narutto",
			},
			"naruto",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loopGetSearchEngine(tt.args.name); got != tt.want {
				t.Errorf("loopGetSearchEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSearchEngine(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"getSearchEngine", args{
				"narutto",
			},
			"naruto",
		},
		{
			"getSearchEngine", args{
				"the hundred",
			},
			"the 100",
		},
		{
			"getSearchEngine", args{
				"Motherland.Fort.Salem",
			},
			"motherland: fort salem",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := getSearchEngine(tt.args.name); got != tt.want {
				t.Errorf("getSearchEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}
