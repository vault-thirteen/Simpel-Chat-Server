package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type LeaveRoomParams struct {
	Auth   *Auth           `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type LeaveRoomResult struct {
	Success
}
