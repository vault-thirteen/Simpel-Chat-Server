package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ResetRoomModeratorsParams struct {
	Auth   *rpc.Auth       `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type ResetRoomModeratorsResult struct {
	rpc.Success
}
