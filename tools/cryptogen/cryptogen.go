package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"os"
)

func generateEd25519Keypair() ed25519.PrivateKey {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	return priv
}

func createTestnetHederaAccount(key ed25519.PrivateKey) hedera.AccountID {
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("HCS_OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("HCS_OPERATOR_PRV_KEY"))
	if err != nil {
		panic(err)
	}

	newKey, err := hedera.Ed25519PrivateKeyFromBytes(key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hedera Account Public Key: %v\n", newKey.PublicKey().String())

	client := hedera.ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	transactionID, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	newAccountID := transactionReceipt.GetAccountID()

	// TODO Fund account

	transactionID, err = hedera.NewCryptoTransferTransaction().
		AddSender(operatorAccountID, hedera.NewHbar(100)).
		AddRecipient(newAccountID, hedera.NewHbar(100)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	_, err = transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	return newAccountID
}

func main() {

	key := generateEd25519Keypair()

	hcsAccId := createTestnetHederaAccount(key)

	fmt.Printf("Hedera Account Id %v\n", hcsAccId.String())

	x509Encoded, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		panic(err)
	}

	pemEncodedPrv := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	fmt.Println("Private key in PEM format")
	fmt.Println(string(pemEncodedPrv))

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/config/key.pem", path), pemEncodedPrv, 0644)
	if err != nil {
		panic(err)
	}

}
