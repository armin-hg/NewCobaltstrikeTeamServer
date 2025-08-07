package client

import (
	"NewCsTeamServer/task"
	"sync"
	"time"
)

// ClientMetadata 存储客户端元数据
type ClientMetadata struct {
	Key          []byte // 16字节随机密钥
	CharsetANSI  uint16 // ANSI代码页
	CharsetOEM   uint16 // OEM代码页
	ClientID     uint32 // 客户端ID
	PID          uint32 // 进程ID
	Port         uint16 // SSH端口（未实现，默认为0）
	Flag         uint8  // 标志位（1: 无, 2: x64 Agent, 4: x64 System, 8: Admin）
	OSMajor      uint8  // 操作系统主版本
	OSMinor      uint8  // 操作系统次版本
	OSBuild      uint16 // 操作系统构建号
	PtrFunc      uint32 // Smart Inject 函数地址（未实现，默认为0）
	PtrGMH       uint32 // GetModuleHandle地址（未实现，默认为0）
	PtrGPA       uint32 // GetProcAddress地址（未实现，默认为0）
	InternalIP   uint32 // 本地IP（小端序）
	IpAddress    string //公网ip地址
	ComputerName string // 计算机名
	UserName     string // 用户名
	LastActive   time.Time
	ProcessName  string          // 进程名
	AESKey       []byte          // AES密钥（16字节）
	HMACKey      []byte          // HMAC密钥（16字节）
	TaskQueue    *task.TaskQueue // 任务队列
	Mutex        sync.Mutex      `json:"-"`
}

type HostList struct {
	ClientID   uint32 // 客户端ID
	PID        uint32 // 进程ID
	Flag       uint8
	InternalIP string //内网ip
	IpAddress  string //公网ip地址

	ComputerName string
	UserName     string // 用户名

	ProcessName string    // 进程名
	LastActive  time.Time `json:"last_active"`
}
