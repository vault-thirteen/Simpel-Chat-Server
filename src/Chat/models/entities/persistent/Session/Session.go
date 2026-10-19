package ses

import (
	"time"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
)

type Session struct {
	common.MetaData
	Id        common.ObjectId `json:"id" gorm:"primarykey;column:id"`
	UserId    common.ObjectId `json:"-" gorm:"uniqueIndex"`
	User      *usr.User       `json:"-"`
	LogInTime time.Time       `json:"logInTime" gorm:"column:logInTime"`
	Token     *string         `json:"-" gorm:"-"`

	lastActivityTimeTS int64 `gorm:"-"`
}

func NewSession(userId common.ObjectId, token *string) (s *Session) {
	// Here we fill some of the fields which are known.
	// Other fields, such as ID, are set by a database.
	s = &Session{
		UserId:    userId,
		LogInTime: time.Now().UTC(),
		Token:     token,
	}

	s.TouchLastActivityTime()

	return s
}

func (s *Session) TouchLastActivityTime() {
	s.lastActivityTimeTS = time.Now().UTC().Unix()
}
func (s *Session) GetLastActivityTimeTS() int64 {
	return s.lastActivityTimeTS
}
