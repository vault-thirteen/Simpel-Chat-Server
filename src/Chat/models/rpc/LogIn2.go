package rpc

type LogIn2Params struct {
	EMailAddress     string `json:"email"`
	RequestId        string `json:"requestId"`
	VerificationCode string `json:"verificationCode"`
	UserPassword     string `json:"userPassword"`
}

type LogIn2Result struct {
	Token string `json:"token"`
	Success
}
