package helper

import (
	"errors"
	"net"
	"reflect"

	"github.com/vault-thirteen/auxie/cio"
	ae "github.com/vault-thirteen/auxie/errors"
)

const (
	Err_NullPointer                      = "null pointer"
	Err_ParameterIsNotSet                = "parameter is not set"
	Err_Invalid                          = "invalid"
	Err_Critical                         = "critical error"
	Err_DestinationIsNotInitialised      = "destination is not initialised"
	Err_PathToConfigurationFileIsNotSet  = "path to configuration file is not set"
	Err_ADC                              = "ADC error"
	Err_SessionsCountMismatch            = "sessions count mismatch"
	Err_RoomCountMismatch                = "room count mismatch"
	Err_DuplicateRoom                    = "duplicate room"
	Err_RoomsListIsAlreadyLoaded         = "rooms list is already loaded"
	Err_PublicRoomCanNotHaveAllowedUsers = "public room can not have allowed users"
	Err_IdIsNotSet                       = "ID is not set"
	Err_DuplicateItem                    = "duplicate item"
	Err_ItemIsNotFound                   = "item is not found"
	Err_RoomIsNotFound                   = "room is not found"
	Err_NoRowsAreAffected                = "no rows were affected"
	Err_UserIsAlreadyInTheRoom           = "user already in the room"
	Err_UserIsNotInTheRoom               = "user is not in the room"
	Err_UserIsNotAllowedToUseThisRoom    = "user is not allowed to use this room"
	Err_UnknownRoomType                  = "unknown room type"
	Err_MessageIsTooLong                 = "message is too long"
	Err_NameIsTooLong                    = "name is too long"
	Err_YouAreNotModerator               = "you are not a moderator"

	Errf_UnsupportedDataType = "unsupported data type: %s"
)

func NewError_Simple2SA(format string, value any) error {
	return ae.ErrorFromFormattedText(cio.FormatMessage(format, value))
}

func NewError_ParameterIsNotSet(parameterName string) error {
	return NewError_Simple2SA(Err_ParameterIsNotSet+": %+v", parameterName)
}

func NewError_InvalidEnumValue(enumName string, value any) error {
	return NewError_Simple2SA(Err_Invalid+" "+enumName+": %+v", value)
}

func NewError_GenericError(text string, value any) error {
	return NewError_Simple2SA(text+": %+v", value)
}

func NewError_UnsupportedDataType(src any) error {
	return NewError_Simple2SA(Errf_UnsupportedDataType, reflect.TypeOf(src).String())
}

func IsNetworkError(err error) (isNetworkError bool) {
	var nerr net.Error
	return errors.As(err, &nerr)
}
