package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ListAllowedRoomUsersParams struct {
	Auth   *rpc.Auth       `json:"auth"`
	RoomId common.ObjectId `json:"roomId"`
}

type ListAllowedRoomUsersResult struct {
	RoomId  common.ObjectId   `json:"roomId"`
	UserIds []common.ObjectId `json:"userIds"`
}
