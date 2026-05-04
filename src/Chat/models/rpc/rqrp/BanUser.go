package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type BanUserParams struct {
	Auth   *rpc.Auth       `json:"auth,omitempty"`
	UserId common.ObjectId `json:"userId" gorm:"uniqueIndex"`
}

type BanUserResult struct {
	rpc.Success
}
