package util

const maxBitLength = 29
const maxUidPerLevel = 0x1FFFFFFF

var bitIndexes = [10][29]int{
	{28, 27, 26, 25, 24, 21, 22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 0, 1, 2, 3, 7, 6, 5, 4, 8, 9, 10, 11},
	{27, 26, 25, 24, 21, 22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 1, 2, 3, 7, 6, 5, 4, 8, 9, 10, 11, 0},
	{26, 25, 24, 21, 22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 2, 3, 7, 6, 5, 4, 8, 9, 10, 11, 0, 1},
	{25, 24, 21, 22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 26, 3, 7, 6, 5, 4, 8, 9, 10, 11, 0, 1, 2},
	{24, 21, 22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 26, 25, 7, 6, 5, 4, 8, 9, 10, 11, 0, 1, 2, 3},
	{21, 22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 26, 25, 24, 6, 5, 4, 8, 9, 10, 11, 0, 1, 2, 3, 7},
	{22, 23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 26, 25, 24, 21, 5, 4, 8, 9, 10, 11, 0, 1, 2, 3, 7, 6},
	{23, 20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 26, 25, 24, 21, 22, 4, 8, 9, 10, 11, 0, 1, 2, 3, 7, 6, 5},
	{20, 19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 26, 25, 24, 21, 22, 23, 8, 9, 10, 11, 0, 1, 2, 3, 7, 6, 5, 4},
	{19, 18, 17, 16, 12, 13, 14, 15, 28, 27, 26, 25, 24, 21, 22, 23, 20, 9, 10, 11, 0, 1, 2, 3, 7, 6, 5, 4, 8},
}

func transBits(ii int, n uint64, re bool) uint64 {
	var bits = make([]byte, maxBitLength)
	for i := maxBitLength - 1; n > 0; n, i = n>>1, i-1 {
		bits[i] = byte(n & 1)
	}

	var nits = make([]byte, maxBitLength)
	for ni, bi := range bitIndexes[ii] {
		if re {
			nits[bi] = bits[ni]
		} else {
			nits[ni] = bits[bi]
		}
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
	var ii = NewRand().IntN(10)
	return bn*10000000000 + uint64(ii)*1000000000 + transBits(ii, counter, false)
}

func DecryptUid(uid uint64) uint64 {
	var pre = uid / 1000000000
	var bn = pre / 10
	var ii = pre % 10
	return (bn-1)*maxUidPerLevel + transBits(int(ii), uid%1000000000, true)
}
