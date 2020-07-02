package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func createTestnetTopic(a, b hedera.PublicKey) hedera.ConsensusTopicID {
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("HCS_OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("HCS_OPERATOR_PRV_KEY"))
	if err != nil {
		panic(err)
	}

	pubKeys := []hedera.PublicKey{a,b}

	client := hedera.ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	thresholdKey := hedera.NewThresholdKey(1).AddAll(pubKeys)

	transactionID, err := hedera.NewConsensusTopicCreateTransaction().
		SetAdminKey(operatorPrivateKey.PublicKey()).
		SetAutoRenewAccountID(operatorAccountID).
		SetSubmitKey(thresholdKey).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	topicID := transactionReceipt.GetConsensusTopicID()

	return topicID

}

func main() {

	a := os.Getenv("A_PUB_KEY")
	b := os.Getenv("B_PUB_KEY")

	aPubKey, err := hedera.Ed25519PublicKeyFromString(a)
	if err != nil {
		panic(err)
	}

	bPubKey, err := hedera.Ed25519PublicKeyFromString(b)
	if err != nil {
		panic(err)
	}

	topicID := createTestnetTopic(aPubKey, bPubKey)

	fmt.Printf("topicID: %v\n", topicID)

}
