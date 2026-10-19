package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	rm "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/room"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type GetRoomParams struct {
	Auth   *rpc.Auth       `json:"auth"`
	RoomId common.ObjectId `json:"roomId"`
}

type GetRoomResult struct {
	Room *rm.Room `json:"room"`
}
