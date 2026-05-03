package pwd

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

const (
	PasswordLengthMin = 8
	PasswordLengthMax = 32
)

type Password struct {
	common.MetaData
	Id     common.ObjectId `gorm:"primarykey"`
	UserId common.ObjectId `gorm:"uniqueIndex"`
	Text   string          `gorm:"column:text"`
}

func IsPasswordValid(password string) bool {
	l := helper.GetStringLengthInBytes(password)

	if (l < PasswordLengthMin) || (l > PasswordLengthMax) {
		return false
	}

	return true
}
