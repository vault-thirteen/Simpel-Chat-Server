package pwd

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type Password struct {
	common.MetaData
	Id     common.ObjectId `gorm:"primarykey"`
	UserId common.ObjectId `gorm:"uniqueIndex"`
	Text   string          `gorm:"column:text"`
}
