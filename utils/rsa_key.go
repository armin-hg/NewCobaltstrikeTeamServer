package utils

import (
	"NewCsTeamServer/config"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/jkeys089/jserial"
	"log"
	"os"
)

// KeyPair 结构体用于存储提取的公钥和私钥
type KeyPair struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func GetRsaKey(path string) error { //TODO 读取后，存放进配置信息里，用于加解密
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("无法读取文件: %v", err)
		return err
	}

	// 解析序列化对象
	objects, err := jserial.ParseSerializedObject(f)
	if err != nil {
		log.Fatalf("解析序列化对象失败: %v", err)
		return err
	}

	// 将解析结果转换为 JSON 格式以便查看结构（调试用）
	//pretty, err := json.MarshalIndent(objects, "", "  ")
	//if err != nil {
	//	log.Fatalf("JSON 序列化失败: %v", err)
	//}

	// 提取 RSA 公钥和私钥
	keyPair, err := extractRSAKeys(objects)
	if err != nil {
		log.Fatalf("提取 RSA 密钥失败: %v", err)
		return err
	}

	// 将公钥和私钥编码为 PEM 格式并打印
	if keyPair.PublicKey != nil {
		publicKeyPEM := encodePublicKeyToPEM(keyPair.PublicKey)
		config.PublicKey = publicKeyPEM //更新公钥
		log.Println("RSA 公钥 (PEM 格式):")
		log.Println(string(publicKeyPEM))
	}
	if keyPair.PrivateKey != nil {
		privateKeyPEM := encodePrivateKeyToPEM(keyPair.PrivateKey)
		log.Println("RSA 私钥 (PEM 格式):")
		config.PrivateKey = privateKeyPEM //更新私钥
		log.Println(string(privateKeyPEM))
	}
	return nil
	// 可选：将密钥保存到文件
	//if keyPair.PublicKey != nil {
	//	err = os.WriteFile("public_key.pem", encodePublicKeyToPEM(keyPair.PublicKey), 0644)
	//	if err != nil {
	//		log.Fatalf("保存公钥到文件失败: %v", err)
	//		return err
	//	}
	//	log.Println("公钥已保存至 public_key.pem")
	//}
	//if keyPair.PrivateKey != nil {
	//	err = os.WriteFile("private_key.pem", encodePrivateKeyToPEM(keyPair.PrivateKey), 0600)
	//	if err != nil {
	//		log.Fatalf("保存私钥到文件失败: %v", err)
	//		return err
	//	}
	//	log.Println("私钥已保存至 private_key.pem")
	//}
}

// extractRSAKeys 从解析的对象中提取 RSA 公钥和私钥
func extractRSAKeys(objects interface{}) (*KeyPair, error) {
	keyPair := &KeyPair{}

	// 打印 objects 类型以调试
	log.Printf("objects 类型: %T", objects)

	// 假设 objects 是一个数组，提取第一个元素
	arr, ok := objects.([]interface{})
	if !ok || len(arr) == 0 {
		return nil, logError("对象格式不正确，期望非空数组，实际类型: %T", objects)
	}
	data := arr[0]
	log.Printf("arr[0] 类型: %T", data)

	// 导航到 java.security.KeyPair
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, logError("data 字段格式不正确，期望 map 类型，实际类型: %T", data)
	}
	log.Printf("dataMap 内容: %+v", dataMap)

	array, ok := dataMap["array"].(map[string]interface{})
	if !ok {
		return nil, logError("array 字段格式不正确，实际类型: %T", dataMap["array"])
	}
	log.Printf("array 内容: %+v", array)

	// 直接从 array["value"] 获取 java.security.KeyPair
	value, ok := array["value"].(map[string]interface{})
	if !ok {
		return nil, logError("value 字段格式不正确，实际类型: %T", array["value"])
	}
	log.Printf("value 内容: %+v", value)

	// 获取 extends 字段
	extends, ok := value["extends"].(map[string]interface{})
	if !ok {
		return nil, logError("extends 字段格式不正确，实际类型: %T", value["extends"])
	}
	log.Printf("extends 内容: %+v", extends)

	// 获取 java.security.KeyPair
	keyPairData, ok := extends["java.security.KeyPair"].(map[string]interface{})
	if !ok {
		return nil, logError("java.security.KeyPair 字段格式不正确，实际类型: %T", extends["java.security.KeyPair"])
	}
	log.Printf("keyPairData 内容: %+v", keyPairData)

	// 提取公钥
	if pubKeyData, exists := keyPairData["publicKey"].(map[string]interface{}); exists {
		log.Printf("publicKey 内容: %+v", pubKeyData)
		pubKeyEncoded, ok := pubKeyData["encoded"].([]interface{})
		if !ok {
			return nil, logError("公钥 encoded 字段格式不正确，实际类型: %T", pubKeyData["encoded"])
		}
		pubKeyBytes := make([]byte, len(pubKeyEncoded))
		for i, v := range pubKeyEncoded {
			var num int64
			switch val := v.(type) {
			case int8:
				num = int64(val)
			case int:
				num = int64(val)
			case float64:
				num = int64(val)
			default:
				return nil, logError("公钥 encoded 数组元素格式不正确，索引 %d，实际类型: %T", i, v)
			}
			pubKeyBytes[i] = byte(num)
		}
		pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
		if err != nil {
			return nil, logError("解析公钥失败: %v", err)
		}
		rsaPubKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, logError("公钥不是 RSA 格式")
		}
		keyPair.PublicKey = rsaPubKey
	}

	// 提取私钥
	if privKeyData, exists := keyPairData["privateKey"].(map[string]interface{}); exists {
		log.Printf("privateKey 内容: %+v", privKeyData)
		privKeyEncoded, ok := privKeyData["encoded"].([]interface{})
		if !ok {
			return nil, logError("私钥 encoded 字段格式不正确，实际类型: %T", privKeyData["encoded"])
		}
		privKeyBytes := make([]byte, len(privKeyEncoded))
		for i, v := range privKeyEncoded {
			var num int64
			switch val := v.(type) {
			case int8:
				num = int64(val)
			case int:
				num = int64(val)
			case float64:
				num = int64(val)
			default:
				return nil, logError("私钥 encoded 数组元素格式不正确，索引 %d，实际类型: %T", i, v)
			}
			privKeyBytes[i] = byte(num)
		}
		privKey, err := x509.ParsePKCS8PrivateKey(privKeyBytes)
		if err != nil {
			return nil, logError("解析私钥失败: %v", err)
		}
		rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
		if !ok {
			return nil, logError("私钥不是 RSA 格式")
		}
		keyPair.PrivateKey = rsaPrivKey
	}

	return keyPair, nil
}

// encodePublicKeyToPEM 将 RSA 公钥编码为 PEM 格式
func encodePublicKeyToPEM(pubKey *rsa.PublicKey) []byte {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatalf("公钥编码失败: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
}

// encodePrivateKeyToPEM 将 RSA 私钥编码为 PEM 格式
func encodePrivateKeyToPEM(privKey *rsa.PrivateKey) []byte {
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		log.Fatalf("私钥编码失败: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})
}

// logError 辅助函数，用于记录错误并返回
func logError(format string, args ...interface{}) error {
	log.Printf(format, args...)
	return fmt.Errorf(format, args...)
}
