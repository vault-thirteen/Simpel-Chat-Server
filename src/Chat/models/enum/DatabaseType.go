package enum

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type DatabaseType string

const (
	DatabaseType_MySQL = "mysql"
)

func (dt DatabaseType) Validate() (err error) {
	switch dt {
	case DatabaseType_MySQL:
		return nil
	default:
		return helper.NewError_InvalidEnumValue(EnumField_DatabaseType, dt)
	}
}

func (dt DatabaseType) ToString() string {
	return string(dt)
}
