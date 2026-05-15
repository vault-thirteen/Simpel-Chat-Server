package rqrp

type SettingsParams struct{}

type SettingsResult struct {
	MessageSizeMax    int `json:"messageSizeMax"`
	PasswordLengthMin int `json:"passwordLengthMin"`
	PasswordLengthMax int `json:"passwordLengthMax"`
}
