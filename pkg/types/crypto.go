package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// CryptoConfig 加密配置
type CryptoConfig struct {
	Key string // 32位密钥
}

var defaultCrypto *CryptoConfig

// InitCrypto 初始化加密配置
func InitCrypto(key string) {
	if key == "" {
		// 从环境变量获取密钥，如果没有则使用默认密钥（生产环境请使用环境变量）
		key = os.Getenv("CRYPTO_KEY")
		if key == "" {
			key = "go-flow-default-32-char-key-123" // 32位默认密钥，生产环境必须更换
		}
	}

	// 确保密钥长度为32位
	if len(key) > 32 {
		key = key[:32]
	} else if len(key) < 32 {
		for len(key) < 32 {
			key += "0"
		}
	}

	defaultCrypto = &CryptoConfig{Key: key}
}

// Encrypt 加密数据
func (c *CryptoConfig) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return "", err
	}

	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密数据
func (c *CryptoConfig) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("密文太短")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GetDefaultCrypto 获取默认加密实例
func GetDefaultCrypto() *CryptoConfig {
	if defaultCrypto == nil {
		InitCrypto("")
	}
	return defaultCrypto
}
