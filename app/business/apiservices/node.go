package apiservices

import (
	"errors"

	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/messaging/libp2p"
)

type NodeService struct {
	p2pClient *libp2p.LibP2PClient
}

func (s *NodeService) Connect(peerAddress string) error {
	if len(peerAddress) == 0 {
		return errors.New("Missing peer address information")
	}

	targetPeerInfo, err := libp2p.MultiAddrToPeerInfo(peerAddress)
	if err != nil {
		return err
	}

	err = libp2p.Connect(s.p2pClient, *targetPeerInfo)

	if err != nil {
		return err
	}

	return nil
}

func NewNodeService(p2pClient *libp2p.LibP2PClient) *NodeService {
	return &NodeService{p2pClient: p2pClient}
}
