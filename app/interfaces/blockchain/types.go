package blockchain

import (
	"context"
)

type BlockchainMessage struct {
	Msg []byte // The actual raw message

	Ctx context.Context
}

type BlockchainMessageReceiver interface {
	Receive(msg *BlockchainMessage)
}

type BlockchainMessageHandler interface {
	Handle(ch <-chan *BlockchainMessage)
}

type BlockchainClient interface {
	Listen(receiver BlockchainMessageReceiver) error

	Send(msg BlockchainMessage) error
}
