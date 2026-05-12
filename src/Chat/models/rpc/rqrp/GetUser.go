package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	usr "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type GetUserParams struct {
	Auth   *rpc.Auth       `json:"auth"`
	UserId common.ObjectId `json:"userId"`
}

type GetUserResult struct {
	User *usr.User2 `json:"user"`
}
