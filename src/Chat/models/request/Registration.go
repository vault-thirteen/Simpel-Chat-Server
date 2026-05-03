package rq

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type Registration struct {
	common.MetaData
	Id               common.ObjectId `json:"id" gorm:"primarykey"`
	Email            string          `json:"email" gorm:"uniqueIndex;size:255"`
	RequestId        string          `json:"requestId" gorm:"uniqueIndex;column:requestId;size:8"`
	VerificationCode string          `json:"verificationCode" gorm:"column:verificationCode;size:6"`
}
