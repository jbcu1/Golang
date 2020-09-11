package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

//GeneratePrivateKey if func for generate private key
func GeneratePrivateKey(bitsize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitsize)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	log.Println("Private key generated")
	return privateKey, nil
}

//EncodePrivteKeyToPEM if func for encode private key to bytes array
func EncodePrivteKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	pivBlock := pem.Block{
		Type:    "RSA Private Key",
		Headers: nil,
		Bytes:   privDER,
	}
	privatePEM := pem.EncodeToMemory(&pivBlock)
	return privatePEM
}

//GeneratedPublicKey is func for generate public key of privatekey
func GeneratedPublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	publicBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	log.Println("Public key generated")
	return publicBytes, nil
}

//WriteKeyToFiles is func for write our keys to files
func WriteKeyToFiles(keyBytes []byte, saveFileTo string) error {

	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}
	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
