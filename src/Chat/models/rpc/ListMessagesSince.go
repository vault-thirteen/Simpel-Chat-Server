package rpc

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	lom "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/ListOfMessages"
)

type ListMessagesSinceParams struct {
	Auth       *Auth           `json:"auth,omitempty"`
	RoomId     common.ObjectId `json:"roomId,omitempty"`
	TimeMarkTS int64           `json:"timeMarkTS,omitempty"`
}

type ListMessagesSinceResult struct {
	Messages *lom.ListOfMessages `json:"messages"`
}
