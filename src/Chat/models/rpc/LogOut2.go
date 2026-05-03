package rpc

type LogOut2Params struct {
	Auth      *Auth  `json:"auth,omitempty"`
	RequestId string `json:"requestId"`
}

type LogOut2Result struct {
	Success
}
