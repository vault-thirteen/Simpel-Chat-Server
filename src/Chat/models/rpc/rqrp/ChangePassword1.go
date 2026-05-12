package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ChangePassword1Params struct {
	Auth *rpc.Auth `json:"auth"`
}

type ChangePassword1Result struct {
	RequestId string `json:"requestId"`
}
