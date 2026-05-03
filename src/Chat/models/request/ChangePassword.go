package rq

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
)

type ChangePassword struct {
	common.MetaData
	Id               common.ObjectId `json:"id" gorm:"primarykey"`
	UserId           common.ObjectId `json:"-" gorm:"uniqueIndex"`
	User             *usr.User       `json:"-"`
	RequestId        string          `json:"requestId" gorm:"uniqueIndex;column:requestId;size:8"`
	VerificationCode string          `json:"verificationCode" gorm:"column:verificationCode;size:6"`
}
