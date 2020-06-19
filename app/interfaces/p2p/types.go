package p2p

import (
	"context"
)

type P2PMessage struct {
	Msg []byte // The actual raw message

	Ctx context.Context // Hold context variables like who has sent it if needed for the future
	// Read more https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39
}

type MessageReceiver interface {
	Receive(msg *P2PMessage)
}

type ChannelMessageHandler interface {
	Handle(ch <-chan *P2PMessage)
}

type Messenger interface {
	Connect(receiver MessageReceiver) error

	Send(msg P2PMessage) error

	Close() error
}
