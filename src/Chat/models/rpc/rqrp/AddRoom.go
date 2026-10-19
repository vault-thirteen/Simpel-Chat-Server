package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type AddRoomParams struct {
	Auth     *rpc.Auth     `json:"auth"`
	RoomType enum.RoomType `json:"roomType"`
	RoomName string        `json:"roomName"`
}

type AddRoomResult struct {
	RoomId common.ObjectId `json:"roomId"`
}
