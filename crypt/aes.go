package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

const HmacHashLen = 16

func AesCBCDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建 AES 加密器失败: %v", err)
	}
	if len(data)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("无效数据长度: %d", len(data))
	}
	decrypted := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, data[:16]) // 使用前 16 字节作为 IV
	mode.CryptBlocks(decrypted, data)
	return decrypted[16:], nil // 移除 IV
}

func HmacHash(data []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key) // 假设 config.HmacKey 可用
	h.Write(data)
	return h.Sum(nil)[:HmacHashLen]
}

// AES解密数据
func AesDecrypt(encryptedData, aesKey, iv, hmacKey, signature []byte) ([]byte, error) {
	hmacCalc := hmac.New(sha256.New, hmacKey)
	hmacCalc.Write(encryptedData)
	if !hmac.Equal(hmacCalc.Sum(nil)[:16], signature) {
		return nil, fmt.Errorf("hmac verification failed")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher error: %v", err)
	}

	if len(encryptedData)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("encrypted data length invalid")
	}

	decrypted := make([]byte, len(encryptedData))
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(decrypted, encryptedData)

	decrypted = pkcs7UnPadding(decrypted)
	return decrypted, nil
}

// AES加密数据
func AesEncrypt(data, aesKey, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher error: %v", err)
	}

	data = pkcs7Padding(data, aes.BlockSize)
	encrypted := make([]byte, len(data))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(encrypted, data)
	return encrypted, nil
}

// 添加HMAC签名
func AddHMAC(data, hmacKey []byte) []byte {
	hmacCalc := hmac.New(sha256.New, hmacKey)
	hmacCalc.Write(data)
	return hmacCalc.Sum(nil)[:16]
}
