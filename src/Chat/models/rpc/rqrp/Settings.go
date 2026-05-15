package rqrp

type Settings struct {
	MessageSizeMax    int `json:"messageSizeMax"`
	PasswordLengthMin int `json:"passwordLengthMin"`
	PasswordLengthMax int `json:"passwordLengthMax"`
}
