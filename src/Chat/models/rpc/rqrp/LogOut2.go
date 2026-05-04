package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type LogOut2Params struct {
	Auth      *rpc.Auth `json:"auth,omitempty"`
	RequestId string    `json:"requestId"`
}

type LogOut2Result struct {
	rpc.Success
}
