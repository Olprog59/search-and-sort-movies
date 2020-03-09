package model

type BuildInfo struct {
	BuildName    string `json:"build_name"`
	BuildVersion string `json:"build_version"`
	BuildHash    string `json:"build_hash"`
	BuildDate    string `json:"build_date"`
	BuildClean   string `json:"build_clean"`
	Os           string `json:"os"`
	Architecture string `json:"architecture"`
}
