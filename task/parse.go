package task

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

// ParseTaskResult 解析任务结果数据
func ParseTaskResult(data []byte) (TaskResult, error) {
	var result TaskResult
	buf := bytes.NewBuffer(data)

	// 检查缓冲区长度
	if buf.Len() < 12 { // Counter(4) + ResultLength(4) + ReplyType(4)
		return result, fmt.Errorf("缓冲区太短: %d 字节", buf.Len())
	}

	// 读取 Counter（4字节）
	counterBytes := make([]byte, 4)
	_, err := buf.Read(counterBytes)
	if err != nil {
		return result, fmt.Errorf("读取 Counter 失败: %v", err)
	}
	counter := binary.BigEndian.Uint32(counterBytes)

	// 读取 ResultLength（4字节）
	resultLenBytes := make([]byte, 4)
	_, err = buf.Read(resultLenBytes)
	if err != nil {
		return result, fmt.Errorf("读取 ResultLength 失败: %v", err)
	}
	resultLen := binary.BigEndian.Uint32(resultLenBytes)

	// 读取 ReplyType（4字节，作为 TaskID）
	replyTypeBytes := make([]byte, 4)
	_, err = buf.Read(replyTypeBytes)
	if err != nil {
		return result, fmt.Errorf("读取 ReplyType 失败: %v", err)
	}
	replyType := binary.BigEndian.Uint32(replyTypeBytes)
	result.TaskID = fmt.Sprintf("%d", replyType)

	// 计算有效 ResultData 长度
	expectedDataLen := int(resultLen) - 4 // ResultLength = ReplyType(4) + ResultData
	if expectedDataLen < 0 {
		return result, fmt.Errorf("无效 ResultLength: %d", resultLen)
	}

	// 读取 ResultData
	remainingLen := buf.Len()
	if remainingLen < expectedDataLen {
		return result, fmt.Errorf("ResultData 长度不足: 期望 %d, 实际 %d", expectedDataLen, remainingLen)
	}

	// 读取有效 Output
	result.Output = make([]byte, expectedDataLen)
	_, err = buf.Read(result.Output)
	if err != nil {
		return result, fmt.Errorf("读取 Output 失败: %v", err)
	}

	// 检查是否有填充数据
	if buf.Len() > 0 {
		log.Printf("警告: 检测到填充数据，长度=%d, 填充=%x", buf.Len(), buf.Bytes())
	}

	log.Printf("parseTaskResult: Counter=%d, ReplyType=%d, TaskID=%s, OutputLength=%d, Output=%s, FullData=%x",
		counter, replyType, result.TaskID, len(result.Output), string(result.Output), data)
	return result, nil
}
