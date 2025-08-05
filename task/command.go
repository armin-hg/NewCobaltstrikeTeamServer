package task

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 构造 Shell 命令数据（CommandType=78），适配官方CommandBuilder
func buildShellCommand(content []byte) []byte {
	var command bytes.Buffer

	// 添加 %COMSPEC%（长度+字符串）
	comspec := []byte("%COMSPEC%")
	commandLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandLenBytes, uint32(len(comspec)))
	command.Write(commandLenBytes)
	command.Write(comspec)

	// 添加命令内容（长度+字符串，/C 前缀）
	commandContent := []byte(" /C " + string(content))
	commandLenBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(commandLenBytes, uint32(len(commandContent)))
	command.Write(commandLenBytes)
	command.Write(commandContent)

	// 添加短整数 0（2字节）
	command.Write([]byte{0, 0})

	result := command.Bytes()
	log.Printf("buildShellCommand: PathLen=%d, Path=%s, CmdLen=%d, Cmd=%s, Result=%x",
		len(comspec), comspec, len(commandContent), commandContent, result)
	return result
}
