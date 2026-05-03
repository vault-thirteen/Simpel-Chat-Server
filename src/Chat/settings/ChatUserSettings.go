package settings

import (
	"errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type ChatUserSettings struct {
	AdministratorIds      []common.ObjectId `json:"administratorIds"`
	IsRegistrationEnabled bool              `json:"isRegistrationEnabled"`
}

func (s *ChatUserSettings) Validate() (err error) {
	if s == nil {
		return errors.New(helper.Err_NullPointer)
	}

	if len(s.AdministratorIds) == 0 {
		return helper.NewError_ParameterIsNotSet("administrators")
	}

	return nil
}

func (s *ChatUserSettings) IsUserAdministrator(userId common.ObjectId) bool {
	for _, id := range s.AdministratorIds {
		if id == userId {
			return true
		}
	}

	return false
}
