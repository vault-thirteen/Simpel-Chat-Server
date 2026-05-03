package usr

import (
	"time"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Password"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

const (
	UserNameLengthMin = 1
	UserNameLengthMax = 255
)

type User struct {
	common.MetaData
	Id           common.ObjectId `gorm:"primarykey"`
	Name         string          `gorm:"uniqueIndex;size:255"`
	EmailAddress string          `gorm:"column:email;uniqueIndex;size:255"`
	Password     *pwd.Password
	RegisterTime time.Time `gorm:"column:regTime"`

	// Due to a restriction of Go programming language forbidding import loops,
	// while a Session object stores a pointer to a User object, a User object
	// can not store a pointer to a Session object. What a poor language !
	// Golang adepts hating Java language, you can continue hating Java
	// language, but in real life Java is still the best programming language
	// in the world.
	//Session      *ses.Session

	IsBanned bool `gorm:"column:isBanned"`
}

func IsUserNameValid(name string) bool {
	l := helper.GetStringLengthInBytes(name)

	if (l < UserNameLengthMin) || (l > UserNameLengthMax) {
		return false
	}

	return true
}

func (u *User) GetMetaData() (md *common.MetaData) {
	return &u.MetaData
}
