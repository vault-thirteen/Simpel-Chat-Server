package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type AddMessageParams struct {
	Auth        *rpc.Auth       `json:"auth,omitempty"`
	RoomId      common.ObjectId `json:"roomId,omitempty"`
	MessageText string          `json:"messageText,omitempty"`
}

type AddMessageResult struct {
	rpc.Success
}
