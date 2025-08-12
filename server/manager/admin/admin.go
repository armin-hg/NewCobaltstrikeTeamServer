package admin

import (
	"NewCsTeamServer/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"sync"
	"time"
)

type ConnectionManager struct {
	admins   map[string]*Admin
	adminsMu sync.RWMutex
	Logger   *golog.Logger
}

var (
	manager *ConnectionManager
	once    sync.Once
)

func GetConnectionManager() *ConnectionManager {
	once.Do(func() {
		manager = &ConnectionManager{
			admins: make(map[string]*Admin),
			Logger: golog.New(),
		}
	})
	return manager
}

type Admin struct {
	Conn       *websocket.Conn `json:"-"`
	IP         string
	ID         string
	Mutex      sync.Mutex
	LastActive time.Time `json:"last_active"`
	WriteChan  chan Message
	CloseChan  chan struct{}
}
type Message struct {
	ID        string `json:"id,omitempty"`
	AdminID   string `json:"admin_id,omitempty"`
	Type      int    `json:"type"`              //命令类型
	Content   string `json:"content,omitempty"` //命令详细
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (m *ConnectionManager) AddAdmin(conn *websocket.Conn) *Admin {
	m.adminsMu.Lock()
	defer m.adminsMu.Unlock()
	adminid := utils.GetUuid()
	if oldAdmin, exists := m.admins[adminid]; exists {

		if oldAdmin.Conn != nil {
			if err := oldAdmin.Conn.Close(); err != nil { //只允许新的连接

			}
		}
		close(oldAdmin.CloseChan)
		delete(m.admins, adminid)
	}

	admin := &Admin{
		ID:         adminid,
		Conn:       conn,
		LastActive: time.Now(),
		WriteChan:  make(chan Message, 100),
		CloseChan:  make(chan struct{}),
	}
	m.admins[adminid] = admin
	go m.WriteAdminMessages(admin)
	return admin
}

// sendAdminMessage 发送消息给管理员
func (m *ConnectionManager) sendAdminMessage(admin *Admin, msg Message) error {
	admin.Mutex.Lock()
	defer admin.Mutex.Unlock()
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	return admin.Conn.WriteMessage(websocket.TextMessage, data)
}

// WriteAdminMessages 处理管理员消息写入
func (m *ConnectionManager) WriteAdminMessages(admin *Admin) {
	for {
		select {
		case msg := <-admin.WriteChan:
			golog.Println("发送消息:", msg)
			if err := m.sendAdminMessage(admin, msg); err != nil {
				m.Logger.Errorf("Failed to send message to admin %s: %v", admin.ID, err)
				admin.CloseChan <- struct{}{}
				return
			}
		case <-admin.CloseChan:
			return
		}
	}
}

// UpdateAdminHeartbeat 更新管理员心跳
func (m *ConnectionManager) UpdateAdminHeartbeat(admin *Admin) {
	admin.Mutex.Lock()
	admin.LastActive = time.Now()
	admin.Mutex.Unlock()
}

// MonitorAdminHeartbeat 监控管理员心跳
func (m *ConnectionManager) MonitorAdminHeartbeat(admin *Admin, timeout time.Duration) {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				admin.Mutex.Lock()
				inactive := time.Since(admin.LastActive)
				admin.Mutex.Unlock()

				if inactive > timeout {
					m.Logger.Infof("Admin connection timeout: ID=%s, inactive for %v, closing connection", admin.ID, inactive)
					if err := admin.Conn.Close(); err != nil {
						m.Logger.Errorf("Failed to close admin connection %s: %v", admin.ID, err)
					}
					m.RemoveAdmin(admin, "")
					return
				}
			}
		}
	}()
}

// RemoveAdmin removes an admin connection safely
func (m *ConnectionManager) RemoveAdmin(admin *Admin, clientIP string) {
	m.adminsMu.Lock()
	defer m.adminsMu.Unlock()

	// Check if admin exists in the map
	if storedAdmin, ok := m.admins[admin.ID]; ok && storedAdmin == admin {
		// Delete from map first to prevent re-entrant calls
		delete(m.admins, admin.ID)

		// Close channels only if they are not nil and not closed
		if admin.CloseChan != nil {
			select {
			case <-admin.CloseChan:
				// Channel is already closed
				m.Logger.Warnf("CloseChan already closed for admin %s", admin.ID)
			default:
				close(admin.CloseChan)
			}
		}
		if admin.WriteChan != nil {
			select {
			case <-admin.WriteChan:
				// Channel is already closed
				m.Logger.Warnf("WriteChan already closed for admin %s", admin.ID)
			default:
				close(admin.WriteChan)
			}
		}

		// Close connection if it exists and is not already closed
		if admin.Conn != nil {
			if err := admin.Conn.Close(); err != nil {
				// Ignore "use of closed network connection" error
				if err.Error() != "use of closed network connection" {
					m.Logger.Errorf("Failed to close admin connection %s: %v", admin.ID, err)
				}
			}
		}

		m.Logger.Infof("Admin %s disconnected", clientIP)
	} else {
		m.Logger.Warnf("Admin %s not found or mismatched for IP %s", admin.ID, clientIP)
	}
}

// SendMessageToAdmin 向指定管理员发送消息
func (m *ConnectionManager) SendMessageToAdmin(adminID string, msg Message) {
	m.adminsMu.RLock()
	admin, ok := m.admins[adminID]
	m.adminsMu.RUnlock()
	if !ok {
		m.Logger.Errorf("Admin %s not found", adminID)
		return
	}
	select {
	case admin.WriteChan <- msg:
	case <-admin.CloseChan:
		m.Logger.Errorf("Admin %s disconnected", adminID)
	}
}

// BroadcastToAdmins 向所有管理员广播消息
func (m *ConnectionManager) BroadcastToAdmins(msg Message) {
	m.adminsMu.RLock()
	defer m.adminsMu.RUnlock()
	msg.Timestamp = time.Now().Unix()
	for adminID, admin := range m.admins {
		select {
		case admin.WriteChan <- msg:
		default:
			m.Logger.Warnf("Admin %s message channel full, dropping message: %v", adminID, msg)
		}
	}
}

// BroadcastToAdminsNoMe 向除指定管理员外的其他管理员广播消息
func (m *ConnectionManager) BroadcastToAdminsNoMe(msg Message, excludeAdminID string) {
	m.adminsMu.RLock()
	defer m.adminsMu.RUnlock()
	for adminID, admin := range m.admins {
		if adminID == excludeAdminID {
			continue
		}
		select {
		case admin.WriteChan <- msg:
		default:
			m.Logger.Warnf("Admin %s message channel full, dropping message: %v", adminID, msg)
		}
	}
}
