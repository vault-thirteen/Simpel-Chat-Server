package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type AddRoomParams struct {
	Auth     *rpc.Auth     `json:"auth,omitempty"`
	RoomType enum.RoomType `json:"roomType,omitempty"`
	RoomName string        `json:"roomName,omitempty"`
}

type AddRoomResult struct {
	rpc.Success
	RoomId common.ObjectId `json:"roomId,omitempty"`
}
