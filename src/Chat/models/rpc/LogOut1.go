package rpc

type LogOut1Params struct {
	Auth *Auth `json:"auth,omitempty"`
}

type LogOut1Result struct {
	RequestId string `json:"requestId"`
}
