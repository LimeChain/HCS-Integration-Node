package hcs

import (
	"context"
	"crypto/ed25519"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/blockchain"
	"github.com/hashgraph/hedera-sdk-go"
	log "github.com/sirupsen/logrus"
)

type HCSClient struct {
	client       *hedera.Client
	mirrorClient *hedera.MirrorClient
	topicID      hedera.ConsensusTopicID
}

/* Example how to use Send - will be removed when implemented
rcpt, err := hcsClient.Send(&blockchain.BlockchainMessage{Msg: []byte(fmt.Sprintf("Hello HCS from Go! Message %v", 1)), Ctx: context.TODO()})
// if err != nil {
// 	panic(err)
// }
// log.Println(rcpt.Status)
*/

func (c *HCSClient) Send(msg *blockchain.BlockchainMessage) (*hedera.TransactionReceipt, error) {
	id, err := hedera.NewConsensusMessageSubmitTransaction().
		SetTopicID(c.topicID).
		SetMessage(msg.Msg).
		Execute(c.client)

	if err != nil {
		return nil, err
	}

	receipt, err := id.GetReceipt(c.client)

	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (c *HCSClient) Listen(receiver blockchain.BlockchainMessageReceiver) error {
	_, err := hedera.NewMirrorConsensusTopicQuery().
		SetTopicID(c.topicID).
		Subscribe(
			*c.mirrorClient,
			func(resp hedera.MirrorConsensusTopicResponse) {
				receiver.Receive(&blockchain.BlockchainMessage{Msg: resp.Message, Ctx: context.TODO()})
			},
			func(err error) {
				log.Errorln(err.Error())
			})

	if err != nil {
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

	return &HCSClient{client: client, mirrorClient: &mirrorClient, topicID: hcsTopicId}
}
