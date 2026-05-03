package rm

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
)

const (
	TimeFormat_RoomForList = "2006-01-02 15:04"
)

type RoomForList struct {
	TimeOfCreation   string          `json:"toc"`
	Id               common.ObjectId `json:"id"`
	Type             enum.RoomType   `json:"typ"`
	Name             string          `json:"nam"`
	ActiveUsersCount int             `json:"auc"`
}

func NewRoomForList(room *Room) (rfl *RoomForList) {
	return &RoomForList{
		TimeOfCreation:   room.CreatedAt.UTC().Format(TimeFormat_RoomForList),
		Id:               room.Id,
		Type:             room.Type,
		Name:             room.Name,
		ActiveUsersCount: len(room.activeUsers),
	}
}
