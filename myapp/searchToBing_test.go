package myapp

import "testing"

func Test_loopGetBingName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"loopGetBingName", args{
				"Boruto 81 Vostfr",
			},
			"boruto+81+vostfr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loopGetBingName(tt.args.name); got != tt.want {
				t.Errorf("loopGetBingName() = %v, want %v", got, tt.want)
			}
		})
	}
}
