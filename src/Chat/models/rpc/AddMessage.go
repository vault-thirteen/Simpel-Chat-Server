package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type AddMessageParams struct {
	Auth        *Auth           `json:"auth,omitempty"`
	RoomId      common.ObjectId `json:"roomId,omitempty"`
	MessageText string          `json:"messageText,omitempty"`
}

type AddMessageResult struct {
	Success
}
