package libp2p

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const p2pStreamName = "/hcs-int-p2p-nodes/1.0.0"

type LibP2PClient struct {
	h                  host.Host
	messagesReadWriter *bufio.ReadWriter
}

func handleIncommingMessage(messagesReadWriter *bufio.ReadWriter, receiver common.MessageReceiver) {
	go func() {
		for {
			msg, err := messagesReadWriter.ReadBytes('\n')
			if err != nil {
				panic(err)
			}

			receiver.Receive(&common.Message{Ctx: context.Background(), Msg: msg})
		}
	}()
}

func (c *LibP2PClient) Listen(receiver common.MessageReceiver) error {
	if c.messagesReadWriter != nil { // I've started the stream and have readwriter available
		handleIncommingMessage(c.messagesReadWriter, receiver)
		return nil
	}

	c.h.SetStreamHandler(p2pStreamName, func(s network.Stream) { // I'm waiting for incomming connection
		c.messagesReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		handleIncommingMessage(c.messagesReadWriter, receiver)
	})
	return nil
}

func (c *LibP2PClient) Send(msg *common.Message) error {
	c.messagesReadWriter.Write(msg.Msg)
	c.messagesReadWriter.WriteByte('\n')
	c.messagesReadWriter.Flush()
	return nil
}

func (c *LibP2PClient) Close() error {
	return nil
}

func NewLibP2PClient(key ed25519.PrivateKey, listenIp, listenPort, peerMultiAddr string) *LibP2PClient {
	libp2pKey, err := crypto.UnmarshalEd25519PrivateKey(key)
	if err != nil {
		panic(err)
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", listenIp, listenPort))

	h, err := libp2p.New(context.Background(), libp2p.ListenAddrs(sourceMultiAddr), libp2p.Identity(libp2pKey))
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

	log.Infof("[LIBP2P] Started libp2p host and listening on: %s \n", addrs[0])

	client := &LibP2PClient{h: h}

	if len(peerMultiAddr) == 0 {
		return client
	}

	maddr, err := multiaddr.NewMultiaddr(peerMultiAddr)
	if err != nil {
		log.Fatalln(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}

	err = h.Connect(context.Background(), *info)
	if err != nil {
		panic(err)
	}

	log.Infof("[LIBP2P] Connected to peer: %s\n", peerMultiAddr)

	s, err := h.NewStream(context.Background(), info.ID, p2pStreamName)
	if err != nil {
		panic(err)
	}

	client.messagesReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	return client

}
