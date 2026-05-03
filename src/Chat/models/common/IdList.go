package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type IdList []ObjectId

func (il *IdList) Scan(src any) (err error) {
	if il == nil {
		return errors.New(helper.Err_DestinationIsNotInitialised)
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
	for _, id := range *il {
		if id == newId {
			return errors.New(helper.Err_DuplicateItem)
		}
	}

	*il = append(*il, newId)

	return nil
}

func (il *IdList) RemoveId(idToRemove ObjectId) (err error) {
	// Find the item.
	var idx int
	var id ObjectId
	var isFound = false

	for idx, id = range *il {
		if id == idToRemove {
			isFound = true
			break
		}
	}

	if !isFound {
		return errors.New(helper.Err_ItemIsNotFound)
	}

	*il = helper.ArrayWithoutItemAt(*il, idx)

	return nil
}

func (il *IdList) HasId(idToCheck ObjectId) bool {
	if il == nil {
		return false
	}

	for _, id := range *il {
		if id == idToCheck {
			return true
		}
	}

	return false
}

func (il *IdList) List() (ids []ObjectId) {
	if il == nil {
		return []ObjectId{}
	}

	return *il
}
