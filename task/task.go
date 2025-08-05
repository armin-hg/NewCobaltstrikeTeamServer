package task

import (
	"sync"
	"time"
)

type Task struct {
	ID        string    // 任务ID
	Type      uint32    // 任务类型（4字节，适配客户端）
	Content   []byte    // 任务内容
	CreatedAt time.Time // 任务创建时间
}

// TaskResult 表示任务执行结果
type TaskResult struct {
	TaskID string
	Output []byte
}

// TaskQueue 定义任务队列
type TaskQueue struct {
	Tasks []Task
	mutex sync.Mutex
}
