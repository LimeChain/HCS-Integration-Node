package libp2p

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const p2pStreamName = "/hcs-int-p2p-nodes/1.0.0"

type LibP2PClient struct {
	h                  host.Host
	messagesReadWriter *bufio.ReadWriter
	receiver           common.MessageReceiver
	streamPairs        map[peer.ID]string
}

func handleIncommingMessage(c *LibP2PClient, receiver common.MessageReceiver) {
	go func() {
		for {
			msg, err := c.messagesReadWriter.ReadBytes('\n')
			if err != nil {
				c.streamPairs[c.h.ID()] = ""
				return
			}

			receiver.Receive(&common.Message{Ctx: context.Background(), Msg: msg})
		}
	}()
}

func (c *LibP2PClient) Listen(receiver common.MessageReceiver) error {
	c.receiver = receiver

	c.h.SetStreamHandler(p2pStreamName, func(s network.Stream) { // I'm waiting for incomming connection
		c.messagesReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		handleIncommingMessage(c, receiver)

		c.streamPairs[c.h.ID()] = s.ID()
	})
	return nil
}

func (c *LibP2PClient) Send(msg *common.Message) error {
	var signedMessageBytes []byte
	var err error

	signedMessageBytes, err = json.Marshal(msg)
	if err != nil {
		print(err)
		return err
	}

	c.messagesReadWriter.Write(signedMessageBytes)
	c.messagesReadWriter.WriteByte('\n')
	c.messagesReadWriter.Flush()
	return nil
}

func (c *LibP2PClient) Close() error {
	return nil
}

func NewLibP2PClient(key ed25519.PrivateKey, listenIp string, listenPort string) *LibP2PClient {
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
	client.streamPairs = make(map[peer.ID]string)

	return client
}

func MultiAddrToPeerInfo(peerMultiAddr string) (*peer.AddrInfo, error) {
	maddr, err := multiaddr.NewMultiaddr(peerMultiAddr)
	if err != nil {
		log.Errorln(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Errorln(err)
	}

	return info, err
}

func Connect(client *LibP2PClient, ai peer.AddrInfo) error {
	if client.receiver == nil {
		errMsg := "Missing a client receiver. Listen function should be executed, firstly."
		log.Errorln(errMsg)
		return errors.New(errMsg)
	}

	client.h.Peerstore().AddAddrs(ai.ID, ai.Addrs, peerstore.TempAddrTTL)

	log.Printf("This is a conversation between %s and %s\n", client.h.ID(), ai.ID)

	s, err := client.h.NewStream(context.Background(), ai.ID, p2pStreamName)
	if err != nil {
		log.Errorln(err)
		return err
	}

	client.messagesReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	handleIncommingMessage(client, client.receiver)

	client.streamPairs[ai.ID] = s.ID()

	return nil
}
