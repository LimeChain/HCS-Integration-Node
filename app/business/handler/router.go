package handler

import (
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
)

type BusinessMessageRouter struct {
	parser   BusinessMessageParser
	handlers map[string]BusinessLogicHandler
}

// Through this implements the ChannelMessageHandler interface needed for the MessageQueue
func (mr *BusinessMessageRouter) Handle(ch <-chan *p2p.P2PMessage) {

	// Waits for event
	for msg := range ch {

		// Parses event type and passes it to the correct BusinessLogicHandler based on the type
		_, err := mr.parser.Parse(msg)
		if err != nil {
			panic(err)
		}

		handler := mr.handlers[string(msg.Msg)]
		handler.Handle(msg)
	}
}

func (mr *BusinessMessageRouter) AddHandler(messageType string, handler BusinessLogicHandler) {
	mr.handlers[messageType] = handler
}

func NewBusinessMessageRouter(parser BusinessMessageParser) *BusinessMessageRouter {
	handlers := make(map[string]BusinessLogicHandler)
	return &BusinessMessageRouter{parser, handlers}
}
