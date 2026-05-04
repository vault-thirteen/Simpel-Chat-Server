package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	lom "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/ListOfMessages"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ListAllMessagesParams struct {
	Auth   *rpc.Auth       `json:"auth,omitempty"`
	RoomId common.ObjectId `json:"roomId,omitempty"`
}

type ListAllMessagesResult struct {
	Messages *lom.ListOfMessages `json:"messages"`
}
