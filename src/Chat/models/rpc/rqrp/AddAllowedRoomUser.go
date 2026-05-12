package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type AddAllowedRoomUserParams struct {
	Auth   *rpc.Auth       `json:"auth"`
	RoomId common.ObjectId `json:"roomId"`
	UserId common.ObjectId `json:"userId"`
}

type AddAllowedRoomUserResult struct{}
