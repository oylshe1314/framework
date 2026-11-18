package framework

var (
	appHash    = "EMPTY"
	configHash = "EMPTY"
	dataHash   = "EMPTY"
)

func GetHash() (string, string, string) {
	return appHash, dataHash, configHash
}

func SetAppHash(hash string) {
	appHash = hash[:5]
}

func SetConfigHash(hash string) {
	configHash = hash[:5]
}

func SetDataHash(hash string) {
	dataHash = hash[:5]
}
