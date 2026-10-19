package rqrp

type VersionParams = struct{}

type VersionResult struct {
	ServerName        string `json:"serverName"`
	ChatFamily        string `json:"chatFamily"`
	AppName           string `json:"appName"`
	AppVersionText    string `json:"appVersion"`
	GolangVersionText string `json:"golang"`
}
