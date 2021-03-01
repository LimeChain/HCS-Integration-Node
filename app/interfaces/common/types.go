package common

import (
	"context"
)

type Message struct {
	Msg []byte // The actual raw message

	Ctx context.Context // Hold context variables like who has sent it if needed for the future
	// Read more https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39
}

type Envelope struct {
	Payload    []byte
	Signature  []byte
	PeerId     string
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
