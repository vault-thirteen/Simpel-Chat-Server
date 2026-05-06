package settings

import (
	"errors"

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
		return errors.New(helper.Err_NullPointer)
	}

	if s.RoomCountMax == 0 {
		return helper.NewError_ParameterIsNotSet("room count limit")
	}

	if s.RoomMessageCountMax == 0 {
		return helper.NewError_ParameterIsNotSet("room message limit")
	}

	if s.MessageSizeMax == 0 {
		return helper.NewError_ParameterIsNotSet("message size limit")
	}

	if s.RoomNameLengthMax == 0 {
		return helper.NewError_ParameterIsNotSet("room name length limit")
	}

	return nil
}
