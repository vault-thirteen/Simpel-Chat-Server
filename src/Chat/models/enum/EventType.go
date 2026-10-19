package enum

const (
	EventType_UserRegistration        = 1
	EventType_UserBan                 = 2
	EventType_UserLogIn               = 3
	EventType_UserLogOut              = 4
	EventType_SessionTimeOut          = 5
	EventType_UserChangePassword      = 6
	EventType_RoomCreation            = 7
	EventType_RoomDeletion            = 8
	EventType_RoomModeratorAddition   = 9
	EventType_RoomModeratorDeletion   = 10
	EventType_RoomModeratorsReset     = 11
	EventType_RoomAllowedUserAddition = 12
	EventType_RoomAllowedUserDeletion = 13
	EventType_RoomAllowedUserReset    = 14
)

type EventType byte

func (et EventType) IsValid() bool {
	switch et {
	case EventType_UserRegistration:
		return true
	case EventType_UserBan:
		return true
	case EventType_UserLogIn:
		return true
	case EventType_UserLogOut:
		return true
	case EventType_SessionTimeOut:
		return true
	case EventType_UserChangePassword:
		return true
	case EventType_RoomCreation:
		return true
	case EventType_RoomDeletion:
		return true
	case EventType_RoomModeratorAddition:
		return true
	case EventType_RoomModeratorDeletion:
		return true
	case EventType_RoomModeratorsReset:
		return true
	case EventType_RoomAllowedUserAddition:
		return true
	case EventType_RoomAllowedUserDeletion:
		return true
	case EventType_RoomAllowedUserReset:
		return true
	}
	return false
}
