package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type EnterRoomParams struct {
	Auth   *rpc.Auth       `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type EnterRoomResult struct {
	rpc.Success
}
