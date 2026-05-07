package re

import (
	"fmt"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	fe "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/rpc/errors/FlexibleError"
)

// RPC error codes.
const (
	Code_Failure                           = 1
	Code_DatabaseError                     = 2
	Code_FieldIsNotSet                     = 3
	Code_FieldValueIsNotValid              = 4
	Code_RequestIdGenerator                = 5
	Code_VerificationCodeGenerator         = 6
	Code_TokenGenerator                    = 7
	Code_MailerError                       = 8
	Code_RegistrationIsDisabled            = 9
	Code_ActionIsNotPermitted              = 10
	Code_UserIsBanned                      = 11
	Code_SessionAlreadyExists              = 12
	Code_WrongPassword                     = 13
	Code_SessionCountLimit                 = 14
	Code_ActiveDataController              = 15
	Code_NotAuthorised                     = 16
	Code_SessionIsNotFound                 = 17
	Code_UserIsNotFound                    = 18
	Code_NewPasswordsAreDifferent          = 19
	Code_NewPasswordMustDifferFromExisting = 20
	Code_CanNotBanOneself                  = 21
	Code_RoomCountLimit                    = 22
	Code_RoomError                         = 23
	Code_UserCanNotUseMultipleRooms        = 24
	Code_UserIsNotUsingAnyRoom             = 25
	Code_RoomDoesNotExist                  = 26
	Code_RoomIsNotFound                    = 27
	Code_UserIsNotAllowedInTheRoom         = 28
	Code_SessionHasTimedOut                = 29
	Code_ShortStringIsTooLong              = 30
)

// RPC error messages.
const (
	Msg_Failure                           = "failure"
	Msg_DatabaseError                     = "database error"
	MsgF_FieldIsNotSet                    = "field is not set: %s"
	MsgF_FieldValueIsNotValid             = "field value is not valid: %s"
	Msg_RequestIdGenerator                = "request ID generator error"
	Msg_VerificationCodeGenerator         = "verification code generator error"
	Msg_TokenGenerator                    = "token generator error"
	Msg_MailerError                       = "mailer error"
	Msg_RegistrationIsDisabled            = "registration is disabled"
	Msg_ActionIsNotPermitted              = "action is not permitted"
	Msg_UserIsBanned                      = "user is banned"
	Msg_SessionAlreadyExists              = "session already exists"
	Msg_WrongPassword                     = "wrong password"
	Msg_SessionCountLimit                 = "session count limit"
	Msg_ActiveDataController              = "active data controller"
	Msg_NotAuthorised                     = "not authorised"
	Msg_SessionIsNotFound                 = "session is not found"
	Msg_UserIsNotFound                    = "user is not found"
	Msg_NewPasswordsAreDifferent          = "new passwords are different"
	Msg_NewPasswordMustDifferFromExisting = "new password must differ from existing password"
	Msg_CanNotBanOneself                  = "can not ban oneself"
	Msg_RoomCountLimit                    = "room count limit"
	Msg_RoomError                         = "room error"
	Msg_UserCanNotUseMultipleRooms        = "user can not use more than one room at the same time"
	Msg_UserIsNotUsingAnyRoom             = "user is not using any room"
	Msg_RoomDoesNotExist                  = "room does not exist"
	Msg_RoomIsNotFound                    = "room is not found"
	Msg_UserIsNotAllowedInTheRoom         = "user is not allowed in this room"
	Msg_SessionHasTimedOut                = "session has timed out"
	Msg_ShortStringIsTooLong              = "short string is too long"
)

func NewRpcError_Failure(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_Failure, Msg_Failure, fe.NewFlexibleError(err).Value())
}
func NewRpcError_DatabaseError(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_DatabaseError, Msg_DatabaseError, fe.NewFlexibleError(err).Value())
}
func NewRpcError_FieldNotSet(fieldName string) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_FieldIsNotSet, fmt.Sprintf(MsgF_FieldIsNotSet, fieldName), nil)
}
func NewRpcError_FieldValueIsNotValid(fieldName string) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_FieldValueIsNotValid, fmt.Sprintf(MsgF_FieldValueIsNotValid, fieldName), nil)
}
func NewRpcError_RequestIdGenerator(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_RequestIdGenerator, Msg_RequestIdGenerator, fe.NewFlexibleError(err).Value())
}
func NewRpcError_VerificationCodeGenerator(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_VerificationCodeGenerator, Msg_VerificationCodeGenerator, fe.NewFlexibleError(err).Value())
}
func NewRpcError_TokenGenerator(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_TokenGenerator, Msg_TokenGenerator, fe.NewFlexibleError(err).Value())
}
func NewRpcError_MailerError(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_MailerError, Msg_MailerError, fe.NewFlexibleError(err).Value())
}
func NewRpcError_RegistrationIsDisabled(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_RegistrationIsDisabled, Msg_RegistrationIsDisabled, fe.NewFlexibleError(err).Value())
}
func NewRpcError_ActionIsNotPermitted(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_ActionIsNotPermitted, Msg_ActionIsNotPermitted, fe.NewFlexibleError(err).Value())
}
func NewRpcError_UserIsBanned(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_UserIsBanned, Msg_UserIsBanned, fe.NewFlexibleError(err).Value())
}
func NewRpcError_SessionAlreadyExists(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_SessionAlreadyExists, Msg_SessionAlreadyExists, fe.NewFlexibleError(err).Value())
}
func NewRpcError_WrongPassword(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_WrongPassword, Msg_WrongPassword, fe.NewFlexibleError(err).Value())
}
func NewRpcError_SessionCountLimit(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_SessionCountLimit, Msg_SessionCountLimit, fe.NewFlexibleError(err).Value())
}
func NewRpcError_ActiveDataController(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_ActiveDataController, Msg_ActiveDataController, fe.NewFlexibleError(err).Value())
}
func NewRpcError_NotAuthorised(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_NotAuthorised, Msg_NotAuthorised, fe.NewFlexibleError(err).Value())
}
func NewRpcError_SessionIsNotFound(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_SessionIsNotFound, Msg_SessionIsNotFound, fe.NewFlexibleError(err).Value())
}
func NewRpcError_UserIsNotFound(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_UserIsNotFound, Msg_UserIsNotFound, fe.NewFlexibleError(err).Value())
}
func NewRpcError_NewPasswordsAreDifferent(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_NewPasswordsAreDifferent, Msg_NewPasswordsAreDifferent, fe.NewFlexibleError(err).Value())
}
func NewRpcError_NewPasswordMustDifferFromExisting(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_NewPasswordMustDifferFromExisting, Msg_NewPasswordMustDifferFromExisting, fe.NewFlexibleError(err).Value())
}
func NewRpcError_CanNotBanOneself(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_CanNotBanOneself, Msg_CanNotBanOneself, fe.NewFlexibleError(err).Value())
}
func NewRpcError_RoomCountLimit(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_RoomCountLimit, Msg_RoomCountLimit, fe.NewFlexibleError(err).Value())
}
func NewRpcError_RoomError(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_RoomError, Msg_RoomError, fe.NewFlexibleError(err).Value())
}
func NewRpcError_UserCanNotUseMultipleRooms(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_UserCanNotUseMultipleRooms, Msg_UserCanNotUseMultipleRooms, fe.NewFlexibleError(err).Value())
}
func NewRpcError_UserIsNotUsingAnyRoom(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_UserIsNotUsingAnyRoom, Msg_UserIsNotUsingAnyRoom, fe.NewFlexibleError(err).Value())
}
func NewRpcError_RoomDoesNotExist(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_RoomDoesNotExist, Msg_RoomDoesNotExist, fe.NewFlexibleError(err).Value())
}
func NewRpcError_RoomIsNotFound(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_RoomIsNotFound, Msg_RoomIsNotFound, fe.NewFlexibleError(err).Value())
}
func NewRpcError_Msg_UserIsNotAllowedInTheRoom(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_UserIsNotAllowedInTheRoom, Msg_UserIsNotAllowedInTheRoom, fe.NewFlexibleError(err).Value())
}
func NewRpcError_SessionHasTimedOut(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_SessionHasTimedOut, Msg_SessionHasTimedOut, fe.NewFlexibleError(err).Value())
}
func NewRpcError_ShortStringIsTooLong(err error) (re *jrm1.RpcError) {
	return jrm1.NewRpcErrorByUser(Code_ShortStringIsTooLong, Msg_ShortStringIsTooLong, fe.NewFlexibleError(err).Value())
}
