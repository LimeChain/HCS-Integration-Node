package hcs

import (
	"crypto/ed25519"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/blockchain"
	"github.com/hashgraph/hedera-sdk-go"
)

type HCSClient struct {
	client  *hedera.Client
	topicID hedera.ConsensusTopicID
}

func (c *HCSClient) Send(msg *blockchain.BlockchainMessage) error {
	return nil // TODO implement
}

func (c *HCSClient) Listen(receiver blockchain.BlockchainMessageReceiver) error {
	return nil // TODO implement
}

func NewHCSClient(account string, key ed25519.PrivateKey, topicID string) *HCSClient {

	hcsPrvKey, err := hedera.Ed25519PrivateKeyFromBytes(key)
	if err != nil {
		panic(err)
	}
	acc, err := hedera.AccountIDFromString(account)
	if err != nil {
		panic(err)
	}

	client := hedera.ClientForTestnet().
		SetOperator(acc, hcsPrvKey)

	hcsTopicId, err := hedera.TopicIDFromString(topicID)
	if err != nil {
		panic(err)
	}

	return &HCSClient{client: client, topicID: hcsTopicId}
}
