package adc

import (
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	rs "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/adc/Rooms"
	ss "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/adc/Sessions"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/database"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/der"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/generator"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Session"
	rm "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/room"
	msg "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/Message"
	rp "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/RoomParameters"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
)

type ActiveDataController struct {
	guard             sync.RWMutex
	serverStartTimeTS int64
	sessions          *ss.Sessions
	rooms             *rs.Rooms
}

func NewActiveDataController(
	db *database.Database,
	der *der.DatabaseErrorReporter,
	criticalErrorsChan *chan error,
	generator *generator.Generator,
	cs *settings.ChatSettings,
	serverStartTimeTS int64,
) (adc *ActiveDataController, err error) {
	adc = &ActiveDataController{
		guard:             sync.RWMutex{},
		serverStartTimeTS: serverStartTimeTS,
	}

	adc.sessions, err = ss.NewSessions(db, der, criticalErrorsChan, generator, cs.Other.SessionCountMax, cs.Other.InactivityDurationMaxSec)
	if err != nil {
		return nil, err
	}

	commonRoomParameters := rp.NewRoomParameters(cs.Message.RoomMessageCountMax, cs.Message.MessageSizeMax)

	adc.rooms, err = rs.NewRooms(db, der, criticalErrorsChan, cs.Message.RoomCountMax, commonRoomParameters, serverStartTimeTS)
	if err != nil {
		return nil, err
	}

	return adc, nil
}

// Session functions.

func (adc *ActiveDataController) StartSession(userId common.ObjectId) (token *string, rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.sessions.StartSession(userId)
}
func (adc *ActiveDataController) GetUserSession(token string) (session *ses.Session, rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.sessions.GetUserSession(token)
}
func (adc *ActiveDataController) StopExistingSession(session *ses.Session, token string) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.sessions.StopExistingSession(session, token)
}
func (adc *ActiveDataController) StopUserSessionIfExists(userId common.ObjectId) (sessionWasFound bool, rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.sessions.StopUserSessionIfExists(userId)
}

// Room functions.

func (adc *ActiveDataController) AddRoom(typé enum.RoomType, name string) (roomId common.ObjectId, rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.AddRoom(typé, name)
}
func (adc *ActiveDataController) DeleteExistingRoom(roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.DeleteExistingRoom(roomId)
}
func (adc *ActiveDataController) ListRooms() (rooms []*rm.RoomForList, rpcErr *jrm1.RpcError) {
	adc.guard.RLock()
	defer adc.guard.RUnlock()
	return adc.rooms.ListRooms()
}

// Room Moderator functions.

func (adc *ActiveDataController) AddRoomModerator(roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.AddRoomModerator(roomId, userId)
}
func (adc *ActiveDataController) DeleteRoomModerator(roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.DeleteRoomModerator(roomId, userId)
}
func (adc *ActiveDataController) ListRoomModerators(roomId common.ObjectId) (userIds []common.ObjectId, rpcErr *jrm1.RpcError) {
	adc.guard.RLock()
	defer adc.guard.RUnlock()
	return adc.rooms.ListRoomModerators(roomId)
}
func (adc *ActiveDataController) ResetRoomModerators(roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.ResetRoomModerators(roomId)
}

// Allowed Room User functions.

func (adc *ActiveDataController) AddAllowedRoomUser(roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.AddAllowedRoomUser(roomId, userId)
}
func (adc *ActiveDataController) DeleteAllowedRoomUser(roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.DeleteAllowedRoomUser(roomId, userId)
}
func (adc *ActiveDataController) ListAllowedRoomUsers(roomId common.ObjectId) (userIds []common.ObjectId, rpcErr *jrm1.RpcError) {
	adc.guard.RLock()
	defer adc.guard.RUnlock()
	return adc.rooms.ListAllowedRoomUsers(roomId)
}
func (adc *ActiveDataController) ResetAllowedRoomUsers(roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.ResetAllowedRoomUsers(roomId)
}

// User Room functions.

func (adc *ActiveDataController) PutUserIntoRoom(userId common.ObjectId, roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.PutUserIntoRoom(userId, roomId)
}
func (adc *ActiveDataController) RemoveUserFromRoom(userId common.ObjectId, roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.RemoveUserFromRoom(userId, roomId)
}
func (adc *ActiveDataController) GetUserRoomId(userId common.ObjectId) (roomId *common.ObjectId, rpcErr *jrm1.RpcError) {
	adc.guard.RLock()
	defer adc.guard.RUnlock()
	return adc.rooms.GetUserRoomId(userId)
}

// Message functions.

func (adc *ActiveDataController) AddMessageIntoRoom(roomId common.ObjectId, userId common.ObjectId, msgText string) (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.rooms.AddMessageIntoRoom(roomId, userId, msgText)
}
func (adc *ActiveDataController) ListAllMessagesInRoom(roomId common.ObjectId, userId common.ObjectId) (msgs []*msg.Message, rpcErr *jrm1.RpcError) {
	adc.guard.RLock()
	defer adc.guard.RUnlock()
	return adc.rooms.ListAllMessagesInRoom(roomId, userId)
}
func (adc *ActiveDataController) ListMessagesInRoomSince(roomId common.ObjectId, userId common.ObjectId, timeMarkTS int64) (msgs []*msg.Message, rpcErr *jrm1.RpcError) {
	adc.guard.RLock()
	defer adc.guard.RUnlock()
	return adc.rooms.ListMessagesInRoomSince(roomId, userId, timeMarkTS)
}

// Auxiliary functions.

func (adc *ActiveDataController) IsUserModerator(roomId common.ObjectId, userId common.ObjectId) (isModerator bool, rpcErr *jrm1.RpcError) {
	adc.guard.RLock()
	defer adc.guard.RUnlock()
	return adc.rooms.IsUserModerator(roomId, userId)
}
func (adc *ActiveDataController) GetServerStartTimeTS() (sstts int64) {
	return adc.serverStartTimeTS
}
func (adc *ActiveDataController) RemoveOutdatedSessions() (rpcErr *jrm1.RpcError) {
	adc.guard.Lock()
	defer adc.guard.Unlock()
	return adc.sessions.RemoveOutdatedSessions()
}
