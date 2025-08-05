package client

import (
	"NewCsTeamServer/task"
	"sync"
)

var GlobalClientManager *ClientManager

// ClientManager 管理客户端
type ClientManager struct {
	clients map[uint32]*ClientMetadata
	mutex   sync.RWMutex
}

// NewClientManager 创建新的客户端管理器
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[uint32]*ClientMetadata),
	}
}

// AddClient 添加或更新客户端
func (cm *ClientManager) AddClient(client *ClientMetadata) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if existing, ok := cm.clients[client.ClientID]; ok {
		// 保留现有任务队列，仅更新元数据
		client.TaskQueue = existing.TaskQueue
	} else {
		// 新客户端，初始化任务队列
		client.TaskQueue = &task.TaskQueue{Tasks: []task.Task{}}
	}
	cm.clients[client.ClientID] = client
}

// GetClient 获取客户端
func (cm *ClientManager) GetClient(clientID uint32) (*ClientMetadata, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	client, exists := cm.clients[clientID]
	return client, exists
}
