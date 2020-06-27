package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

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
