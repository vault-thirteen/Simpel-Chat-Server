package rpc

type RegisterUser1Params struct {
	EMailAddress string `json:"email"`
}

type RegisterUser1Result struct {
	RequestId string `json:"requestId"`
}
