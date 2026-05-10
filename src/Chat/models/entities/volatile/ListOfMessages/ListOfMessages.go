package lom

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	msg "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/volatile/Message"
)

type ListOfMessages struct {
	OpTimeTS          int64             `json:"opTimeTS"`
	RoomId            common.ObjectId   `json:"roomId"`
	ServerStartTimeTS int64             `json:"SSTTS"`
	SinceTS           *int64            `json:"sinceTS,omitempty"`
	Count             int               `json:"count"`
	AuthorIds         []common.ObjectId `json:"authorIds"`
	TOCDTSs           []int64           `json:"TOCDTSs"`
	Contents          []string          `json:"contents"`
}

func NewListOfMessages(roomId common.ObjectId, msgs []*msg.Message, nowTS int64, serverStartTimeTS int64, sinceTS *int64) (l *ListOfMessages) {
	if msgs == nil {
		return nil
	}

	l = &ListOfMessages{
		OpTimeTS:          nowTS,
		RoomId:            roomId,
		ServerStartTimeTS: serverStartTimeTS,
		SinceTS:           sinceTS,
		Count:             len(msgs),
		AuthorIds:         make([]common.ObjectId, 0, len(msgs)),
		TOCDTSs:           make([]int64, 0, len(msgs)),
		Contents:          make([]string, 0, len(msgs)),
	}

	for _, m := range msgs {
		l.AuthorIds = append(l.AuthorIds, m.AuthorId())
		l.TOCDTSs = append(l.TOCDTSs, m.TimeOfCreationDTS())
		l.Contents = append(l.Contents, m.Content())
	}

	return l
}
