package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	lom "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/ListOfMessages"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ListMessagesSinceParams struct {
	Auth       *rpc.Auth       `json:"auth"`
	RoomId     common.ObjectId `json:"roomId"`
	TimeMarkTS int64           `json:"timeMarkTS"`
}

type ListMessagesSinceResult struct {
	Messages *lom.ListOfMessages `json:"messages"`
}
