package pwd

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type PasswordController struct {
	PasswordLengthMin int
	PasswordLengthMax int
}

func NewPasswordController(pwdLenMin int, pwdLenMax int) *PasswordController {
	return &PasswordController{
		PasswordLengthMin: pwdLenMin,
		PasswordLengthMax: pwdLenMax,
	}
}

func (pc *PasswordController) IsPasswordTextValid(pwd string) bool {
	l := helper.GetStringLengthInBytes(pwd)

	if (l < pc.PasswordLengthMin) || (l > pc.PasswordLengthMax) {
		return false
	}

	return true
}
