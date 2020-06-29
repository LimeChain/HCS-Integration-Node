package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func getPeerPublicKey() ed25519.PublicKey {
	peerPubKeyHex := os.Getenv("PEER_PUB_KEY")
	if len(peerPubKeyHex) == 0 {
		log.Panic("No peer pub key supplied")
	}

	peerPubKeyString, err := hex.DecodeString(peerPubKeyHex)
	if err != nil {
		panic(err)
	}

	return ed25519.PublicKey(peerPubKeyString)
}

func getPrivateKey() ed25519.PrivateKey {

	keyPath := os.Getenv("PRIVATE_KEY_FILE")

	info, err := os.Stat(keyPath)
	if os.IsNotExist(err) {
		panic(fmt.Errorf("The key does not exist"))
	}
	if info.IsDir() {
		panic(fmt.Errorf("The key is actually directory"))
	}
	pemKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}

	pemBlock, _ := pem.Decode(pemKey)
	if pemBlock == nil {
		panic(fmt.Errorf("Not a PEM encoded file"))
	}

	prvKeyBytes, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}

	priv, ok := prvKeyBytes.(ed25519.PrivateKey)

	if !ok {
		panic("The inputted key is not ED25519 private key")
	}

	return priv

}
