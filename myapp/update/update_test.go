package update

import (
	"testing"
)

func TestApplication_checkIfNewVersion(t *testing.T) {
	type fields struct {
		Version    string
		OldVersion string
		Name       string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"checkIfNewVersion", fields{
				OldVersion: "0.9.1.35",
				Version:    "0.9.3.0",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				Version:    tt.fields.Version,
				OldVersion: tt.fields.OldVersion,
				Name:       tt.fields.Name,
			}
			//a.LaunchAppCheckUpdate("0.9.1.35", "search-and-sort-movies-linux-amd64")
			if got := a.checkIfNewVersion(); got != tt.want {
				t.Errorf("checkIfNewVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
