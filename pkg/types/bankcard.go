package types

import (
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"strconv"
)

// BankCard 银行卡号类型
type BankCard struct {
	BaseSensitive
}

// NewBankCard 创建新的银行卡号
func NewBankCard(cardNumber string) BankCard {
	return BankCard{
		BaseSensitive: NewBaseSensitive(cardNumber),
	}
}

// Mask 银行卡号脱敏 (保留前4位和后4位)
func (bc BankCard) Mask() string {
	if bc.value == "" {
		return ""
	}
	return MaskString(bc.value, 4, 4, '*')
}

// IsValid 验证银行卡号格式 (使用Luhn算法)
func (bc BankCard) IsValid() bool {
	if bc.value == "" {
		return false
	}

	// 检查是否只包含数字
	cardRegex := regexp.MustCompile(`^\d+$`)
	if !cardRegex.MatchString(bc.value) {
		return false
	}

	// 检查长度（一般银行卡号为13-19位）
	if len(bc.value) < 13 || len(bc.value) > 19 {
		return false
	}

	// 使用Luhn算法验证
	return bc.luhnCheck()
}

// luhnCheck Luhn算法校验
func (bc BankCard) luhnCheck() bool {
	var sum int
	nDigits := len(bc.value)
	parity := nDigits % 2

	for i, digit := range bc.value {
		d, err := strconv.Atoi(string(digit))
		if err != nil {
			return false
		}

		if i%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}

	return sum%10 == 0
}

// String 返回脱敏后的字符串
func (bc BankCard) String() string {
	return bc.Mask()
}

// GetCardType 获取卡片类型
func (bc BankCard) GetCardType() string {
	if len(bc.value) < 1 {
		return "未知卡型"
	}

	// 根据卡号长度和首位数字判断卡片类型
	firstDigit := string(bc.value[0])
	length := len(bc.value)

	switch {
	case firstDigit == "4" && (length == 13 || length == 16 || length == 19):
		return "Visa"
	case (firstDigit == "5" || bc.value[:2] >= "22" && bc.value[:2] <= "27") && length == 16:
		return "MasterCard"
	case (bc.value[:2] == "34" || bc.value[:2] == "37") && length == 15:
		return "American Express"
	case bc.value[:2] == "62" && (length >= 16 && length <= 19):
		return "银联"
	default:
		return "其他"
	}
}

// GetBankName 获取银行名称（根据BIN码判断）
func (bc BankCard) GetBankName() string {
	if len(bc.value) < 6 {
		return "未知银行"
	}

	// BIN码（Bank Identification Number）前6位
	bin := bc.value[:6]

	// 这里只是示例，实际应用中需要完整的BIN码数据库
	switch {
	case bin >= "620000" && bin <= "629999":
		return bc.getUnionPayBank(bin)
	case bin[:4] == "4000":
		return "招商银行"
	case bin[:4] == "4367":
		return "建设银行"
	case bin[:4] == "4390":
		return "工商银行"
	case bin[:4] == "5187":
		return "民生银行"
	case bin[:4] == "5221":
		return "中国银行"
	case bin[:4] == "5451":
		return "广发银行"
	case bin[:4] == "5571":
		return "平安银行"
	default:
		return "其他银行"
	}
}

// getUnionPayBank 获取银联卡对应的银行
func (bc BankCard) getUnionPayBank(bin string) string {
	// 简化的银联BIN码对应表
	switch {
	case bin >= "622700" && bin <= "622799":
		return "工商银行"
	case bin >= "621700" && bin <= "621799":
		return "建设银行"
	case bin >= "621000" && bin <= "621099":
		return "邮储银行"
	case bin >= "622800" && bin <= "622899":
		return "农业银行"
	case bin >= "621600" && bin <= "621699":
		return "中国银行"
	case bin >= "622300" && bin <= "622399":
		return "交通银行"
	case bin >= "622600" && bin <= "622699":
		return "招商银行"
	case bin >= "622900" && bin <= "622999":
		return "兴业银行"
	case bin >= "623000" && bin <= "623099":
		return "光大银行"
	case bin >= "623100" && bin <= "623199":
		return "华夏银行"
	case bin >= "623200" && bin <= "623299":
		return "民生银行"
	case bin >= "623300" && bin <= "623399":
		return "广发银行"
	case bin >= "623400" && bin <= "623499":
		return "平安银行"
	case bin >= "623500" && bin <= "623599":
		return "浦发银行"
	default:
		return "银联卡"
	}
}

// FormatCard 格式化显示银行卡号（每4位一组）
func (bc BankCard) FormatCard() string {
	masked := bc.Mask()
	if len(masked) == 0 {
		return ""
	}

	var formatted string
	for i, char := range masked {
		if i > 0 && i%4 == 0 {
			formatted += " "
		}
		formatted += string(char)
	}
	return formatted
}

// GetBIN 获取BIN码（前6位）
func (bc BankCard) GetBIN() string {
	if len(bc.value) < 6 {
		return ""
	}
	return bc.value[:6]
}

// 实现 GORM Scanner 接口
func (bc *BankCard) Scan(value interface{}) error {
	return ScanSensitive(bc, value)
}

// 实现 GORM Valuer 接口
func (bc BankCard) Value() (driver.Value, error) {
	return ValueSensitive(&bc)
}

// 实现 JSON 序列化（返回脱敏数据）
func (bc BankCard) MarshalJSON() ([]byte, error) {
	return json.Marshal(bc.Mask())
}

// 实现 JSON 反序列化
func (bc *BankCard) UnmarshalJSON(data []byte) error {
	var cardNumber string
	if err := json.Unmarshal(data, &cardNumber); err != nil {
		return err
	}
	bc.SetValue(cardNumber)
	return nil
}
