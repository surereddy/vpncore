package misc

import (
	"encoding/hex"
	"strings"
	"strconv"
)



func BytesEqual(x, y []byte) bool {
	// copy from net.ip.go:bytesEqual()

	if len(x) != len(y) {
		return false
	}
	for i, b := range x {
		if y[i] != b {
			return false
		}
	}
	return true
}

func ByteToHexString(value byte) string {
	return hex.EncodeToString([]byte{value})
}

func BytesToUint16(value []byte) uint16 {
	return uint16(value[0]) << 8 | uint16(value[1])
}

func BytesToUint32(value []byte) uint32 {
	return uint32(value[0]) << 24 |
		uint32(value[1]) << 16 |
		uint32(value[2]) << 8 |
		uint32(value[3])
}

func BytesToInt64(value []byte) int64 {
	return int64(value[0]) << 56 |
		int64(value[1]) << 48 |
		int64(value[2]) << 40 |
		int64(value[3]) << 32 |
		int64(value[4]) << 24 |
		int64(value[5]) << 16 |
		int64(value[6]) << 8 |
		int64(value[7])
}

func BytesToHexString(value []byte) string {
	strs := make([]string, len(value))
	for i, b := range value {
		strs[i] = hex.EncodeToString([]byte{b})
	}
	return "[" + strings.Join(strs, ",") + "]"
}

func Uint16ToBytes(value uint16, b []byte) []byte {
	return append(b, byte(value >> 8), byte(value))
}

func Uint16ToString(value uint16) string {
	return strconv.Itoa(int(value))
}

func Uint32ToBytes(value uint32, b []byte) []byte {
	return append(b, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value))
}

func Uint32ToString(value uint32) string {
	return strconv.FormatUint(uint64(value), 10)
}

func IntToBytes(value int, b []byte) []byte {
	return append(b, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value))
}

func IntToString(value int) string {
	return Int64ToString(int64(value))
}

func Int64ToBytes(value int64, b []byte) []byte {
	return append(b,
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value))
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}


