package queue

import (
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type MessageQueue struct {
	messageChannel chan *common.Message
}

func (q *MessageQueue) Receive(msg *common.Message) {
	q.messageChannel <- msg
}

func New(ch chan *common.Message, handler common.ChannelMessageHandler) *MessageQueue {
	q := MessageQueue{messageChannel: ch}

	go handler.Handle(q.messageChannel)
	return &q
}
