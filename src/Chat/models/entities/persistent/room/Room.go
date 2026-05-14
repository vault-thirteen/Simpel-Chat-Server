package rm

import (
	"errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/Message"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/Messages"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/RoomParameters"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Room struct {
	common.MetaData
	Id   common.ObjectId `json:"id" gorm:"primarykey"`
	Type enum.RoomType   `json:"type" gorm:"column:type;type:tinyint"`
	Name string          `json:"name" gorm:"unique;size:255"`

	// IDs of users, who may change the list of allowed users of the chat room.
	// This field is shown to everyone.
	Moderators *common.IdList `json:"moderators" gorm:"column:moderators"`

	// IDs of users, who may read messages from and write messages into this
	// chat room. NULL for public rooms, as anyone can use them.
	// This field is shown only to moderators of the room.
	AllowedUserIds *common.IdList `json:"allowedUserIds,omitempty" gorm:"column:allowedUserIds"`

	// Common room settings.
	parameters        *rp.RoomParameters `gorm:"-"`
	serverStartTimeTS int64              `gorm:"-"`

	// Users who are actively interacting with the chat room, i.e. are reading
	// or writing messages. Before interacting with a room, a user must become
	// an active user of a room.
	activeUsers []*usr.User `gorm:"-"`

	// Message controller.
	messages *msgs.Messages `gorm:"-"`
}

func NewRoom(
	typé enum.RoomType,
	name string,
	roomParameters *rp.RoomParameters,
	serverStartTimeTS int64,
) (r *Room, err error) {
	r = &Room{
		//Id: 0, // ID is set by database.
		Type:           typé,
		Name:           name,
		Moderators:     nil,
		AllowedUserIds: nil,
	}

	r.initPrivateFields(roomParameters, serverStartTimeTS)

	err = r.Validate()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Room) InitPrivateFields(roomParameters *rp.RoomParameters, serverStartTimeTS int64) {
	r.initPrivateFields(roomParameters, serverStartTimeTS)
}
func (r *Room) initPrivateFields(roomParameters *rp.RoomParameters, serverStartTimeTS int64) {
	r.parameters = roomParameters
	r.serverStartTimeTS = serverStartTimeTS
	r.activeUsers = make([]*usr.User, 0)
	r.messages = msgs.NewMessages(roomParameters, serverStartTimeTS)
}

func (r *Room) Validate() (err error) {
	err = r.Type.Validate()
	if err != nil {
		return helper.NewError_InvalidEnumValue(enum.EnumField_RoomType, r.Type)
	}

	if len(r.Name) == 0 {
		return helper.NewError_ParameterIsNotSet("room name")
	}

	if len(r.Name) > r.parameters.RoomNameLengthLimit() {
		return helper.NewError_GenericError(helper.Err_NameIsTooLong, len(r.Name))
	}

	if r.parameters == nil {
		return helper.NewError_ParameterIsNotSet("room parameters")
	}

	if r.serverStartTimeTS <= 0 {
		return helper.NewError_ParameterIsNotSet("server start time")
	}

	if r.activeUsers == nil {
		return helper.NewError_ParameterIsNotSet("active users")
	}

	if r.messages == nil {
		return helper.NewError_ParameterIsNotSet("messages")
	}

	return nil
}

func (r *Room) IsUserModerator(userId common.ObjectId) bool {
	return r.Moderators.HasId(userId)
}
func (r *Room) IsUserAllowed(userId common.ObjectId) bool {
	return r.AllowedUserIds.HasId(userId)
}
func (r *Room) IsUserInsideRoom(userId common.ObjectId) bool {
	return r.hasUserWithId(userId)
}

func (r *Room) PutUserIntoRoom(user *usr.User) (err error) {
	if user == nil {
		return errors.New(helper.Err_NullPointer)
	}

	// Check for duplicate.
	if r.hasUserWithId(user.Id) {
		return errors.New(helper.Err_UserIsAlreadyInTheRoom)
	}

	switch r.Type {
	case enum.RoomType_Public:
		{
			r.addUser(user)
			return nil
		}

	case enum.RoomType_Private:
		{
			if !r.IsUserAllowed(user.Id) {
				return errors.New(helper.Err_UserIsNotAllowedToUseThisRoom)
			}
			r.addUser(user)
			return nil
		}

	default:
		return errors.New(helper.Err_UnknownRoomType)
	}
}
func (r *Room) RemoveUserFromRoom(user *usr.User) (err error) {
	if user == nil {
		return errors.New(helper.Err_NullPointer)
	}

	// User exists ?
	if !r.hasUserWithId(user.Id) {
		return errors.New(helper.Err_UserIsNotInTheRoom)
	}

	r.deleteUser(user)

	return nil
}
func (r *Room) GetFirstActiveUser() (user *usr.User) {
	if len(r.activeUsers) == 0 {
		return nil
	}

	return r.activeUsers[0]
}
func (r *Room) GetActiveUserIds() (activeUserIds []common.ObjectId) {
	if len(r.activeUsers) == 0 {
		return nil
	}

	activeUserIds = make([]common.ObjectId, 0, len(r.activeUsers))
	for _, user := range r.activeUsers {
		activeUserIds = append(activeUserIds, user.Id)
	}

	return activeUserIds
}

func (r *Room) AddMessage(userId common.ObjectId, msgText string) (err error) {
	if !r.hasUserWithId(userId) {
		return errors.New(helper.Err_UserIsNotInTheRoom)
	}

	m := msg.NewMessage(userId, msgText, r.serverStartTimeTS)

	err = r.messages.AddMessage(m)
	if err != nil {
		return err
	}

	return nil
}
func (r *Room) ListAllMessages(userId common.ObjectId) (msgs []*msg.Message, err error) {
	if !r.hasUserWithId(userId) {
		return nil, errors.New(helper.Err_UserIsNotInTheRoom)
	}

	return r.messages.GetAllMessages(), nil
}
func (r *Room) ListMessagesSince(userId common.ObjectId, timeMarkTS int64) (msgs []*msg.Message, err error) {
	if !r.hasUserWithId(userId) {
		return nil, errors.New(helper.Err_UserIsNotInTheRoom)
	}

	return r.messages.GetMessagesSince(timeMarkTS), nil
}

func (r *Room) hasUserWithId(userId common.ObjectId) bool {
	for _, user := range r.activeUsers {
		if user.Id == userId {
			return true
		}
	}

	return false
}
func (r *Room) addUser(user *usr.User) {
	r.activeUsers = append(r.activeUsers, user)
}
func (r *Room) deleteUser(user *usr.User) {
	// Find the index of the item to remove.
	var idx int
	var au *usr.User
	var isFound = false

	for idx, au = range r.activeUsers {
		if au.Id == user.Id {
			isFound = true
			break
		}
	}

	if !isFound {
		return
	}

	r.activeUsers = helper.ArrayWithoutItemAt(r.activeUsers, idx)
}

func (r *Room) GetMetaData() (md *common.MetaData) {
	return &r.MetaData
}
