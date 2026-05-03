package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type ResetRoomModeratorsParams struct {
	Auth   *Auth           `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type ResetRoomModeratorsResult struct {
	Success
}
