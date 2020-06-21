package handler

import (
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"log"
)

type BusinessMessageRouter struct {
	parser   BusinessMessageParser
	handlers map[string]BusinessLogicHandler
}

// Through this implements the ChannelMessageHandler interface needed for the MessageQueue
func (mr *BusinessMessageRouter) Handle(ch <-chan *p2p.P2PMessage) {

	// Waits for event
	for msg := range ch {
		if err := mr.handleMessage(msg); err != nil {
			log.Println(err.Error())
		}
	}
}

func (mr *BusinessMessageRouter) handleMessage(msg *p2p.P2PMessage) error {
	// Parses event type and passes it to the correct BusinessLogicHandler based on the type
	bMsg, err := mr.parser.Parse(msg)
	if err != nil {
		return err
	}

	handler := mr.handlers[bMsg.Type]
	if err := handler.Handle(msg); err != nil {
		return err
	}

	return nil
}

func (mr *BusinessMessageRouter) AddHandler(messageType string, handler BusinessLogicHandler) {
	mr.handlers[messageType] = handler
}

func NewBusinessMessageRouter(parser BusinessMessageParser) *BusinessMessageRouter {
	handlers := make(map[string]BusinessLogicHandler)
	return &BusinessMessageRouter{parser, handlers}
}
