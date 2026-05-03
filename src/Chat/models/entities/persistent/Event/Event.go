package ev

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
)

type Event struct {
	common.MetaData
	Id           common.ObjectId  `gorm:"primarykey"`
	ActorId      common.ObjectId  `gorm:"column:actorId"`
	Type         enum.EventType   `gorm:"column:type"`
	TargetUserId *common.ObjectId `gorm:"column:targetUserId"`
}

func (e *Event) HasValidType() bool {
	return e.Type.IsValid()
}
