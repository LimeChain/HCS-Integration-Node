package router

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockParser struct {
	Type string
	Err  error
}

func (mp *MockParser) Parse(msg *p2p.P2PMessage) (*messages.BusinessMessage, error) {
	return &messages.BusinessMessage{Type: mp.Type}, mp.Err
}

type MockHandler struct {
	h func() error
}

func (m *MockHandler) Handle(msg *p2p.P2PMessage) error {
	return m.h()
}

func NewMockHandler(h func() error) *MockHandler {
	return &MockHandler{
		h: h,
	}
}

func TestNewBusinessMessageRouter(t *testing.T) {

	var parser MockParser

	r := NewBusinessMessageRouter(&parser)

	assert.Equal(t, r.parser, &parser, "It did not save the correct parser")
	assert.Empty(t, r.handlers, "It did have some built in handlers")

}

func TestAddHandler(t *testing.T) {

	var parser MockParser

	r := NewBusinessMessageRouter(&parser)

	var handler MockHandler

	r.AddHandler("proposal", &handler)
	assert.Equal(t, r.handlers["proposal"], &handler, "It did not set handler correctly")
	assert.Equal(t, len(r.handlers), 1, "It did have some built in handlers")

}

func TestHandle(t *testing.T) {

	okParser := MockParser{
		Type: "proposal",
		Err:  nil,
	}

	parserError := errors.New("Test Parse Error")

	badParser := MockParser{
		Type: "",
		Err:  parserError,
	}

	handlerCalled := false

	okHandler := NewMockHandler(func() error {
		handlerCalled = true
		return nil
	})

	handlerError := errors.New("Test Handler Error")

	badHandler := NewMockHandler(func() error {
		handlerCalled = true
		return handlerError
	})

	cases := []struct {
		testName            string
		parser              BusinessMessageParser
		handler             *MockHandler
		expectHandlerCalled bool
		expectedError       error
	}{
		{"Should Work Ok", &okParser, okHandler, true, nil},
		{"Should throw error on parser problem", &badParser, okHandler, false, parserError},
		{"Should throw error on handler problem", &okParser, badHandler, true, handlerError},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			handlerCalled = false
			r := NewBusinessMessageRouter(tc.parser)
			r.AddHandler("proposal", tc.handler)

			req := struct {
				Type string
			}{
				Type: "proposal",
			}

			reqBytes, err := json.Marshal(req)
			assert.Nil(t, err, "could not marshal the request")

			reqMsg := p2p.P2PMessage{
				Msg: reqBytes,
				Ctx: context.Background(),
			}

			err = r.handleMessage(&reqMsg)

			assert.Equal(t, handlerCalled, tc.expectHandlerCalled)
			assert.Equal(t, err, tc.expectedError)
		})
	}

}
