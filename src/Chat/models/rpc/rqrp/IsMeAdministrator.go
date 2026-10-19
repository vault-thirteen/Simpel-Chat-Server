package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type IsMeAdministratorParams struct {
	Auth *rpc.Auth `json:"auth"`
}

type IsMeAdministratorResult struct {
	IsAdministrator bool `json:"isAdministrator"`
}
