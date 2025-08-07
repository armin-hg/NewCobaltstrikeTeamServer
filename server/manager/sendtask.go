package manager

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/server/public"
	"NewCsTeamServer/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// HTTP处理任务下发请求
func handleIssueTask(c *gin.Context) {
	var req struct {
		ClientID uint32 `json:"client_id"`
		Type     uint32 `json:"type"` // 4字节命令类型，适配客户端
		Content  string `json:"content"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, public.ApiResponse{
			Code: 1,
			Msg:  "err",
		})

		return
	}
	client, exists := client.GlobalClientManager.GetClient(req.ClientID)
	if !exists {
		c.JSON(http.StatusBadRequest, public.ApiResponse{
			Code: 1,
			Msg:  "err",
		})

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

	log.Printf("Task issued to client %d: ID=%s, Type=%d, Content=%s", req.ClientID, task.ID, task.Type, task.Content)
	c.String(http.StatusOK, "Task issued successfully")
	c.JSON(http.StatusOK, public.ApiResponse{Code: 0, Msg: "任务下发成功"})
	return
}
