package hcs

import (
	"context"
	"crypto/ed25519"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/hashgraph/hedera-sdk-go"
	log "github.com/sirupsen/logrus"
)

const SequenceNumberKey = "SequenceNumber"

type HCSClient struct {
	client       *hedera.Client
	mirrorClient *hedera.MirrorClient
	topicID      hedera.ConsensusTopicID
}

func (c *HCSClient) Send(msg *common.Message) error {
	id, err := hedera.NewConsensusMessageSubmitTransaction().
		SetTopicID(c.topicID).
		SetMessage(msg.Msg).
		Execute(c.client)

	if err != nil {
		return err
	}

	_, err = id.GetReceipt(c.client)

	if err != nil {
		return err
	}

	log.Infof("Sent message to HCS with Id :%s\n", id.String())

	return nil
}

func (c *HCSClient) Listen(receiver common.MessageReceiver) error {
	_, err := hedera.NewMirrorConsensusTopicQuery().
		SetTopicID(c.topicID).
		Subscribe(
			*c.mirrorClient,
			func(resp hedera.MirrorConsensusTopicResponse) {
				ctx := context.WithValue(context.Background(), SequenceNumberKey, resp.SequenceNumber)
				receiver.Receive(&common.Message{Msg: resp.Message, Ctx: ctx})
			},
			func(err error) {
				log.Errorln(err.Error())
			})

	if err != nil {
		return err
	}
	return nil
}

func (c *HCSClient) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}

	if err := c.mirrorClient.Close(); err != nil {
		return err
	}
	return nil
}

func NewHCSClient(account string, key ed25519.PrivateKey, mirrorNodeAddress, topicID string) *HCSClient {

	hcsPrvKey, err := hedera.Ed25519PrivateKeyFromBytes(key)
	if err != nil {
		panic(err)
	}

	acc, err := hedera.AccountIDFromString(account)
	if err != nil {
		panic(err)
	}

	mirrorClient, err := hedera.NewMirrorClient(mirrorNodeAddress)
	if err != nil {
		panic(err)
	}

	client := hedera.ClientForTestnet().
		SetOperator(acc, hcsPrvKey)

	hcsTopicId, err := hedera.TopicIDFromString(topicID)
	if err != nil {
		panic(err)
	}

	log.Infof("[HCS] HCS Client started with account ID: %s\n", account)

	return &HCSClient{client: client, mirrorClient: &mirrorClient, topicID: hcsTopicId}
}
