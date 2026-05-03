package rpc

import (
	"encoding/json"
	"errors"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	"gorm.io/gorm"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/adc"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/database"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/der"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/generator"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/mailer"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Event"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Password"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Session"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
	rm "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/room"
	lom "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/ListOfMessages"
	msg "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/Message"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/request"
	rpcm "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
	re "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/rpc/errors"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type RpcController struct {
	db               *database.Database
	mailer           *mailer.Mailer
	generator        *generator.Generator
	adc              *adc.ActiveDataController
	der              *der.DatabaseErrorReporter
	chatUserSettings *settings.ChatUserSettings
}

func NewRpcController(
	db *database.Database,
	mailer *mailer.Mailer,
	generator *generator.Generator,
	adc *adc.ActiveDataController,
	der *der.DatabaseErrorReporter,
	chatUserSettings *settings.ChatUserSettings,
) (rc *RpcController) {
	rc = &RpcController{
		db:               db,
		mailer:           mailer,
		generator:        generator,
		adc:              adc,
		der:              der,
		chatUserSettings: chatUserSettings,
	}

	return rc
}

func (rc *RpcController) GetRpcFunctions() []jrm1.RpcFunction {
	return []jrm1.RpcFunction{
		rc.Ping,

		// Auth functions.
		rc.RegisterUser1,
		rc.RegisterUser2,
		rc.LogIn1,
		rc.LogIn2,
		rc.LogOut1,
		rc.LogOut2,
		rc.ChangePassword1,
		rc.ChangePassword2,
		rc.BanUser,

		// Room functions.
		rc.AddRoom,
		rc.DeleteRoom,
		rc.ListRooms,

		// Room Moderator functions.
		rc.AddRoomModerator,
		rc.DeleteRoomModerator,
		rc.ListRoomModerators,
		rc.ResetRoomModerators,

		// Allowed Room User functions.
		rc.AddAllowedRoomUser,
		rc.DeleteAllowedRoomUser,
		rc.ListAllowedRoomUsers,
		rc.ResetAllowedRoomUsers,

		// User Room functions.
		rc.EnterRoom,
		rc.LeaveRoom,
		rc.GetMyRoomId,

		// Message functions.
		rc.AddMessage,
		rc.ListAllMessages,
		rc.ListMessagesSince,
	}
}

// Ping.

func (rc *RpcController) Ping(_ *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	result = rpcm.PingResult{
		Success: rpcm.Success{OK: true},
	}
	return result, nil
}

// Auth functions.

func (rc *RpcController) RegisterUser1(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.RegisterUser1Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.RegisterUser1Result
	r, rpcErr = rc.registerUser1(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) registerUser1(p *rpcm.RegisterUser1Params) (result *rpcm.RegisterUser1Result, rpcErr *jrm1.RpcError) {
	// Is registration enabled ?
	{
		if !rc.chatUserSettings.IsRegistrationEnabled {
			return nil, re.NewRpcError_RegistrationIsDisabled(nil)
		}
	}

	// Check input data.
	{
		if len(p.EMailAddress) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_EMailAddress)
		}
		if !helper.IsEmailAddressValid(p.EMailAddress) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_EMailAddress)
		}
	}

	// Is e-mail address already used ?
	{
		isEmailAddressUsed, err := rc.db.IsEmailAddressUsed(p.EMailAddress)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if isEmailAddressUsed {
			// Make hackers' life a bit more difficult.
			err = helper.SleepBeforeReturningFakeRequestId()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			var fakeRequestId *string
			fakeRequestId, err = rc.generator.RIDG().CreatePassword()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.RegisterUser1Result{RequestId: *fakeRequestId}
			return result, nil
		}
	}

	// Start user registration.
	{
		requestId, err := rc.generator.RIDG().CreatePassword()
		if err != nil {
			return nil, re.NewRpcError_RequestIdGenerator(err)
		}

		var verificationCode *string
		verificationCode, err = rc.generator.VCG().CreatePassword()
		if err != nil {
			return nil, re.NewRpcError_VerificationCodeGenerator(err)
		}

		err = rc.mailer.SendVerificationCode(*verificationCode, p.EMailAddress)
		if err != nil {
			return nil, re.NewRpcError_MailerError(err)
		}

		var reg = rq.Registration{
			Email:            p.EMailAddress,
			RequestId:        *requestId,
			VerificationCode: *verificationCode,
		}

		err = rc.db.CreateRegistrationRequest(&reg)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		result = &rpcm.RegisterUser1Result{RequestId: *requestId}
		return result, nil
	}
}
func (rc *RpcController) RegisterUser2(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.RegisterUser2Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.RegisterUser2Result
	r, rpcErr = rc.registerUser2(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) registerUser2(p *rpcm.RegisterUser2Params) (result *rpcm.RegisterUser2Result, rpcErr *jrm1.RpcError) {
	// Is registration enabled ?
	{
		if !rc.chatUserSettings.IsRegistrationEnabled {
			return nil, re.NewRpcError_RegistrationIsDisabled(nil)
		}
	}

	// Check input data.
	{
		if len(p.EMailAddress) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_EMailAddress)
		}
		if !helper.IsEmailAddressValid(p.EMailAddress) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_EMailAddress)
		}
		if len(p.RequestId) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RequestId)
		}
		if len(p.VerificationCode) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_VerificationCode)
		}
		if len(p.UserName) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserName)
		}
		if !usr.IsUserNameValid(p.UserName) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_UserName)
		}
		if len(p.UserPassword) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserPassword)
		}
		if !pwd.IsPasswordValid(p.UserPassword) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_UserPassword)
		}
	}

	// Check the verification code.
	{
		var rr = &rq.Registration{RequestId: p.RequestId}
		err := rc.db.FindRegistrationRequest(rr)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Fraud.
				err = helper.SleepBeforeFraudResponse()
				if err != nil {
					return nil, re.NewRpcError_RequestIdGenerator(err)
				}

				result = &rpcm.RegisterUser2Result{Success: rpcm.Success{OK: false}}
				return result, nil
			}

			return nil, rc.der.DatabaseError(err)
		}

		if (rr.Email != p.EMailAddress) || (rr.VerificationCode != p.VerificationCode) {
			err = helper.SleepBeforeReportingFailedVerificationCode()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.RegisterUser2Result{Success: rpcm.Success{OK: false}}
			return result, nil
		}
	}

	var user *usr.User

	// Register the user.
	{
		user = &usr.User{
			Name:         p.UserName,
			EmailAddress: p.EMailAddress,
			RegisterTime: time.Now().UTC(),
		}

		rr := new(rq.Registration)
		err := rc.db.GetRegistrationRequestByEmail(p.EMailAddress, rr)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		err = rc.db.CreateUser(user, p.UserPassword)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		err = rc.db.DeleteRegistrationRequest(rr)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		err = rc.mailer.SendRegistrationSuccess(user.EmailAddress)
		if err != nil {
			return nil, re.NewRpcError_MailerError(err)
		}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_UserRegistration,
			TargetUserId: nil,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	result = &rpcm.RegisterUser2Result{Success: rpcm.Success{OK: true}}
	return result, nil
}
func (rc *RpcController) LogIn1(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.LogIn1Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.LogIn1Result
	r, rpcErr = rc.logIn1(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) logIn1(p *rpcm.LogIn1Params) (result *rpcm.LogIn1Result, rpcErr *jrm1.RpcError) {
	// Check input data.
	{
		if len(p.EMailAddress) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_EMailAddress)
		}
		if !helper.IsEmailAddressValid(p.EMailAddress) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_EMailAddress)
		}
	}

	// Does this user exist ?
	{
		userExists, err := rc.db.ExistsUserWithEmail(p.EMailAddress)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if !userExists {
			// Make hackers' life a bit more difficult.
			err = helper.SleepBeforeReturningFakeRequestId()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			var fakeRequestId *string
			fakeRequestId, err = rc.generator.RIDG().CreatePassword()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.LogIn1Result{RequestId: *fakeRequestId}
			return result, nil
		}
	}

	// Is user banned ?
	{
		isUserWithEmailBanned, err := rc.db.IsUserWithEmailBanned(p.EMailAddress)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if isUserWithEmailBanned {
			return nil, re.NewRpcError_UserIsBanned(err)
		}
	}

	// Is user already logged in ?
	{
		isLoggedIn, err := rc.db.IsUserWithEmailLoggedIn(p.EMailAddress)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if isLoggedIn {
			// Make hackers' life a bit more difficult.
			err = helper.SleepBeforeReturningFakeRequestId()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			var fakeRequestId *string
			fakeRequestId, err = rc.generator.RIDG().CreatePassword()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.LogIn1Result{RequestId: *fakeRequestId}
			return result, nil
		}
	}

	// Start logging in.
	{
		requestId, err := rc.generator.RIDG().CreatePassword()
		if err != nil {
			return nil, re.NewRpcError_RequestIdGenerator(err)
		}

		var verificationCode *string
		verificationCode, err = rc.generator.VCG().CreatePassword()
		if err != nil {
			return nil, re.NewRpcError_VerificationCodeGenerator(err)
		}

		err = rc.mailer.SendVerificationCode(*verificationCode, p.EMailAddress)
		if err != nil {
			return nil, re.NewRpcError_MailerError(err)
		}

		var lir = rq.LogIn{
			Email:            p.EMailAddress,
			RequestId:        *requestId,
			VerificationCode: *verificationCode,
		}

		err = rc.db.CreateLogInRequest(&lir)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		result = &rpcm.LogIn1Result{RequestId: *requestId}
		return result, nil
	}
}
func (rc *RpcController) LogIn2(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.LogIn2Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.LogIn2Result
	r, rpcErr = rc.logIn2(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) logIn2(p *rpcm.LogIn2Params) (result *rpcm.LogIn2Result, rpcErr *jrm1.RpcError) {
	// Check input data.
	{
		if len(p.EMailAddress) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_EMailAddress)
		}
		if !helper.IsEmailAddressValid(p.EMailAddress) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_EMailAddress)
		}
		if len(p.RequestId) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RequestId)
		}
		if len(p.VerificationCode) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_VerificationCode)
		}
		if len(p.UserPassword) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserPassword)
		}
		if !pwd.IsPasswordValid(p.UserPassword) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_UserPassword)
		}
	}

	// Does this user exist ?
	{
		userExists, err := rc.db.ExistsUserWithEmail(p.EMailAddress)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if !userExists {
			// Make hackers' life a bit more difficult.
			err = helper.SleepBeforeReturningFakeRequestId()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.LogIn2Result{Success: rpcm.Success{OK: false}}
			return result, nil
		}
	}

	// Is user banned ?
	{
		isUserWithEmailBanned, err := rc.db.IsUserWithEmailBanned(p.EMailAddress)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if isUserWithEmailBanned {
			return nil, re.NewRpcError_UserIsBanned(err)
		}
	}

	// Is user already logged in ?
	{
		isLoggedIn, err := rc.db.IsUserWithEmailLoggedIn(p.EMailAddress)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if isLoggedIn {
			// Make hackers' life a bit more difficult.
			err = helper.SleepBeforeReturningFakeRequestId()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.LogIn2Result{Success: rpcm.Success{OK: false}}
			return result, nil
		}
	}

	// Check the verification code.
	{
		var lir = &rq.LogIn{RequestId: p.RequestId}
		err := rc.db.FindLogInRequest(lir)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Fraud.
				err = helper.SleepBeforeFraudResponse()
				if err != nil {
					return nil, re.NewRpcError_RequestIdGenerator(err)
				}

				result = &rpcm.LogIn2Result{Success: rpcm.Success{OK: false}}
				return result, nil
			}

			return nil, rc.der.DatabaseError(err)
		}

		if (lir.Email != p.EMailAddress) || (lir.VerificationCode != p.VerificationCode) {
			err = helper.SleepBeforeReportingFailedVerificationCode()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.LogIn2Result{Success: rpcm.Success{OK: false}}
			return result, nil
		}
	}

	var user = &usr.User{
		EmailAddress: p.EMailAddress,
	}

	// Check password.
	{
		err := rc.db.FindUserWithEmail(user)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		var ok bool
		ok, err = rc.db.CheckUserPassword(user, &p.UserPassword)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if !ok {
			return nil, re.NewRpcError_WrongPassword(nil)
		}
	}

	var token *string

	// Log the user in.
	{
		lir := new(rq.LogIn)
		err := rc.db.GetLogInRequestByEmail(p.EMailAddress, lir)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		err = rc.db.DeleteLogInRequest(lir)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		token, rpcErr = rc.adc.StartSession(user.Id)
		if rpcErr != nil {
			return nil, rpcErr
		}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_UserLogIn,
			TargetUserId: nil,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	result = &rpcm.LogIn2Result{
		Token:   *token,
		Success: rpcm.Success{OK: true},
	}
	return result, nil
}
func (rc *RpcController) LogOut1(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.LogOut1Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.LogOut1Result
	r, rpcErr = rc.logOut1(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) logOut1(p *rpcm.LogOut1Params) (result *rpcm.LogOut1Result, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Start logging out.
	{
		requestId, err := rc.generator.RIDG().CreatePassword()
		if err != nil {
			return nil, re.NewRpcError_RequestIdGenerator(err)
		}

		var lor = rq.LogOut{
			UserId:    session.UserId,
			RequestId: *requestId,
		}

		err = rc.db.CreateLogOutRequest(&lor)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		result = &rpcm.LogOut1Result{RequestId: *requestId}
		return result, nil
	}
}
func (rc *RpcController) LogOut2(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.LogOut2Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.LogOut2Result
	r, rpcErr = rc.logOut2(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) logOut2(p *rpcm.LogOut2Params) (result *rpcm.LogOut2Result, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if len(p.RequestId) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RequestId)
		}
	}

	var lor *rq.LogOut

	// Check for fraud.
	{
		lor = new(rq.LogOut)
		err := rc.db.GetLogOutRequestByRequestId(p.RequestId, lor)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if session.UserId != lor.UserId {
			// Fraud.
			err = helper.SleepBeforeFraudResponse()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.LogOut2Result{Success: rpcm.Success{OK: false}}
			return result, nil
		}
	}

	// Log the user out.
	{
		err := rc.db.DeleteLogOutRequest(lor)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		rpcErr = rc.adc.StopExistingSession(session, p.Auth.Token)
		if rpcErr != nil {
			return nil, rpcErr
		}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      session.UserId,
			Type:         enum.EventType_UserLogOut,
			TargetUserId: nil,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	result = &rpcm.LogOut2Result{Success: rpcm.Success{OK: true}}
	return result, nil
}
func (rc *RpcController) ChangePassword1(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ChangePassword1Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ChangePassword1Result
	r, rpcErr = rc.changePassword1(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) changePassword1(p *rpcm.ChangePassword1Params) (result *rpcm.ChangePassword1Result, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Start changing password.
	{
		requestId, err := rc.generator.RIDG().CreatePassword()
		if err != nil {
			return nil, re.NewRpcError_RequestIdGenerator(err)
		}

		var verificationCode *string
		verificationCode, err = rc.generator.VCG().CreatePassword()
		if err != nil {
			return nil, re.NewRpcError_VerificationCodeGenerator(err)
		}

		err = rc.mailer.SendVerificationCode(*verificationCode, user.EmailAddress)
		if err != nil {
			return nil, re.NewRpcError_MailerError(err)
		}

		var pcr = rq.ChangePassword{
			UserId:           user.Id,
			RequestId:        *requestId,
			VerificationCode: *verificationCode,
		}

		err = rc.db.CreatePasswordChangeRequest(&pcr)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		result = &rpcm.ChangePassword1Result{RequestId: *requestId}
		return result, nil
	}
}
func (rc *RpcController) ChangePassword2(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ChangePassword2Params
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ChangePassword2Result
	r, rpcErr = rc.changePassword2(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) changePassword2(p *rpcm.ChangePassword2Params) (result *rpcm.ChangePassword2Result, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if len(p.RequestId) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RequestId)
		}
		if len(p.VerificationCode) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_VerificationCode)
		}
		if len(p.UserPassword) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserPassword)
		}
		if !pwd.IsPasswordValid(p.UserPassword) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_UserPassword)
		}
		if len(p.NewUserPassword1) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserPassword)
		}
		if !pwd.IsPasswordValid(p.NewUserPassword1) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_UserPassword)
		}
		if len(p.NewUserPassword2) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserPassword)
		}
		if !pwd.IsPasswordValid(p.NewUserPassword2) {
			return nil, re.NewRpcError_FieldValueIsNotValid(rpcm.Field_UserPassword)
		}
		if p.NewUserPassword1 != p.NewUserPassword2 {
			return nil, re.NewRpcError_NewPasswordsAreDifferent(nil)
		}
	}

	var user *usr.User
	user, rpcErr = rc.getUserWithPassword(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check password.
	if p.UserPassword != user.Password.Text {
		return nil, re.NewRpcError_UserIsNotFound(nil)
	}
	if (user.Password.Text == p.NewUserPassword1) || (user.Password.Text == p.NewUserPassword2) {
		return nil, re.NewRpcError_NewPasswordMustDifferFromExisting(nil)
	}

	var pcr *rq.ChangePassword
	// Check the verification code.
	{
		pcr = &rq.ChangePassword{RequestId: p.RequestId}
		err := rc.db.FindPasswordChangeRequest(pcr)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Fraud.
				err = helper.SleepBeforeFraudResponse()
				if err != nil {
					return nil, re.NewRpcError_RequestIdGenerator(err)
				}

				result = &rpcm.ChangePassword2Result{Success: rpcm.Success{OK: false}}
				return result, nil
			}

			return nil, rc.der.DatabaseError(err)
		}

		if pcr.VerificationCode != p.VerificationCode {
			err = helper.SleepBeforeReportingFailedVerificationCode()
			if err != nil {
				return nil, re.NewRpcError_RequestIdGenerator(err)
			}

			result = &rpcm.ChangePassword2Result{Success: rpcm.Success{OK: false}}
			return result, nil
		}
	}

	// Fraud check.
	if pcr.UserId != user.Id {
		return nil, re.NewRpcError_UserIsNotFound(nil)
	}

	// Change the password.
	{
		err := rc.db.ChangeUserPassword(user, &p.NewUserPassword1)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		err = rc.db.DeletePasswordChangeRequest(pcr)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		// Log the user out.
		rpcErr = rc.adc.StopExistingSession(session, p.Auth.Token)
		if rpcErr != nil {
			return nil, rpcErr
		}

		err = rc.mailer.SendPasswordChangeSuccess(user.EmailAddress)
		if err != nil {
			return nil, re.NewRpcError_MailerError(err)
		}
	}

	// Events' report.
	{
		// 1. 'Password-change' event.
		{
			event := &ev.Event{
				ActorId:      session.UserId,
				Type:         enum.EventType_UserChangePassword,
				TargetUserId: nil,
			}

			err := rc.db.CreateEvent(event)
			if err != nil {
				return nil, rc.der.DatabaseError(err)
			}
		}

		// 2. 'Log-out' event.
		{
			event := &ev.Event{
				ActorId:      session.UserId,
				Type:         enum.EventType_UserLogOut,
				TargetUserId: nil,
			}

			err := rc.db.CreateEvent(event)
			if err != nil {
				return nil, rc.der.DatabaseError(err)
			}
		}
	}

	result = &rpcm.ChangePassword2Result{Success: rpcm.Success{OK: true}}
	return result, nil
}
func (rc *RpcController) BanUser(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.BanUserParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.BanUserResult
	r, rpcErr = rc.banUser(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) banUser(p *rpcm.BanUserParams) (result *rpcm.BanUserResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		if !rc.chatUserSettings.IsUserAdministrator(user.Id) {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Check input data.
	{
		if p.UserId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserId)
		}
		if p.UserId == user.Id {
			return nil, re.NewRpcError_CanNotBanOneself(nil)
		}
	}

	// Does target user exist ?
	{
		targetUserExists, err := rc.db.ExistsUserWithId(&usr.User{Id: p.UserId})
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		if !targetUserExists {
			return nil, re.NewRpcError_UserIsNotFound(nil)
		}
	}

	var sessionWasFound bool

	// Ban the user, close the session.
	{
		var targetUser = &usr.User{Id: p.UserId}

		err := rc.db.BanUserById(targetUser)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}

		// Log the user out, if session exists.
		sessionWasFound, rpcErr = rc.adc.StopUserSessionIfExists(targetUser.Id)
		if rpcErr != nil {
			return nil, rpcErr
		}
	}

	// Events' report.
	{
		// 1. 'Ban' event.
		{
			event := &ev.Event{
				ActorId:      user.Id,
				Type:         enum.EventType_UserBan,
				TargetUserId: &p.UserId,
			}

			err := rc.db.CreateEvent(event)
			if err != nil {
				return nil, rc.der.DatabaseError(err)
			}
		}

		// 2. 'Log-out' event.
		{
			if sessionWasFound {
				event := &ev.Event{
					ActorId:      user.Id,
					Type:         enum.EventType_UserLogOut,
					TargetUserId: &p.UserId,
				}

				err := rc.db.CreateEvent(event)
				if err != nil {
					return nil, rc.der.DatabaseError(err)
				}
			}
		}
	}

	result = &rpcm.BanUserResult{Success: rpcm.Success{OK: true}}
	return result, nil
}

// Room functions.

func (rc *RpcController) AddRoom(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.AddRoomParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.AddRoomResult
	r, rpcErr = rc.addRoom(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) addRoom(p *rpcm.AddRoomParams) (result *rpcm.AddRoomResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		if !rc.chatUserSettings.IsUserAdministrator(user.Id) {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Check input data.
	{
		if p.RoomType == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomType)
		}
		if len(p.RoomName) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomName)
		}
	}

	// Add a room.
	{
		var roomId common.ObjectId
		roomId, rpcErr = rc.adc.AddRoom(p.RoomType, p.RoomName)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.AddRoomResult{
			Success: rpcm.Success{OK: true},
			RoomId:  roomId,
		}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomCreation,
			TargetUserId: nil,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}
func (rc *RpcController) DeleteRoom(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.DeleteRoomParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.DeleteRoomResult
	r, rpcErr = rc.deleteRoom(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) deleteRoom(p *rpcm.DeleteRoomParams) (result *rpcm.DeleteRoomResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		if !rc.chatUserSettings.IsUserAdministrator(user.Id) {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	// Delete a room.
	{
		rpcErr = rc.adc.DeleteExistingRoom(p.RoomId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.DeleteRoomResult{Success: rpcm.Success{OK: true}}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomDeletion,
			TargetUserId: nil,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}
func (rc *RpcController) ListRooms(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ListRoomsParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ListRoomsResult
	r, rpcErr = rc.listRooms(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) listRooms(p *rpcm.ListRoomsParams) (result *rpcm.ListRoomsResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Perform an action.
	{
		var rooms []*rm.RoomForList
		rooms, rpcErr = rc.adc.ListRooms()
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.ListRoomsResult{Rooms: rooms}
	}

	return result, nil
}

// Room Moderator functions.

func (rc *RpcController) AddRoomModerator(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.AddRoomModeratorParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.AddRoomModeratorResult
	r, rpcErr = rc.addRoomModerator(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) addRoomModerator(p *rpcm.AddRoomModeratorParams) (result *rpcm.AddRoomModeratorResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		if !rc.chatUserSettings.IsUserAdministrator(user.Id) {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
		if p.UserId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserId)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.AddRoomModerator(p.RoomId, p.UserId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.AddRoomModeratorResult{Success: rpcm.Success{OK: true}}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomModeratorAddition,
			TargetUserId: &p.UserId,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}
func (rc *RpcController) DeleteRoomModerator(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.DeleteRoomModeratorParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.DeleteRoomModeratorResult
	r, rpcErr = rc.deleteRoomModerator(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) deleteRoomModerator(p *rpcm.DeleteRoomModeratorParams) (result *rpcm.DeleteRoomModeratorResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		if !rc.chatUserSettings.IsUserAdministrator(user.Id) {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
		if p.UserId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserId)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.DeleteRoomModerator(p.RoomId, p.UserId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.DeleteRoomModeratorResult{Success: rpcm.Success{OK: true}}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomModeratorDeletion,
			TargetUserId: &p.UserId,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}
func (rc *RpcController) ListRoomModerators(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ListRoomModeratorsParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ListRoomModeratorsResult
	r, rpcErr = rc.listRoomModerators(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) listRoomModerators(p *rpcm.ListRoomModeratorsParams) (result *rpcm.ListRoomModeratorsResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		if !rc.chatUserSettings.IsUserAdministrator(user.Id) {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	// Perform an action.
	{
		var userIds []common.ObjectId
		userIds, rpcErr = rc.adc.ListRoomModerators(p.RoomId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.ListRoomModeratorsResult{UserIds: userIds}
	}

	return result, nil
}
func (rc *RpcController) ResetRoomModerators(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ResetRoomModeratorsParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ResetRoomModeratorsResult
	r, rpcErr = rc.resetRoomModerators(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) resetRoomModerators(p *rpcm.ResetRoomModeratorsParams) (result *rpcm.ResetRoomModeratorsResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		if !rc.chatUserSettings.IsUserAdministrator(user.Id) {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.ResetRoomModerators(p.RoomId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.ResetRoomModeratorsResult{Success: rpcm.Success{OK: true}}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomModeratorsReset,
			TargetUserId: nil,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}

// Allowed Room User functions.

func (rc *RpcController) AddAllowedRoomUser(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.AddAllowedRoomUserParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.AddAllowedRoomUserResult
	r, rpcErr = rc.addAllowedRoomUser(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) addAllowedRoomUser(p *rpcm.AddAllowedRoomUserParams) (result *rpcm.AddAllowedRoomUserResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
		if p.UserId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserId)
		}
	}

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		var isModerator bool
		isModerator, rpcErr = rc.adc.IsUserModerator(p.RoomId, user.Id)
		if rpcErr != nil {
			return nil, rpcErr
		}

		if !isModerator {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.AddAllowedRoomUser(p.RoomId, p.UserId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.AddAllowedRoomUserResult{Success: rpcm.Success{OK: true}}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomAllowedUserAddition,
			TargetUserId: &p.UserId,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}
func (rc *RpcController) DeleteAllowedRoomUser(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.DeleteAllowedRoomUserParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.DeleteAllowedRoomUserResult
	r, rpcErr = rc.deleteAllowedRoomUser(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) deleteAllowedRoomUser(p *rpcm.DeleteAllowedRoomUserParams) (result *rpcm.DeleteAllowedRoomUserResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
		if p.UserId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_UserId)
		}
	}

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		var isModerator bool
		isModerator, rpcErr = rc.adc.IsUserModerator(p.RoomId, user.Id)
		if rpcErr != nil {
			return nil, rpcErr
		}

		if !isModerator {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.DeleteAllowedRoomUser(p.RoomId, p.UserId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.DeleteAllowedRoomUserResult{Success: rpcm.Success{OK: true}}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomAllowedUserDeletion,
			TargetUserId: &p.UserId,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}
func (rc *RpcController) ListAllowedRoomUsers(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ListAllowedRoomUsersParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ListAllowedRoomUsersResult
	r, rpcErr = rc.listAllowedRoomUsers(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) listAllowedRoomUsers(p *rpcm.ListAllowedRoomUsersParams) (result *rpcm.ListAllowedRoomUsersResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		var isModerator bool
		isModerator, rpcErr = rc.adc.IsUserModerator(p.RoomId, user.Id)
		if rpcErr != nil {
			return nil, rpcErr
		}

		if !isModerator {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Perform an action.
	{
		var userIds []common.ObjectId
		userIds, rpcErr = rc.adc.ListAllowedRoomUsers(p.RoomId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.ListAllowedRoomUsersResult{UserIds: userIds}
	}

	return result, nil
}
func (rc *RpcController) ResetAllowedRoomUsers(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ResetAllowedRoomUsersParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ResetAllowedRoomUsersResult
	r, rpcErr = rc.resetAllowedRoomUsers(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) resetAllowedRoomUsers(p *rpcm.ResetAllowedRoomUsersParams) (result *rpcm.ResetAllowedRoomUsersResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	var user *usr.User
	user, rpcErr = rc.getUser(session)
	if rpcErr != nil {
		return nil, rpcErr
	}

	// Check permissions.
	{
		var isModerator bool
		isModerator, rpcErr = rc.adc.IsUserModerator(p.RoomId, user.Id)
		if rpcErr != nil {
			return nil, rpcErr
		}

		if !isModerator {
			return nil, re.NewRpcError_ActionIsNotPermitted(nil)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.ResetAllowedRoomUsers(p.RoomId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.ResetAllowedRoomUsersResult{Success: rpcm.Success{OK: true}}
	}

	// Event report.
	{
		event := &ev.Event{
			ActorId:      user.Id,
			Type:         enum.EventType_RoomAllowedUserReset,
			TargetUserId: nil,
		}

		err := rc.db.CreateEvent(event)
		if err != nil {
			return nil, rc.der.DatabaseError(err)
		}
	}

	return result, nil
}

// User Room functions.

func (rc *RpcController) EnterRoom(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.EnterRoomParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.EnterRoomResult
	r, rpcErr = rc.enterRoom(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) enterRoom(p *rpcm.EnterRoomParams) (result *rpcm.EnterRoomResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.PutUserIntoRoom(session.UserId, p.RoomId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.EnterRoomResult{Success: rpcm.Success{OK: true}}
	}

	return result, nil
}
func (rc *RpcController) LeaveRoom(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.LeaveRoomParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.LeaveRoomResult
	r, rpcErr = rc.leaveRoom(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) leaveRoom(p *rpcm.LeaveRoomParams) (result *rpcm.LeaveRoomResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.RemoveUserFromRoom(session.UserId, p.RoomId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.LeaveRoomResult{Success: rpcm.Success{OK: true}}
	}

	return result, nil
}
func (rc *RpcController) GetMyRoomId(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.GetMyRoomIdParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.GetMyRoomIdResult
	r, rpcErr = rc.getMyRoomId(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) getMyRoomId(p *rpcm.GetMyRoomIdParams) (result *rpcm.GetMyRoomIdResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Perform an action.
	{
		var roomId *common.ObjectId
		roomId, rpcErr = rc.adc.GetUserRoomId(session.UserId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.GetMyRoomIdResult{RoomId: roomId}
	}

	return result, nil
}

// Message functions.

func (rc *RpcController) AddMessage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.AddMessageParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.AddMessageResult
	r, rpcErr = rc.addMessage(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) addMessage(p *rpcm.AddMessageParams) (result *rpcm.AddMessageResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
		if len(p.MessageText) == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_MessageText)
		}
	}

	// Perform an action.
	{
		rpcErr = rc.adc.AddMessageIntoRoom(p.RoomId, session.UserId, p.MessageText)
		if rpcErr != nil {
			return nil, rpcErr
		}

		result = &rpcm.AddMessageResult{Success: rpcm.Success{OK: true}}
	}

	return result, nil
}
func (rc *RpcController) ListAllMessages(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ListAllMessagesParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ListAllMessagesResult
	r, rpcErr = rc.listAllMessages(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) listAllMessages(p *rpcm.ListAllMessagesParams) (result *rpcm.ListAllMessagesResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
	}

	// Perform an action.
	{
		var rawMsgs []*msg.Message
		rawMsgs, rpcErr = rc.adc.ListAllMessagesInRoom(p.RoomId, session.UserId)
		if rpcErr != nil {
			return nil, rpcErr
		}

		sstts := rc.adc.GetServerStartTimeTS()

		result = &rpcm.ListAllMessagesResult{Messages: lom.NewListOfMessages(p.RoomId, rawMsgs, sstts, nil)}
	}

	return result, nil
}
func (rc *RpcController) ListMessagesSince(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, rpcErr *jrm1.RpcError) {
	var p *rpcm.ListMessagesSinceParams
	rpcErr = jrm1.ParseParameters(params, &p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var r *rpcm.ListMessagesSinceResult
	r, rpcErr = rc.listMessagesSince(p)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return r, nil
}
func (rc *RpcController) listMessagesSince(p *rpcm.ListMessagesSinceParams) (result *rpcm.ListMessagesSinceResult, rpcErr *jrm1.RpcError) {
	var session *ses.Session
	session, rpcErr = rc.getUserSession(p.Auth)
	if rpcErr != nil {
		return nil, rpcErr
	}

	session.TouchLastActivityTime()

	// Check input data.
	{
		if p.RoomId == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_RoomId)
		}
		if p.TimeMarkTS == 0 {
			return nil, re.NewRpcError_FieldNotSet(rpcm.Field_TimeMarkTS)
		}
	}

	// Perform an action.
	{
		var rawMsgs []*msg.Message
		rawMsgs, rpcErr = rc.adc.ListMessagesInRoomSince(p.RoomId, session.UserId, p.TimeMarkTS)
		if rpcErr != nil {
			return nil, rpcErr
		}

		sstts := rc.adc.GetServerStartTimeTS()

		result = &rpcm.ListMessagesSinceResult{Messages: lom.NewListOfMessages(p.RoomId, rawMsgs, sstts, &p.TimeMarkTS)}
	}

	return result, nil
}

// Helper functions.
func (rc *RpcController) getUserSession(auth *rpcm.Auth) (session *ses.Session, rpcErr *jrm1.RpcError) {
	if auth == nil {
		return nil, re.NewRpcError_NotAuthorised(nil)
	}

	if len(auth.Token) == 0 {
		return nil, re.NewRpcError_FieldNotSet(rpcm.Field_AuthToken)
	}

	session, rpcErr = rc.adc.GetUserSession(auth.Token)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return session, nil
}
func (rc *RpcController) getUser(session *ses.Session) (user *usr.User, rpcErr *jrm1.RpcError) {
	user = &usr.User{Id: session.UserId}
	err := rc.db.GetUserById(user)
	if err != nil {
		return nil, rc.der.DatabaseError(err)
	}

	if session.UserId != user.Id {
		return nil, re.NewRpcError_UserIsNotFound(nil)
	}

	return user, nil
}
func (rc *RpcController) getUserWithPassword(session *ses.Session) (user *usr.User, rpcErr *jrm1.RpcError) {
	user = &usr.User{Id: session.UserId}
	err := rc.db.GetUserWithPasswordById(user)
	if err != nil {
		return nil, rc.der.DatabaseError(err)
	}

	if session.UserId != user.Id {
		return nil, re.NewRpcError_UserIsNotFound(nil)
	}

	if user.Password == nil {
		return nil, re.NewRpcError_UserIsNotFound(nil)
	}

	return user, nil
}
