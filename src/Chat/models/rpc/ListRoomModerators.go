package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type ListRoomModeratorsParams struct {
	Auth   *Auth           `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type ListRoomModeratorsResult struct {
	UserIds []common.ObjectId `json:"userIds"`
}
