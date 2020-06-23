package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

func main() {

	key := generateEd25519Keypair()

	x509Encoded, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		panic(err)
	}

	pemEncodedPrv := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

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
