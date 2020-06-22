package json

import (
	"context"
	"encoding/json"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
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

	var parser JSONBusinessMesssageParser

	msg, err := parser.Parse(&reqMsg)

	assert.Nil(t, err, "The parser could not parse the p2p message")
	assert.NotNil(t, msg, "The resulting message was")
	assert.Equal(t, msg.Type, "proposal", "The parser parsed the type incorrectly")

}
