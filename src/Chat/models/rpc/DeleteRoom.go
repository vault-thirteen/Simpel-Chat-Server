package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type DeleteRoomParams struct {
	Auth   *Auth           `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type DeleteRoomResult struct {
	Success
}
