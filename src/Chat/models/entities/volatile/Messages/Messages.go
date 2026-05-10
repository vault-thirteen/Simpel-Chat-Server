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
	var timeBorderDTS = timeBorderTS - m.serverStartTimeTS

	if len(m.messages) < 1 {
		return out
	}

	// Read messages from the left or right corner of array depending on the
	// requested time border.

	firstMsgTimeDTS := m.messages[0].TimeOfCreationDTS()
	lastMsgTimeDTS := m.messages[len(m.messages)-1].TimeOfCreationDTS()
	avgMsgTimeDTS := (lastMsgTimeDTS - firstMsgTimeDTS) / 2

	if timeBorderDTS <= avgMsgTimeDTS {
		// Read messages from the left end.
		out = make([]*msg.Message, 0, m.parameters.MessageCountLimit())

		for _, message := range m.messages {
			if message.TimeOfCreationDTS() < timeBorderDTS {
				continue
			}

			out = append(out, message)
		}
	} else {
		// Read messages from the right end.
		buf := make([]*msg.Message, 0, m.parameters.MessageCountLimit())

		i := len(m.messages) - 1
		message := m.messages[i]
		for message.TimeOfCreationDTS() >= timeBorderDTS {
			buf = append(buf, message)

			// Next left index.
			i--
			if i < 0 {
				break
			}
			message = m.messages[i]
		}

		out = make([]*msg.Message, 0, len(buf))
		for i = len(buf) - 1; i >= 0; i-- {
			out = append(out, buf[i])
		}
	}

	return out
}

func (m *Messages) GetAllMessages() []*msg.Message {
	return m.messages
}

func (m *Messages) MessagesCount() int {
	return len(m.messages)
}
