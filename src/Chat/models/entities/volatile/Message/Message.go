package msg

import (
	"time"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
)

type Message struct {
	authorId common.ObjectId `gorm:"-"`
	content  string          `gorm:"-"`

	// Delta Unix Timestamp, since Server's Start Time.
	timeOfCreationDTS int64 `gorm:"-"`
}

func NewMessage(authorId common.ObjectId, content string, serverStartTimeTS int64) (m *Message) {
	// N.B. While Message object has no access to size limits, message's size
	// limit is controlled by Message Controller a.k.a Messages object.

	if (authorId == 0) ||
		(len(content) == 0) ||
		(serverStartTimeTS <= 0) {
		return nil
	}

	nowTS := time.Now().UTC().Unix()

	m = &Message{
		authorId:          authorId,
		content:           content,
		timeOfCreationDTS: nowTS - serverStartTimeTS,
	}

	return m
}

func (m *Message) AuthorId() common.ObjectId {
	return m.authorId
}
func (m *Message) Content() string {
	return m.content
}
func (m *Message) TimeOfCreationDTS() int64 {
	return m.timeOfCreationDTS
}
