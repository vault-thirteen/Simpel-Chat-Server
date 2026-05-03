package helper

import (
	"errors"
	"fmt"
	"net"
	"reflect"
)

const (
	Err_DllControllerIsNotInitialised    = "DLL controller is not initialised"
	Err_NullPointer                      = "null pointer"
	Err_ParameterIsNotSet                = "parameter is not set"
	Err_Invalid                          = "invalid"
	Err_Critical                         = "critical error"
	Err_DestinationIsNotInitialised      = "destination is not initialised"
	ErrF_UnsupportedDataType             = "unsupported data type: %s"
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
)

func NewError_Simple2SA(format string, value any) error {
	return errors.New(fmt.Sprintf(format, value))
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
	return NewError_Simple2SA(ErrF_UnsupportedDataType, reflect.TypeOf(src).String())
}

func NewError_WrappedError(text string, err error) error {
	return fmt.Errorf(text+": %w", err)
}

func IsNetworkError(err error) (isNetworkError bool) {
	var nerr net.Error
	return errors.As(err, &nerr)
}
