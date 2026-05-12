package rs

import (
	"errors"
	"sort"
	"sync"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/database"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/der"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	usr "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
	rm "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/room"
	msg "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/Message"
	rp "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/RoomParameters"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	re "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/rpc/errors"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Rooms struct {
	guard              sync.RWMutex
	db                 *database.Database
	der                *der.DatabaseErrorReporter
	criticalErrorsChan *chan error

	roomsCountMax        int
	count                int
	commonRoomParameters *rp.RoomParameters
	serverStartTimeTS    int64

	// List of all rooms, accessed by room's ID.
	// N.B.: Null pointer to a room is not allowed here, i.e.
	// item removals are done via deletion of a map key.
	roomById map[common.ObjectId]*rm.Room

	// List of rooms used by each user, accessed by user's ID.
	// Each user may enter only a single room at the same time.
	// N.B.: Null pointer to a room is not allowed here, i.e.
	// item removals are done via deletion of a map key.
	roomByUserId map[common.ObjectId]*rm.Room
}

func NewRooms(
	db *database.Database,
	der *der.DatabaseErrorReporter,
	criticalErrorsChan *chan error,
	roomsCountMax int,
	commonRoomParameters *rp.RoomParameters,
	serverStartTimeTS int64,
) (r *Rooms, err error) {
	r = &Rooms{
		guard:              sync.RWMutex{},
		db:                 db,
		der:                der,
		criticalErrorsChan: criticalErrorsChan,

		roomsCountMax:        roomsCountMax,
		count:                0,
		commonRoomParameters: commonRoomParameters,
		serverStartTimeTS:    serverStartTimeTS,
		roomById:             make(map[common.ObjectId]*rm.Room),
		roomByUserId:         make(map[common.ObjectId]*rm.Room),
	}

	// As opposed to a normal room constructor which is used most of the time,
	// the initialisation of rooms at server's start is quite a complicated
	// process. Initial rooms' loading is performed in two steps:
	//  - the first step loads public data from database,
	//  - the second step prepares the rest data fields.

	err = r.loadRoomsFromDb()
	if err != nil {
		return nil, err
	}

	err = r.prepareRooms()
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Room functions.

func (r *Rooms) AddRoom(typé enum.RoomType, name string) (roomId common.ObjectId, rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	// Check limits.
	if r.count >= r.roomsCountMax {
		return 0, re.NewRpcError_RoomCountLimit(nil)
	}

	// As opposed to initial rooms' loading at server start,
	// here we use a normal room constructor.
	room, err := rm.NewRoom(typé, name, r.commonRoomParameters, r.serverStartTimeTS)
	if err != nil {
		return 0, re.NewRpcError_RoomError(err)
	}

	rpcErr = r.addRoom(room)
	if rpcErr != nil {
		return 0, rpcErr
	}

	rpcErr = r.internalSelfCheck()
	if rpcErr != nil {
		return 0, rpcErr
	}

	return room.Id, nil
}
func (r *Rooms) DeleteExistingRoom(roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	rpcErr = r.unbindRoom(roomId)
	if rpcErr != nil {
		return rpcErr
	}

	rpcErr = r.deleteRoom(&rm.Room{Id: roomId})
	if rpcErr != nil {
		return rpcErr
	}

	rpcErr = r.internalSelfCheck()
	if rpcErr != nil {
		return rpcErr
	}

	return nil
}
func (r *Rooms) ListRooms() (rooms []*rm.RoomForList, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	roomIds := make([]common.ObjectId, 0, len(r.roomById))
	for key, _ := range r.roomById {
		roomIds = append(roomIds, key)
	}

	sort.Slice(roomIds, func(i, j int) bool { return roomIds[i] < roomIds[j] })

	roomsSortedById := make([]*rm.RoomForList, 0, len(roomIds))
	for _, id := range roomIds {
		roomsSortedById = append(roomsSortedById, rm.NewRoomForList(r.roomById[id]))
	}

	return roomsSortedById, nil
}

// Room Moderator functions.

func (r *Rooms) AddRoomModerator(roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	if room.Moderators == nil {
		x := make(common.IdList, 0)
		room.Moderators = &x
	}

	err := room.Moderators.AddId(userId)
	if err != nil {
		return re.NewRpcError_RoomError(err)
	}

	err = r.db.SaveRoom(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	return nil
}
func (r *Rooms) DeleteRoomModerator(roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	if room.Moderators == nil {
		x := make(common.IdList, 0)
		room.Moderators = &x
	}

	err := room.Moderators.RemoveId(userId)
	if err != nil {
		return re.NewRpcError_RoomError(err)
	}

	err = r.db.SaveRoom(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	return nil
}
func (r *Rooms) ListRoomModerators(roomId common.ObjectId) (userIds []common.ObjectId, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return nil, re.NewRpcError_RoomError(err)
	}

	userIds = room.Moderators.List()

	return userIds, nil
}
func (r *Rooms) ResetRoomModerators(roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	room.Moderators = nil

	err := r.db.ResetRoomModerators(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	return nil
}

// Allowed Room User functions.

func (r *Rooms) AddAllowedRoomUser(callerId common.ObjectId, roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	if !room.IsUserModerator(callerId) {
		err := errors.New(helper.Err_YouAreNotModerator)
		return re.NewRpcError_RoomError(err)
	}

	if room.Type == enum.RoomType_Public {
		err := errors.New(helper.Err_PublicRoomCanNotHaveAllowedUsers)
		return re.NewRpcError_RoomError(err)
	}

	if room.AllowedUserIds == nil {
		x := make(common.IdList, 0)
		room.AllowedUserIds = &x
	}

	err := room.AllowedUserIds.AddId(userId)
	if err != nil {
		return re.NewRpcError_RoomError(err)
	}

	err = r.db.SaveRoom(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	return nil
}
func (r *Rooms) DeleteAllowedRoomUser(callerId common.ObjectId, roomId common.ObjectId, userId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	if !room.IsUserModerator(callerId) {
		err := errors.New(helper.Err_YouAreNotModerator)
		return re.NewRpcError_RoomError(err)
	}

	if room.Type == enum.RoomType_Public {
		err := errors.New(helper.Err_PublicRoomCanNotHaveAllowedUsers)
		return re.NewRpcError_RoomError(err)
	}

	if room.AllowedUserIds == nil {
		x := make(common.IdList, 0)
		room.AllowedUserIds = &x
	}

	err := room.AllowedUserIds.RemoveId(userId)
	if err != nil {
		return re.NewRpcError_RoomError(err)
	}

	err = r.db.SaveRoom(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	return nil
}
func (r *Rooms) ListAllowedRoomUsers(callerId common.ObjectId, roomId common.ObjectId) (userIds []common.ObjectId, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return nil, re.NewRpcError_RoomError(err)
	}

	if !room.IsUserModerator(callerId) {
		err := errors.New(helper.Err_YouAreNotModerator)
		return nil, re.NewRpcError_RoomError(err)
	}

	userIds = room.AllowedUserIds.List()

	return userIds, nil
}
func (r *Rooms) ResetAllowedRoomUsers(callerId common.ObjectId, roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	if !room.IsUserModerator(callerId) {
		err := errors.New(helper.Err_YouAreNotModerator)
		return re.NewRpcError_RoomError(err)
	}

	room.AllowedUserIds = nil

	err := r.db.ResetAllowedRoomUsers(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	return nil
}

// User Room functions.

func (r *Rooms) PutUserIntoRoom(userId common.ObjectId, roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	var room *rm.Room

	// Is user already using another room ? Find the room. Check the room.
	{
		_, userHasRoom := r.roomByUserId[userId]
		if userHasRoom {
			return re.NewRpcError_UserCanNotUseMultipleRooms(nil)
		}

		var ok bool
		room, ok = r.roomById[roomId]
		if !ok {
			return re.NewRpcError_RoomDoesNotExist(nil)
		}

		if room.Id != roomId {
			return re.NewRpcError_RoomIsNotFound(nil)
		}

		switch room.Type {
		case enum.RoomType_Private:
			if !room.IsUserAllowed(userId) {
				return re.NewRpcError_Msg_UserIsNotAllowedInTheRoom(nil)
			}
		}
	}

	var user *usr.User

	// Find the user.
	{
		user = &usr.User{Id: userId}
		err := r.db.GetUserById(user)
		if err != nil {
			return r.der.DatabaseError(err)
		}

		if user.Id != userId {
			r.der.DatabaseError(nil)
		}
	}

	// Put user into room.
	{
		err := room.PutUserIntoRoom(user)
		if err != nil {
			return re.NewRpcError_RoomError(err)
		}

		r.bindUserWithRoom(room, user.Id)
	}

	return nil
}
func (r *Rooms) RemoveUserFromRoom(userId common.ObjectId, roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	var room *rm.Room

	// Is user using any room ? Find the room. Check the room.
	{
		var userHasRoom bool
		room, userHasRoom = r.roomByUserId[userId]
		if !userHasRoom {
			return re.NewRpcError_UserIsNotUsingAnyRoom(nil)
		}

		if room.Id != roomId {
			return re.NewRpcError_RoomIsNotFound(nil)
		}
	}

	var user *usr.User

	// Find the user.
	{
		user = &usr.User{Id: userId}
		err := r.db.GetUserById(user)
		if err != nil {
			return r.der.DatabaseError(err)
		}

		if user.Id != userId {
			r.der.DatabaseError(nil)
		}
	}

	// Remove user from room.
	{
		err := room.RemoveUserFromRoom(user)
		if err != nil {
			return re.NewRpcError_RoomError(err)
		}

		r.unbindUserFromRoom(user.Id)
	}

	return nil
}
func (r *Rooms) GetUserRoomId(userId common.ObjectId) (roomId *common.ObjectId, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, userHasRoom := r.roomByUserId[userId]
	if !userHasRoom {
		return nil, nil
	}

	return &room.Id, nil
}
func (r *Rooms) GetRoom(callerId common.ObjectId, roomId common.ObjectId) (room *rm.Room, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, ok := r.roomById[roomId]
	if !ok {
		return nil, nil
	}

	if !room.IsUserModerator(callerId) {
		room.AllowedUserIds = nil
	}

	return room, nil
}
func (r *Rooms) GetRoomUsers(roomId common.ObjectId) (activeUserIds []common.ObjectId, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, ok := r.roomById[roomId]
	if !ok {
		return nil, nil
	}

	return room.GetActiveUserIds(), nil
}

// Message functions.

func (r *Rooms) AddMessageIntoRoom(roomId common.ObjectId, userId common.ObjectId, msgText string) (rpcErr *jrm1.RpcError) {
	r.guard.Lock()
	defer r.guard.Unlock()

	if len(msgText) > r.commonRoomParameters.MessageSizeLimit() {
		err := errors.New(helper.Err_MessageIsTooLong)
		return re.NewRpcError_RoomError(err)
	}

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	err := room.AddMessage(userId, msgText)
	if err != nil {
		return re.NewRpcError_RoomError(err)
	}

	return nil
}
func (r *Rooms) ListAllMessagesInRoom(roomId common.ObjectId, userId common.ObjectId) (msgs []*msg.Message, nowTS int64, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return nil, -1, re.NewRpcError_RoomError(err)
	}

	// We return a timestamp of the memory read operation for cases of slow
	// network communication in order for the client to be able to adjust its
	// behaviour to these slow network conditions.
	nowTS = time.Now().Unix()

	var err error
	msgs, err = room.ListAllMessages(userId)
	if err != nil {
		return nil, -1, re.NewRpcError_RoomError(err)
	}

	return msgs, nowTS, nil
}
func (r *Rooms) ListMessagesInRoomSince(roomId common.ObjectId, userId common.ObjectId, timeMarkTS int64) (msgs []*msg.Message, nowTS int64, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return nil, -1, re.NewRpcError_RoomError(err)
	}

	// We return a timestamp of the memory read operation for cases of slow
	// network communication in order for the client to be able to adjust its
	// behaviour to these slow network conditions.
	nowTS = time.Now().Unix()

	var err error
	msgs, err = room.ListMessagesSince(userId, timeMarkTS)
	if err != nil {
		return nil, -1, re.NewRpcError_RoomError(err)
	}

	return msgs, nowTS, nil
}

// Auxiliary functions.

func (r *Rooms) IsUserModerator(roomId common.ObjectId, userId common.ObjectId) (isModerator bool, rpcErr *jrm1.RpcError) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return false, re.NewRpcError_RoomError(err)
	}

	isModerator = room.IsUserModerator(userId)

	return isModerator, nil
}

func (r *Rooms) loadRoomsFromDb() (err error) {
	// In this first step we load public fields from database:
	// - ID,
	// - Type,
	// - Name,
	// - Moderators,
	// - AllowedUserIds.

	if len(r.roomById) != 0 {
		return errors.New(helper.Err_RoomsListIsAlreadyLoaded)
	}

	var rooms []*rm.Room
	rooms, err = r.db.ListAllRooms()
	if err != nil {
		return err
	}

	for _, room := range rooms {
		_, isDuplicate := r.roomById[room.Id]
		if isDuplicate {
			return errors.New(helper.Err_DuplicateRoom)
		}

		r.roomById[room.Id] = room
	}

	r.count = len(r.roomById)

	return nil
}
func (r *Rooms) prepareRooms() (err error) {
	// In this second step we prepare private fields (Go language calls them
	// "unexported") which are not saved in database:
	// - parameters,
	// - serverStartTimeTS,
	// - activeUsers,
	// - messages.

	for _, room := range r.roomById {
		room.InitPrivateFields(r.commonRoomParameters, r.serverStartTimeTS)
	}

	r.count = len(r.roomById)

	return nil
}
func (r *Rooms) addRoom(room *rm.Room) (rpcErr *jrm1.RpcError) {
	err := r.db.CreateRoom(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	_, isDuplicate := r.roomById[room.Id]
	if isDuplicate {
		err = errors.New(helper.Err_DuplicateRoom)

		defer func() {
			*(r.criticalErrorsChan) <- helper.NewError_WrappedError(helper.Err_ADC, err)
		}()

		return re.NewRpcError_RoomError(err)
	}

	r.roomById[room.Id] = room
	r.count++

	return nil
}
func (r *Rooms) deleteRoom(room *rm.Room) (rpcErr *jrm1.RpcError) {
	err := r.db.DeleteRoom(room)
	if err != nil {
		return r.der.DatabaseError(err)
	}

	delete(r.roomById, room.Id)
	r.count--

	return nil
}
func (r *Rooms) unbindRoom(roomId common.ObjectId) (rpcErr *jrm1.RpcError) {
	room, ok := r.roomById[roomId]
	if !ok {
		err := errors.New(helper.Err_RoomIsNotFound)
		return re.NewRpcError_RoomError(err)
	}

	var activeUser *usr.User
	for {
		activeUser = room.GetFirstActiveUser()
		if activeUser == nil {
			break
		}

		err := room.RemoveUserFromRoom(activeUser)
		if err != nil {
			return re.NewRpcError_RoomError(err)
		}

		r.unbindUserFromRoom(activeUser.Id)
	}

	return nil
}

func (r *Rooms) internalSelfCheck() (rpcErr *jrm1.RpcError) {
	// 1. Compare room count in memory.
	if r.count != len(r.roomById) {
		defer func() {
			*(r.criticalErrorsChan) <- helper.NewError_WrappedError(helper.Err_ADC, errors.New(helper.Err_RoomCountMismatch))
		}()

		return re.NewRpcError_ActiveDataController(nil)
	}

	// 2. Compare room count in database and memory.
	roomsCountInDb, err := r.db.CountAllRooms()
	if err != nil {
		return r.der.DatabaseError(err)
	}

	if roomsCountInDb != r.count {
		defer func() {
			*(r.criticalErrorsChan) <- helper.NewError_WrappedError(helper.Err_ADC, errors.New(helper.Err_RoomCountMismatch))
		}()

		return re.NewRpcError_ActiveDataController(nil)
	}

	return nil
}
func (r *Rooms) bindUserWithRoom(room *rm.Room, userId common.ObjectId) {
	r.roomByUserId[userId] = room
}
func (r *Rooms) unbindUserFromRoom(userId common.ObjectId) {
	delete(r.roomByUserId, userId)
}
