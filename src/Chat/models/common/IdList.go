package common

import (
	"database/sql/driver"
	"encoding/json"
	"slices"

	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type IdList []ObjectId

func (il *IdList) Scan(src any) (err error) {
	if il == nil {
		return errors.ErrorFromText(helper.Err_DestinationIsNotInitialised)
	}

	switch src.(type) {
	case []byte:
		{
			data := new(IdList)

			err = json.Unmarshal(src.([]byte), data)
			if err != nil {
				return err
			}

			if data != nil {
				*il = *data
			}

			return nil
		}

	case nil:
		return nil

	default:
		return helper.NewError_UnsupportedDataType(src)
	}
}

func (il *IdList) Value() (dv driver.Value, err error) {
	if il == nil {
		return nil, nil
	}

	var buf []byte
	buf, err = json.Marshal(il)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (il *IdList) AsArray() (ids []ObjectId) {
	return *il
}

func (il *IdList) AddId(newId ObjectId) (err error) {
	// Check for duplicates.
	if slices.Contains(*il, newId) {
		return errors.ErrorFromText(helper.Err_DuplicateItem)
	}

	*il = append(*il, newId)

	return nil
}

func (il *IdList) RemoveId(idToRemove ObjectId) (err error) {
	// Find the item.
	idx := slices.Index(*il, idToRemove)

	if idx < 0 {
		return errors.ErrorFromText(helper.Err_ItemIsNotFound)
	}

	*il = slices.Delete(*il, idx, idx+1)

	return nil
}

func (il *IdList) HasId(idToCheck ObjectId) bool {
	if il == nil {
		return false
	}

	return slices.Contains(*il, idToCheck)
}

func (il *IdList) List() (ids []ObjectId) {
	if il == nil {
		return []ObjectId{}
	}

	return *il
}
