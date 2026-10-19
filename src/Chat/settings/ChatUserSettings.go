package settings

import (
	"slices"

	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type ChatUserSettings struct {
	AdministratorIds      []common.ObjectId `json:"administratorIds"`
	IsRegistrationEnabled bool              `json:"isRegistrationEnabled"`
	PasswordLengthMin     int               `json:"passwordLengthMin"`
	PasswordLengthMax     int               `json:"passwordLengthMax"`
}

func (s *ChatUserSettings) Validate() (err error) {
	if s == nil {
		return errors.ErrorFromText(helper.Err_NullPointer)
	}

	if len(s.AdministratorIds) == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_Administrators)
	}

	if s.PasswordLengthMin == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_PasswordMinimalLength)
	}

	if s.PasswordLengthMax == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_PasswordMaximalLength)
	}

	return nil
}

func (s *ChatUserSettings) IsUserAdministrator(userId common.ObjectId) bool {
	return slices.Contains(s.AdministratorIds, userId)
}
