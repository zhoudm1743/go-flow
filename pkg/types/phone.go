package types

import (
	"database/sql/driver"
	"encoding/json"
	"regexp"
)

// Phone 手机号类型
type Phone struct {
	BaseSensitive
}

// NewPhone 创建新的手机号
func NewPhone(phone string) Phone {
	return Phone{
		BaseSensitive: NewBaseSensitive(phone),
	}
}

// Mask 手机号脱敏 (保留前3位和后4位)
func (p Phone) Mask() string {
	if p.value == "" {
		return ""
	}
	return MaskString(p.value, 3, 4, '*')
}

// IsValid 验证手机号格式
func (p Phone) IsValid() bool {
	if p.value == "" {
		return false
	}

	// 中国手机号正则
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(p.value)
}

// String 返回脱敏后的字符串
func (p Phone) String() string {
	return p.Mask()
}

// GetArea 获取手机号归属地区（前3位）
func (p Phone) GetArea() string {
	if len(p.value) < 3 {
		return ""
	}
	return p.value[:3]
}

// GetProvider 获取运营商信息
func (p Phone) GetProvider() string {
	if len(p.value) < 3 {
		return "未知"
	}

	prefix := p.value[:3]
	switch prefix {
	case "130", "131", "132", "145", "155", "156", "166", "171", "175", "176", "185", "186", "196":
		return "中国联通"
	case "133", "134", "135", "136", "137", "138", "139", "147", "150", "151", "152", "157", "158", "159", "172", "178", "182", "183", "184", "187", "188", "195", "198":
		return "中国移动"
	case "140", "141", "142", "143", "144", "146", "149", "153", "173", "174", "177", "180", "181", "189", "190", "191", "193", "199":
		return "中国电信"
	default:
		return "未知运营商"
	}
}

// 实现 GORM Scanner 接口
func (p *Phone) Scan(value interface{}) error {
	return ScanSensitive(p, value)
}

// 实现 GORM Valuer 接口
func (p Phone) Value() (driver.Value, error) {
	return ValueSensitive(&p)
}

// 实现 JSON 序列化（返回脱敏数据）
func (p Phone) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Mask())
}

// 实现 JSON 反序列化
func (p *Phone) UnmarshalJSON(data []byte) error {
	var phone string
	if err := json.Unmarshal(data, &phone); err != nil {
		return err
	}
	p.SetValue(phone)
	return nil
}
