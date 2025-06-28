package types

import (
	"time"
)

// User 用户模型示例
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Phone     Phone     `gorm:"type:text" json:"phone"`     // 手机号（加密存储）
	IDCard    IDCard    `gorm:"type:text" json:"id_card"`   // 身份证号（加密存储）
	BankCard  BankCard  `gorm:"type:text" json:"bank_card"` // 银行卡号（加密存储）
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserInfo 用户信息（包含解析后的敏感数据信息）
type UserInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`          // 脱敏后的手机号
	Provider string `json:"phone_provider"` // 运营商
	IDCard   string `json:"id_card"`        // 脱敏后的身份证号
	Age      int    `json:"age"`            // 从身份证解析的年龄
	Gender   string `json:"gender"`         // 从身份证解析的性别
	Region   string `json:"region"`         // 从身份证解析的地区
	BankCard string `json:"bank_card"`      // 脱敏后的银行卡号
	BankName string `json:"bank_name"`      // 银行名称
	CardType string `json:"card_type"`      // 卡片类型
	Email    string `json:"email"`
}

// ToUserInfo 转换为用户信息（脱敏版本）
func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		ID:       u.ID,
		Name:     u.Name,
		Phone:    u.Phone.Mask(),
		Provider: u.Phone.GetProvider(),
		IDCard:   u.IDCard.Mask(),
		Age:      u.IDCard.GetAge(),
		Gender:   u.IDCard.GetGenderString(),
		Region:   u.IDCard.GetRegionName(),
		BankCard: u.BankCard.FormatCard(),
		BankName: u.BankCard.GetBankName(),
		CardType: u.BankCard.GetCardType(),
		Email:    u.Email,
	}
}

// ValidateData 验证所有敏感数据
func (u *User) ValidateData() map[string]string {
	errors := make(map[string]string)

	if !u.Phone.IsValid() {
		errors["phone"] = "手机号格式不正确"
	}

	if !u.IDCard.IsValid() {
		errors["id_card"] = "身份证号格式不正确"
	}

	if !u.BankCard.IsValid() {
		errors["bank_card"] = "银行卡号格式不正确"
	}

	return errors
}

/*
使用示例：

1. 在启动时初始化敏感数据类型模块：
   import "github.com/zhoudm1743/go-flow/types"

   func init() {
       // 使用默认密钥
       types.Init()

       // 或使用自定义密钥
       // types.InitWithKey("your-32-char-encryption-key-here")
   }

2. 在模型中使用敏感数据类型：
   type User struct {
       ID       uint           `gorm:"primaryKey"`
       Name     string         `json:"name"`
       Phone    types.Phone    `gorm:"type:text" json:"phone"`
       IDCard   types.IDCard   `gorm:"type:text" json:"id_card"`
       BankCard types.BankCard `gorm:"type:text" json:"bank_card"`
   }

3. 创建用户：
   user := User{
       Name:     "张三",
       Phone:    types.NewPhone("13812345678"),
       IDCard:   types.NewIDCard("110101199001011234"),
       BankCard: types.NewBankCard("6225881234567890"),
   }

   // 验证数据
   if errors := user.ValidateData(); len(errors) > 0 {
       // 处理验证错误
       fmt.Println("验证失败:", errors)
       return
   }

   // 保存到数据库（数据会自动加密存储）
   db.Create(&user)

4. 查询用户：
   var user User
   db.First(&user, 1)

   // 获取脱敏信息
   userInfo := user.ToUserInfo()
   fmt.Printf("用户信息: %+v\n", userInfo)

   // 输出示例：
   // 手机号: 138****5678 (中国移动)
   // 身份证: 110101****1234 (男, 24岁, 北京市)
   // 银行卡: 6225 **** **** 7890 (招商银行, 银联)

5. 直接操作敏感数据：
   phone := types.NewPhone("13812345678")
   fmt.Println("原始手机号:", phone.Raw())        // 13812345678
   fmt.Println("脱敏手机号:", phone.Mask())       // 138****5678
   fmt.Println("运营商:", phone.GetProvider())    // 中国移动

   idCard := types.NewIDCard("110101199001011234")
   fmt.Println("脱敏身份证:", idCard.Mask())       // 110101****1234
   fmt.Println("年龄:", idCard.GetAge())          // 34
   fmt.Println("性别:", idCard.GetGenderString()) // 男

   bankCard := types.NewBankCard("6225881234567890")
   fmt.Println("格式化卡号:", bankCard.FormatCard()) // 6225 **** **** 7890
   fmt.Println("银行:", bankCard.GetBankName())      // 招商银行

6. JSON序列化（自动脱敏）：
   userJson, _ := json.Marshal(user)
   // 输出的JSON中敏感字段已经自动脱敏

7. 数据库字段类型建议：
   - 敏感字段使用 TEXT 类型存储加密数据
   - 如需要索引，可以单独存储hash值或者使用部分明文（如手机号后4位）

8. 安全注意事项：
   - 生产环境必须使用环境变量 CRYPTO_KEY 设置32位密钥
   - 定期轮换加密密钥
   - 在日志中避免输出原始敏感数据
   - 数据库备份时注意加密密钥的管理
*/
