package libp2p

import (
	"context"
	"crypto/ed25519"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	log "github.com/sirupsen/logrus"
)

type LibP2PClient struct {
	h host.Host
}

func (p *LibP2PClient) Listen(receiver common.MessageReceiver) error {
	return nil
}

func (c *LibP2PClient) Send(msg *common.Message) error {
	return nil
}

func (p *LibP2PClient) Close() error {
	return nil
}

// TODO implement Messenger
func NewLibP2PClient(key ed25519.PrivateKey) *LibP2PClient { // TODO add more params
	libp2pKey, err := crypto.UnmarshalEd25519PrivateKey(key)
	if err != nil {
		panic(err)
	}

	h, err := libp2p.New(context.Background(), libp2p.Identity(libp2pKey))
	if err != nil {
		panic(err)
	}

	peerInfo := peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	log.Infof("[LIBP2P] Started libp2p host: %s\n", addrs[0])

	return &LibP2PClient{h: h}

}
