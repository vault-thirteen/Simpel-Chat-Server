package rqrp

type RegisterUser2Params struct {
	EMailAddress     string `json:"email"`
	RequestId        string `json:"requestId"`
	VerificationCode string `json:"verificationCode"`
	UserName         string `json:"userName"`
	UserPassword     string `json:"userPassword"`
}

type RegisterUser2Result struct{}
