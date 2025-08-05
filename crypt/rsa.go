package crypt

import (
	"NewCsTeamServer/config"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 去除填充
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	padding := int(data[length-1])
	return data[:length-padding]
}
func RsaEncrypt(origData []byte) ([]byte, error) {

	block, _ := pem.Decode(config.PrivateKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(config.PublicKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv := privInterface.(*rsa.PrivateKey)
	return rsa.DecryptPKCS1v15(rand.Reader, priv, origData)
}
