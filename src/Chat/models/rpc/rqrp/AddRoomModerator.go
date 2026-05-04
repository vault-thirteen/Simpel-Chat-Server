package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type AddRoomModeratorParams struct {
	Auth   *rpc.Auth       `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
	UserId common.ObjectId `json:"userId,omitempty"`
}

type AddRoomModeratorResult struct {
	rpc.Success
}
