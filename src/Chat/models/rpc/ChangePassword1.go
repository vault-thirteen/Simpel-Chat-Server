package rpc

type ChangePassword1Params struct {
	Auth *Auth `json:"auth,omitempty"`
}

type ChangePassword1Result struct {
	RequestId string `json:"requestId"`
}
