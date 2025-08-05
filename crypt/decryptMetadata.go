package crypt

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/config"
	"NewCsTeamServer/task"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"strings"
)

// 解密元数据（RSA解密）
func DecryptMetadata(encodedData string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %v", err)
	}

	block, _ := pem.Decode(config.PrivateKey)
	if block == nil {
		return nil, fmt.Errorf("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key error: %v", err)
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("rsa decrypt error: %v", err)
	}

	// 验证MagicHead（0x0000BEEF）
	if len(plaintext) < 8 || binary.BigEndian.Uint32(plaintext[:4]) != 0x0000BEEF {
		return nil, fmt.Errorf("invalid metadata header")
	}

	// 验证长度字段
	dataLen := binary.BigEndian.Uint32(plaintext[4:8])
	if int(dataLen) != len(plaintext)-8 {
		return nil, fmt.Errorf("invalid metadata length")
	}

	return plaintext[8:], nil
}

// 解析元数据
func ParseMetadata(data []byte) (*client.ClientMetadata, error) {
	if len(data) < 51 {
		return nil, fmt.Errorf("metadata too short")
	}

	client := &client.ClientMetadata{
		Key:         data[0:16],
		CharsetANSI: binary.BigEndian.Uint16(data[16:18]),
		CharsetOEM:  binary.BigEndian.Uint16(data[18:20]),
		ClientID:    binary.BigEndian.Uint32(data[20:24]),
		PID:         binary.BigEndian.Uint32(data[24:28]),
		Port:        binary.BigEndian.Uint16(data[28:30]),
		Flag:        data[30],
		OSMajor:     data[31],
		OSMinor:     data[32],
		OSBuild:     binary.BigEndian.Uint16(data[33:35]),
		PtrFunc:     binary.BigEndian.Uint32(data[35:39]),
		PtrGMH:      binary.BigEndian.Uint32(data[39:43]),
		PtrGPA:      binary.BigEndian.Uint32(data[43:47]),
		InternalIP:  binary.LittleEndian.Uint32(data[47:51]),
		TaskQueue:   &task.TaskQueue{Tasks: []task.Task{}},
	}

	// 解析InfoString（Computer\tUser\tProcess）
	info := string(data[51:])
	infoParts := strings.Split(info, "\t")
	if len(infoParts) >= 3 {
		client.ComputerName = infoParts[0]
		client.UserName = infoParts[1]
		client.ProcessName = infoParts[2]
	} else {
		return nil, fmt.Errorf("invalid info string format")
	}

	// 生成AES和HMAC密钥
	hash := sha256.Sum256(client.Key)
	client.AESKey = hash[:16]
	client.HMACKey = hash[16:]

	return client, nil
}
