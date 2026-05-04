package rqrp

import (
	rm "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/room"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ListRoomsParams struct {
	Auth *rpc.Auth `json:"auth,omitempty"`
}

type ListRoomsResult struct {
	Rooms []*rm.RoomForList `json:"rooms"`
}
