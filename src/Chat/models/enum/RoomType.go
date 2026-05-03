package enum

import (
	"database/sql/driver"
	"errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"

	"github.com/vault-thirteen/auxie/number"
)

type RoomType byte

const (
	RoomType_Public  = 1
	RoomType_Private = 2
)

func (rt RoomType) Validate() (err error) {
	switch rt {
	case RoomType_Public:
		return nil
	case RoomType_Private:
		return nil
	default:
		return helper.NewError_InvalidEnumValue(EnumField_RoomType, rt)
	}
}

func (rt *RoomType) Scan(src any) (err error) {
	if rt == nil {
		return errors.New(helper.Err_DestinationIsNotInitialised)
	}

	switch src.(type) {
	case []byte:
		{
			var b byte
			b, err = number.ParseUint8(string(src.([]byte)))
			if err != nil {
				return err
			}

			*rt = RoomType(b)

			return nil
		}

	case int64:
		*rt = RoomType(byte(src.(int64)))
		return nil

	case nil:
		return nil

	default:
		return helper.NewError_UnsupportedDataType(src)
	}
}

func (rt RoomType) Value() (dv driver.Value, err error) {
	return int64(rt), nil
}
