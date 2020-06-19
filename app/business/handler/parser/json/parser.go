package json

import (
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
)

type JSONBusinessMesssageParser struct{}

func (p *JSONBusinessMesssageParser) Parse(msg *p2p.P2PMessage) (*handler.BusinessMessage, error) {
	fmt.Println("Parser", msg.Msg)
	return nil, nil
}
