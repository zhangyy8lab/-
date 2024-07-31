package utils

import (
	"encoding/binary"
	"os"
)

func Uint64ToHex(num uint64) []byte {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, num)

	return buffer
}

// FileIsExist 判断db是否存在
func FileIsExist(dbPath string) bool {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}
	return true
}
