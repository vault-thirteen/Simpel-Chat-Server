package rpc

type LogIn1Params struct {
	EMailAddress string `json:"email"`
}

type LogIn1Result struct {
	RequestId string `json:"requestId"`
}
