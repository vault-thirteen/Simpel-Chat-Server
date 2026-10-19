package usr

import (
	"time"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type User2 struct {
	Id           common.ObjectId `json:"id"`
	Name         string          `json:"name"`
	EmailAddress *string         `json:"emailAddress,omitempty"`
	RegisterTime time.Time       `json:"registerTime"`
	IsBanned     bool            `json:"isBanned"`
}

func NewUser2(u *User, showEmail bool) (u2 *User2) {
	if u == nil {
		return nil
	}

	u2 = &User2{
		Id:           u.Id,
		Name:         u.Name,
		EmailAddress: nil,
		RegisterTime: u.RegisterTime,
		IsBanned:     u.IsBanned,
	}

	if showEmail {
		u2.EmailAddress = &u.EmailAddress
	}

	return u2
}

func NewUsers2(users []*User, showEmail bool) (users2 []*User2) {
	users2 = make([]*User2, 0, len(users))

	for _, u := range users {
		users2 = append(users2, NewUser2(u, showEmail))
	}

	return users2
}
