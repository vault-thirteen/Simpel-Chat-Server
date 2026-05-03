package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type EnterRoomParams struct {
	Auth   *Auth           `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type EnterRoomResult struct {
	Success
}
