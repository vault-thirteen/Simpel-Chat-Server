package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type GetRoomUsersParams struct {
	Auth   *rpc.Auth       `json:"auth"`
	RoomId common.ObjectId `json:"roomId"`
}

type GetRoomUsersResult struct {
	RoomId        common.ObjectId   `json:"roomId"`
	ActiveUserIds []common.ObjectId `json:"activeUserIds"`
}
