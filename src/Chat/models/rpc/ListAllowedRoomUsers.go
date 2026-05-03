package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type ListAllowedRoomUsersParams struct {
	Auth   *Auth           `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type ListAllowedRoomUsersResult struct {
	UserIds []common.ObjectId `json:"userIds"`
}
