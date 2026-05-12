package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type GetMyRoomIdParams struct {
	Auth *rpc.Auth `json:"auth"`
}

type GetMyRoomIdResult struct {
	UserId common.ObjectId  `json:"userId"`
	RoomId *common.ObjectId `json:"roomId"`
}
