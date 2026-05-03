package rp

type RoomParameters struct {
	messageCountMax int `gorm:"-"`
	messageSizeMax  int `gorm:"-"`
}

func NewRoomParameters(messageCountMax, messageSizeMax int) *RoomParameters {
	return &RoomParameters{
		messageCountMax: messageCountMax,
		messageSizeMax:  messageSizeMax,
	}
}

func (rp *RoomParameters) MessageCountLimit() int {
	return rp.messageCountMax
}

func (rp *RoomParameters) MessageSizeLimit() int {
	return rp.messageSizeMax
}
