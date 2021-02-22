package common

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

type _Message struct {
	Msg []byte // The actual raw message

	Ctx context.Context // Hold context variables like who has sent it if needed for the future
	// Read more https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39
}

type Message struct {
	Msg []byte // The actual raw message

	Ctx context.Context // Hold context variables like who has sent it if needed for the future
	// Read more https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39
}

type SignedMessage struct {
	Msg        []byte
	Signature  []byte
	PeerId     peer.ID
	PubKeyData []byte
}

type MessageReceiver interface {
	Receive(msg *Message)
}

type ChannelMessageHandler interface {
	Handle(ch <-chan *Message)
}

type Messenger interface {
	Listen(receiver MessageReceiver) error

	Send(msg *Message) error

	Close() error
}
