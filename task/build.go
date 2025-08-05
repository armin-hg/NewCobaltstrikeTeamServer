package task

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

// BuildCommandData 构造命令数据，适配官方Shell方法
func BuildCommandData(task Task) ([]byte, error) {
	var command bytes.Buffer

	// 添加 CommandType（4字节）
	commandTypeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandTypeBytes, task.Type)
	command.Write(commandTypeBytes)

	// 构造 CommandData
	var commandData []byte
	switch task.Type {
	case 78:
		commandData = buildShellCommand(task.Content)
	case 2:
		return nil, fmt.Errorf("session type (CommandType=2) not implemented")
	default:
		return nil, fmt.Errorf("unsupported command type: %d", task.Type)
	}

	// 添加 CommandDataLength（4字节）
	commandLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandLenBytes, uint32(len(commandData)))
	command.Write(commandLenBytes)
	command.Write(commandData)

	log.Printf("buildCommandData: Task ID=%s, Type=%d, CommandDataLength=%d, CommandData=%x", task.ID, task.Type, len(commandData), command.Bytes())
	return command.Bytes(), nil
}

// BuildTaskPacket 构造任务数据包，适配客户端DecryptPacket和ParsePacket
func BuildTaskPacket(tasks []Task) ([]byte, error) {
	var packet bytes.Buffer

	// 添加时间戳（4字节）
	timestamp := make([]byte, 4)
	binary.BigEndian.PutUint32(timestamp, uint32(time.Now().Unix()))
	packet.Write(timestamp)

	// 预留长度字段（4字节）
	lengthBytes := make([]byte, 4)
	packet.Write(lengthBytes)

	// 添加命令
	totalCommandLen := 0
	for _, task := range tasks {
		commandData, err := BuildCommandData(task)
		if err != nil {
			return nil, err
		}
		packet.Write(commandData)
		totalCommandLen += len(commandData)
		log.Printf("buildTaskPacket: Task ID=%s, Type=%d, CommandData=%x", task.ID, task.Type, commandData)
	}

	// 更新长度字段（仅CommandType + CommandData）
	totalLen := uint32(totalCommandLen)
	binary.BigEndian.PutUint32(packet.Bytes()[4:8], totalLen)
	log.Printf("buildTaskPacket: Total length=%d, Packet=%x", totalLen, packet.Bytes())

	return packet.Bytes(), nil
}
