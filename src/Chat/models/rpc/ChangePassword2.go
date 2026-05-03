package rpc

type ChangePassword2Params struct {
	Auth             *Auth  `json:"auth,omitempty"`
	RequestId        string `json:"requestId"`
	VerificationCode string `json:"verificationCode"`
	UserPassword     string `json:"userPassword"`
	NewUserPassword1 string `json:"newUserPassword1"`
	NewUserPassword2 string `json:"newUserPassword2"`
}

type ChangePassword2Result struct {
	Success
}
