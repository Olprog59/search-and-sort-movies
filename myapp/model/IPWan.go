package model

type IPWan struct {
	Method     string `json:"method"`
	Port       int    `json:"port"`
	Forwarded  string `json:"forwarded"`
	RemoteHost string `json:"remote_host"`
	Mime       string `json:"mime"`
	Via        string `json:"via"`
	Encoding   string `json:"encoding"`
	Language   string `json:"language"`
	IPAddr     string `json:"ip_addr"`
	UserAgent  string `json:"user_agent"`
}
