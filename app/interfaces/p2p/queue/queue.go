package queue

import (
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
)

type MessageQueue struct {
	messageChannel chan *p2p.P2PMessage
}

func (q *MessageQueue) Receive(msg *p2p.P2PMessage) {
	q.messageChannel <- msg
}

func New(ch chan *p2p.P2PMessage, handler p2p.ChannelMessageHandler) *MessageQueue {
	q := MessageQueue{messageChannel: ch}

	go handler.Handle(q.messageChannel)
	return &q
}
