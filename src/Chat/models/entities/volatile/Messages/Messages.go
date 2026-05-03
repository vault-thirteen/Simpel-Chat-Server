package msgs

import (
	"errors"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/Message"
	rp "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/RoomParameters"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Messages struct {
	parameters        *rp.RoomParameters `gorm:"-"`
	serverStartTimeTS int64              `gorm:"-"`
	messages          []*msg.Message     `gorm:"-"`
}

func NewMessages(parameters *rp.RoomParameters, serverStartTimeTS int64) (m *Messages) {
	if (parameters == nil) || (serverStartTimeTS <= 0) {
		return nil
	}

	m = &Messages{
		parameters:        parameters,
		serverStartTimeTS: serverStartTimeTS,
	}

	m.resetMessages()

	return m
}

func (m *Messages) resetMessages() {
	m.messages = make([]*msg.Message, 0, m.parameters.MessageCountLimit())
}

func (m *Messages) AddMessage(message *msg.Message) (err error) {
	// Checks.
	// N.B. Other checks are performed by message constructor.
	{
		if message == nil {
			return errors.New(helper.Err_NullPointer)
		}

		if len(message.Content()) > m.parameters.MessageSizeLimit() {
			return errors.New(helper.Err_MessageIsTooLong)
		}
	}

	// Reset messages on overflow.
	if len(m.messages) >= m.parameters.MessageCountLimit() {
		m.resetMessages()
	}

	m.messages = append(m.messages, message)

	return nil
}

func (m *Messages) GetMessagesSince(timeBorderTS int64) (out []*msg.Message) {
	out = make([]*msg.Message, 0, m.parameters.MessageCountLimit())

	var mtoc int64
	for _, message := range m.messages {
		mtoc = m.serverStartTimeTS + message.TimeOfCreationDTS()

		if mtoc < timeBorderTS {
			continue
		}

		out = append(out, message)
	}

	return out
}

func (m *Messages) GetAllMessages() []*msg.Message {
	return m.messages
}

func (m *Messages) MessagesCount() int {
	return len(m.messages)
}
