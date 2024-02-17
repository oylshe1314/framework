package util

const maxBitLength = 29
const maxUidPerLevel = 0x1FFFFFFF

var bitIndexes = []int{
	28, 27, 26, 25, 24, 21, 22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 0, 1, 2, 3, 7, 6, 5, 4, 8, 9, 10, 11,
}

func transBits(n uint64) uint64 {
	var bits = make([]byte, maxBitLength)
	for i := maxBitLength - 1; n > 0; n, i = n>>1, i-1 {
		bits[i] = byte(n & 1)
	}

	var nits = make([]byte, maxBitLength)
	for ni, bi := range bitIndexes {
		nits[ni] = bits[bi]
	}

	n = 0
	for _, nit := range nits {
		n = n<<1 | uint64(nit)
	}

	return n
}

func EncryptUid(counter uint64) uint64 {
	var bn = counter/maxUidPerLevel + 1
	counter %= maxUidPerLevel
	return bn*1000000000 + transBits(counter)
}
