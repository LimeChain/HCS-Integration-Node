package libp2p

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

			rawMessage := strings.TrimSpace(string(msg))
			fmt.Println(rawMessage)

			if len(rawMessage) == 0 {
				return
			}

			var signedMessage common.P2PMessage
			err = json.Unmarshal([]byte(rawMessage), &signedMessage)
			fmt.Println(err)

			fmt.Println(string(signedMessage.Msg))

			receiver.Receive(&common.Message{Ctx: context.Background(), Msg: signedMessage.Msg})
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
	signature, err := c.signData(msg.Msg)
	if err != nil {
		return err
	}

	fmt.Println(signature)

	nodePubKey, _ := c.h.Peerstore().PubKey(c.h.ID()).Bytes()

	msg.Signature = signature
	msg.PubKeyData = nodePubKey
	msg.PeerId = c.h.ID()

	peerId := c.h.ID()

	p2pMessage := &common.P2PMessage{Signature: signature, PubKeyData: nodePubKey, PeerId: peerId, Msg: msg.Msg}

	signedMessage := new(bytes.Buffer)
	json.NewEncoder(signedMessage).Encode(p2pMessage)
	signedMessageBytes := signedMessage.Bytes()

	c.messagesReadWriter.Write(signedMessageBytes)
	c.messagesReadWriter.WriteByte('\n')
	c.messagesReadWriter.Flush()
	return nil
}

// sign binary data using the local node's private key
func (c *LibP2PClient) signData(data []byte) ([]byte, error) {
	key := c.h.Peerstore().PrivKey(c.h.ID())
	res, err := key.Sign(data)
	return res, err
}

// func (c *LibP2PClient) SendSignedMessage(msg *common.Message) error {
// 	signature, err := c.signData(msg.Msg)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(signature)

// 	nodePubKey, _ := c.h.Peerstore().PubKey(c.h.ID()).Bytes()

// 	signedMessage := &common.SignedMessage{Signature: signature, PeerId: c.h.ID(), NodePubKey: nodePubKey}
// 	reqBodyBytes := new(bytes.Buffer)
// 	json.NewEncoder(reqBodyBytes).Encode(signedMessage)

// 	signedMessageBytes := reqBodyBytes.Bytes() // this is the []byte

// 	c.messagesReadWriter.Write(signedMessageBytes)
// 	c.messagesReadWriter.WriteByte('\n')
// 	c.messagesReadWriter.Flush()
// 	return nil
// }

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
