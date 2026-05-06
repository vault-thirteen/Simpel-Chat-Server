package rp

type RoomParameters struct {
	messageCountMax   int `gorm:"-"`
	messageSizeMax    int `gorm:"-"`
	roomNameLengthMax int `gorm:"-"`
}

func NewRoomParameters(messageCountMax, messageSizeMax, roomNameLengthMax int) *RoomParameters {
	return &RoomParameters{
		messageCountMax:   messageCountMax,
		messageSizeMax:    messageSizeMax,
		roomNameLengthMax: roomNameLengthMax,
	}
}

func (rp *RoomParameters) MessageCountLimit() int {
	return rp.messageCountMax
}
func (rp *RoomParameters) MessageSizeLimit() int {
	return rp.messageSizeMax
}
func (rp *RoomParameters) RoomNameLengthLimit() int { return rp.roomNameLengthMax }
