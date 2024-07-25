package utils

import "encoding/binary"

func Uint64ToHex(num uint64) []byte {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, num)
	return buffer
}
