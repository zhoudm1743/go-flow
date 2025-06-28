package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// SensitiveData 敏感数据接口
type SensitiveData interface {
	// Mask 脱敏显示
	Mask() string
	// Raw 获取原始数据
	Raw() string
	// IsValid 验证数据有效性
	IsValid() bool
	// Encrypt 加密存储
	Encrypt() (string, error)
	// Decrypt 解密数据
	Decrypt(encryptedData string) error
}

// BaseSensitive 基础敏感数据结构
type BaseSensitive struct {
	value     string
	encrypted bool
}

// NewBaseSensitive 创建基础敏感数据
func NewBaseSensitive(value string) BaseSensitive {
	return BaseSensitive{
		value:     value,
		encrypted: false,
	}
}

// Raw 获取原始数据
func (bs *BaseSensitive) Raw() string {
	return bs.value
}

// SetValue 设置值
func (bs *BaseSensitive) SetValue(value string) {
	bs.value = value
	bs.encrypted = false
}

// IsValid 基础验证（子类可以重写）
func (bs *BaseSensitive) IsValid() bool {
	return bs.value != ""
}

// Mask 基础脱敏（子类可以重写）
func (bs *BaseSensitive) Mask() string {
	return MaskString(bs.value, 1, 1, '*')
}

// Encrypt 加密存储
func (bs *BaseSensitive) Encrypt() (string, error) {
	if bs.value == "" {
		return "", nil
	}

	crypto := GetDefaultCrypto()
	encrypted, err := crypto.Encrypt(bs.value)
	if err != nil {
		return "", err
	}

	bs.encrypted = true
	return encrypted, nil
}

// Decrypt 解密数据
func (bs *BaseSensitive) Decrypt(encryptedData string) error {
	if encryptedData == "" {
		bs.value = ""
		bs.encrypted = false
		return nil
	}

	crypto := GetDefaultCrypto()
	decrypted, err := crypto.Decrypt(encryptedData)
	if err != nil {
		return err
	}

	bs.value = decrypted
	bs.encrypted = false
	return nil
}

// MaskString 通用脱敏方法
func MaskString(value string, keepStart, keepEnd int, maskChar rune) string {
	if value == "" {
		return ""
	}

	runes := []rune(value)
	length := len(runes)

	if length <= keepStart+keepEnd {
		// 如果字符串太短，只显示部分字符
		if length <= 2 {
			return string(maskChar)
		}
		return string(runes[0]) + strings.Repeat(string(maskChar), length-2) + string(runes[length-1])
	}

	result := make([]rune, length)

	// 保留开头
	for i := 0; i < keepStart && i < length; i++ {
		result[i] = runes[i]
	}

	// 中间用掩码字符
	for i := keepStart; i < length-keepEnd; i++ {
		result[i] = maskChar
	}

	// 保留结尾
	for i := length - keepEnd; i < length; i++ {
		result[i] = runes[i]
	}

	return string(result)
}

// 实现GORM的Scanner和Valuer接口的通用方法

// ScanSensitive 通用扫描方法
func ScanSensitive(dest SensitiveData, value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("无法扫描 %T 到敏感数据类型", value)
	}

	// 尝试解密数据
	return dest.Decrypt(str)
}

// ValueSensitive 通用值方法
func ValueSensitive(src SensitiveData) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}

	// 加密存储
	encrypted, err := src.Encrypt()
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}
