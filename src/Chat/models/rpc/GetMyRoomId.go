package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type GetMyRoomIdParams struct {
	Auth *Auth `json:"auth,omitempty"`
}

type GetMyRoomIdResult struct {
	RoomId *common.ObjectId `json:"roomId"`
}
