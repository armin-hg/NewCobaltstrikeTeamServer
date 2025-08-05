package task

// AddTask 添加任务到队列
func (tq *TaskQueue) AddTask(task Task) {
	tq.mutex.Lock()
	defer tq.mutex.Unlock()
	tq.Tasks = append(tq.Tasks, task)
}

// PopTask 获取并移除队列中的第一个任务
func (tq *TaskQueue) PopTask() (Task, bool) {
	tq.mutex.Lock()
	defer tq.mutex.Unlock()
	if len(tq.Tasks) == 0 {
		return Task{}, false
	}
	task := tq.Tasks[0]
	tq.Tasks = tq.Tasks[1:]
	return task, true
}
