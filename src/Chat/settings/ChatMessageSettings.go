package settings

import (
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type ChatMessageSettings struct {
	RoomCountMax        int `json:"roomCountMax"`
	RoomMessageCountMax int `json:"roomMessageCountMax"`
	MessageSizeMax      int `json:"messageSizeMax"`
	RoomNameLengthMax   int `json:"roomNameLengthMax"`
}

func (s *ChatMessageSettings) Validate() (err error) {
	if s == nil {
		return errors.ErrorFromText(helper.Err_NullPointer)
	}

	if s.RoomCountMax == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_RoomCountLimit)
	}

	if s.RoomMessageCountMax == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_RoomMessageLimit)
	}

	if s.MessageSizeMax == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_MessageSizeLimit)
	}

	if s.RoomNameLengthMax == 0 {
		return helper.NewError_ParameterIsNotSet(ParameterName_RoomNameLengthLimit)
	}

	return nil
}
