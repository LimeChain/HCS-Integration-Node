package queue

import (
	"context"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockChannelMessageHandler struct{}

func (mcm *MockChannelMessageHandler) Handle(ch <-chan *p2p.P2PMessage) {

}

func TestNew(t *testing.T) {

	ch := make(chan *p2p.P2PMessage)

	var fr MockChannelMessageHandler

	q := New(ch, &fr)

	assert.Equal(t, q.messageChannel, ch, "New did not store the correct channel")
}

func TestReceive(t *testing.T) {

	ch := make(chan *p2p.P2PMessage)

	var fr MockChannelMessageHandler

	q := New(ch, &fr)

	go func() {
		q.Receive(&p2p.P2PMessage{Msg: []byte("rfp"), Ctx: context.Background()})
		q.Receive(&p2p.P2PMessage{Msg: []byte("rf2"), Ctx: context.Background()})
	}()

	res1, res2 := <-ch, <-ch
	assert.Equal(t, string(res1.Msg), "rfp", "The first element of the queue did not match what was sent")
	assert.Equal(t, string(res2.Msg), "rf2", "The second element of the queue did not match what was sent")
}
