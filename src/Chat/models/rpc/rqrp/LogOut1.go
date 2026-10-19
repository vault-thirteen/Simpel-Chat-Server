package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type LogOut1Params struct {
	Auth *rpc.Auth `json:"auth"`
}

type LogOut1Result struct {
	RequestId string `json:"requestId"`
}
