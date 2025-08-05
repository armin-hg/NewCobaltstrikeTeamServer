package server

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/crypt"
	"NewCsTeamServer/task"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// HandleBeaconResult 处理客户端的任务结果
func HandleBeaconResult(w http.ResponseWriter, r *http.Request, cm *client.ClientManager) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数中的 ClientID
	query := r.URL.Query()
	clientIDStr := query.Get("id")
	if clientIDStr == "" {
		http.Error(w, "缺少 client_id 参数", http.StatusBadRequest)
		return
	}
	clientID, err := strconv.ParseUint(clientIDStr, 10, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("无效的 client_id: %v", err), http.StatusBadRequest)
		return
	}

	// 检查 ClientID 是否存在
	client, exists := cm.GetClient(uint32(clientID))
	if !exists {
		log.Printf("处理结果失败: 未找到客户端 %d", clientID)
		http.Error(w, fmt.Sprintf("未找到客户端 %d", clientID), http.StatusNotFound)
		return
	}

	// 读取 POST 请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("读取请求体失败: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Printf("handleBeaconResult: 客户端 %d, 原始请求体=%x", clientID, body)

	// 验证 SendLength
	if len(body) < 4+crypt.HmacHashLen {
		http.Error(w, fmt.Sprintf("无效请求体: 太短 (%d 字节)", len(body)), http.StatusBadRequest)
		return
	}
	//sendLen := binary.BigEndian.Uint32(body[:4])
	//if int(sendLen) != len(body)-4 {
	//	log.Printf("无效 SendLength: 期望 %d, 实际 %d", sendLen, len(body)-4)
	//	http.Error(w, "无效 SendLength", http.StatusBadRequest)
	//	return
	//}

	// 分离加密数据和 HMAC
	encryptedData := body[4 : len(body)-crypt.HmacHashLen]
	receivedHMAC := body[len(body)-crypt.HmacHashLen:]
	log.Printf("handleBeaconResult: 客户端 %d, 加密数据=%x, 收到 HMAC=%x", clientID, encryptedData, receivedHMAC)

	// 验证 HMAC
	expectedHMAC := crypt.HmacHash(encryptedData, client.HMACKey)
	if !bytes.Equal(receivedHMAC, expectedHMAC) {
		log.Printf("HMAC 验证失败，客户端 %d", clientID)
		http.Error(w, "无效 HMAC 签名", http.StatusUnauthorized)
		return
	}

	// 解密数据（补充 IV 前缀）
	iv := []byte("abcdefghijklmnop")
	decryptedData, err := crypt.AesCBCDecrypt(append(iv, encryptedData...), client.AESKey)
	if err != nil {
		log.Printf("解密失败，客户端 %d: %v", clientID, err)
		http.Error(w, fmt.Sprintf("解密失败: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("handleBeaconResult: 客户端 %d, 解密数据=%x", clientID, decryptedData)

	// 解析任务结果
	result, err := task.ParseTaskResult(decryptedData)
	if err != nil {
		log.Printf("解析任务结果失败，客户端 %d: %v", clientID, err)
		http.Error(w, fmt.Sprintf("无效任务结果: %v", err), http.StatusBadRequest)
		return
	}

	// 存储结果
	//cm.AddResult(uint32(clientID), result)
	log.Printf("收到客户端 %d 的结果: TaskID=%s, Output=%s", clientID, result.TaskID, string(result.Output))

	// 返回成功响应
	// 返回成功响应
	w.Write([]byte(""))
}
