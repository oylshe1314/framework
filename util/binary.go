package util

func Compose2Uint16(h, l uint16) uint32 {
	return uint32(h)<<16 | uint32(l)
}

func Compose2Uint32(h, l uint32) uint64 {
	return uint64(h)<<32 | uint64(l)
}

func Split2Uint16(n uint32) (uint16, uint16) {
	return uint16(n >> 16 & 0xFFFF), uint16(n & 0xFFFF)
}

func Split2Uint32(n uint64) (uint32, uint32) {
	return uint32(n >> 32 & 0xFFFFFFFF), uint32(n & 0xFFFFFFFF)
}
