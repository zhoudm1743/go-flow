package types

import (
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"strconv"
	"time"
)

// IDCard 身份证号类型
type IDCard struct {
	BaseSensitive
}

// NewIDCard 创建新的身份证号
func NewIDCard(idCard string) IDCard {
	return IDCard{
		BaseSensitive: NewBaseSensitive(idCard),
	}
}

// Mask 身份证号脱敏 (保留前6位和后4位)
func (id IDCard) Mask() string {
	if id.value == "" {
		return ""
	}
	return MaskString(id.value, 6, 4, '*')
}

// IsValid 验证身份证号格式
func (id IDCard) IsValid() bool {
	if id.value == "" {
		return false
	}

	// 18位身份证正则
	idCardRegex := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
	return idCardRegex.MatchString(id.value)
}

// String 返回脱敏后的字符串
func (id IDCard) String() string {
	return id.Mask()
}

// GetBirthDate 获取出生日期
func (id IDCard) GetBirthDate() (time.Time, error) {
	if len(id.value) < 14 {
		return time.Time{}, nil
	}

	year, err := strconv.Atoi(id.value[6:10])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(id.value[10:12])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(id.value[12:14])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

// GetAge 获取年龄
func (id IDCard) GetAge() int {
	birthDate, err := id.GetBirthDate()
	if err != nil {
		return 0
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// 如果还没到生日，年龄减1
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

// GetGender 获取性别 (true: 男, false: 女)
func (id IDCard) GetGender() (bool, error) {
	if len(id.value) < 17 {
		return false, nil
	}

	genderCode, err := strconv.Atoi(string(id.value[16]))
	if err != nil {
		return false, err
	}

	// 奇数为男性，偶数为女性
	return genderCode%2 == 1, nil
}

// GetGenderString 获取性别字符串
func (id IDCard) GetGenderString() string {
	gender, err := id.GetGender()
	if err != nil {
		return "未知"
	}

	if gender {
		return "男"
	}
	return "女"
}

// GetRegion 获取地区编码（前6位）
func (id IDCard) GetRegion() string {
	if len(id.value) < 6 {
		return ""
	}
	return id.value[:6]
}

// GetRegionName 获取地区名称（简化版本）
func (id IDCard) GetRegionName() string {
	region := id.GetRegion()
	if len(region) < 2 {
		return "未知地区"
	}

	// 这里只是示例，实际应用中可以建立完整的地区编码表
	province := region[:2]
	switch province {
	case "11":
		return "北京市"
	case "12":
		return "天津市"
	case "13":
		return "河北省"
	case "14":
		return "山西省"
	case "15":
		return "内蒙古自治区"
	case "21":
		return "辽宁省"
	case "22":
		return "吉林省"
	case "23":
		return "黑龙江省"
	case "31":
		return "上海市"
	case "32":
		return "江苏省"
	case "33":
		return "浙江省"
	case "34":
		return "安徽省"
	case "35":
		return "福建省"
	case "36":
		return "江西省"
	case "37":
		return "山东省"
	case "41":
		return "河南省"
	case "42":
		return "湖北省"
	case "43":
		return "湖南省"
	case "44":
		return "广东省"
	case "45":
		return "广西壮族自治区"
	case "46":
		return "海南省"
	case "50":
		return "重庆市"
	case "51":
		return "四川省"
	case "52":
		return "贵州省"
	case "53":
		return "云南省"
	case "54":
		return "西藏自治区"
	case "61":
		return "陕西省"
	case "62":
		return "甘肃省"
	case "63":
		return "青海省"
	case "64":
		return "宁夏回族自治区"
	case "65":
		return "新疆维吾尔自治区"
	default:
		return "未知地区"
	}
}

// 实现 GORM Scanner 接口
func (id *IDCard) Scan(value interface{}) error {
	return ScanSensitive(id, value)
}

// 实现 GORM Valuer 接口
func (id IDCard) Value() (driver.Value, error) {
	return ValueSensitive(&id)
}

// 实现 JSON 序列化（返回脱敏数据）
func (id IDCard) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.Mask())
}

// 实现 JSON 反序列化
func (id *IDCard) UnmarshalJSON(data []byte) error {
	var idCard string
	if err := json.Unmarshal(data, &idCard); err != nil {
		return err
	}
	id.SetValue(idCard)
	return nil
}
