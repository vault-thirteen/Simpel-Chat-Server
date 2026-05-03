package rpc

import (
	rm "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/room"
)

type ListRoomsParams struct {
	Auth *Auth `json:"auth,omitempty"`
}

type ListRoomsResult struct {
	Rooms []*rm.RoomForList `json:"rooms"`
}
