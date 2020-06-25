package json

import (
	"encoding/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type JSONBusinessMesssageParser struct{}

func (p *JSONBusinessMesssageParser) Parse(msg *common.Message) (*messages.BusinessMessage, error) {
	var res messages.BusinessMessage
	err := json.Unmarshal(msg.Msg, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
