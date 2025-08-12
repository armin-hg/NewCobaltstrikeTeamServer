package utils

import (
	"encoding/binary"
	"github.com/google/uuid"
	"net"
	"strconv"
)

// Uint32ToIPString 将小端序 uint32 转换为 IP 地址字符串
func Uint32ToIPString(ipUint32 uint32) string {
	// 创建 4 字节缓冲区
	buf := make([]byte, 4)
	// 小端序写入 uint32
	binary.LittleEndian.PutUint32(buf, ipUint32)
	// 转换为 net.IP 并格式化为字符串
	return net.IP(buf).String()
}

// StringToUint32 将字符串转换为 uint32
func StringToUint32(s string) uint32 {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

func GetUuid() string {
	return uuid.New().String()
}
