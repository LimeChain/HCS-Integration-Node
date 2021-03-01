package apiservices

import (
	"errors"

	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type NodeService struct {
	messenger common.Messenger
}

func (s *NodeService) Connect(peerAddress string) (bool, error) {
	if len(peerAddress) == 0 {
		return false, errors.New("Missing peer address information")
	}

	connected, err := s.messenger.Connect(peerAddress)

	if err != nil {
		return false, err
	}

	return connected, nil
}

func NewNodeService(messenger common.Messenger) *NodeService {
	return &NodeService{messenger: messenger}
}
