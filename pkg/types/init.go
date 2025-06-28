package types

// Init 初始化types模块
func Init() {
	InitCrypto("")
}

// InitWithKey 使用自定义密钥初始化
func InitWithKey(cryptoKey string) {
	InitCrypto(cryptoKey)
}
