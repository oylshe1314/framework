package util

import "encoding/binary"

func BytesToUint16(p []byte) uint16 {
	return binary.LittleEndian.Uint16(p)
}

func BytesToUint32(p []byte) uint32 {
	return binary.LittleEndian.Uint32(p)
}

func BytesToUint64(p []byte) uint64 {
	return binary.LittleEndian.Uint64(p)
}

func Uint16ToBytes(p []byte, v uint16) []byte {
	return binary.LittleEndian.AppendUint16(p, v)
}

func Uint32ToBytes(p []byte, v uint32) []byte {
	return binary.LittleEndian.AppendUint32(p, v)
}

func Uint64ToBytes(p []byte, v uint64) []byte {
	return binary.LittleEndian.AppendUint64(p, v)
}

func PutUint16ToBytes(p []byte, v uint16) {
	binary.LittleEndian.PutUint16(p, v)
}

func PutUint32ToBytes(p []byte, v uint32) {
	binary.LittleEndian.PutUint32(p, v)
}

func PutUint64ToBytes(p []byte, v uint64) {
	binary.LittleEndian.PutUint64(p, v)
}

func Compose2uint16(h, l uint16) uint32 {
	return uint32(h)<<16 | uint32(l)
}

func Compose2uint32(h, l uint32) uint64 {
	return uint64(h)<<32 | uint64(l)
}

func Split2uint16(value uint32) (uint16, uint16) {
	return uint16(value >> 16 & 0xFFFF), uint16(value & 0xFFFF)
}

func Split2uint32(value uint64) (uint32, uint32) {
	return uint32(value >> 32 & 0xFFFFFFFF), uint32(value & 0xFFFFFFFF)
}
