package validate

// 用户注册结构体示例
type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,chinese_name" label:"姓名"`
	Phone    string `json:"phone" validate:"required,phone" label:"手机号"`
	Email    string `json:"email" validate:"required,email" label:"邮箱"`
	Password string `json:"password" validate:"required,strong_password" label:"密码"`
	IDCard   string `json:"id_card" validate:"omitempty,idcard" label:"身份证号"`
	Age      int    `json:"age" validate:"required,min=18,max=100" label:"年龄"`
}

// 用户信息更新结构体示例
type UserUpdateRequest struct {
	ID    uint   `json:"id" validate:"required,min=1" label:"用户ID"`
	Name  string `json:"name" validate:"omitempty,chinese_name" label:"姓名"`
	Phone string `json:"phone" validate:"omitempty,phone" label:"手机号"`
	Email string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Age   int    `json:"age" validate:"omitempty,min=18,max=100" label:"年龄"`
}

// 查询参数结构体示例
type UserQueryRequest struct {
	Page     int    `form:"page" validate:"required,min=1" label:"页码"`
	PageSize int    `form:"page_size" validate:"required,min=1,max=100" label:"每页数量"`
	Keyword  string `form:"keyword" validate:"omitempty,min=1,max=50" label:"关键词"`
}

/*
使用示例：

1. 在启动时初始化验证器（通常在main.go或bootstrap中）：
   import "github.com/zhoudm1743/go-flow/pkg/validate"

   func init() {
       if err := validate.InitValidator(); err != nil {
           log.Fatal("初始化验证器失败:", err)
       }
   }

2. 在控制器中使用验证：
   func RegisterUser(c *gin.Context) {
       var req UserRegisterRequest

       // 使用VerifyJSON进行验证
       if err := util.VerifyUtil.VerifyJSON(c, &req); err != nil {
           c.JSON(400, gin.H{"error": err})
           return
       }

       // 业务逻辑处理...
   }

3. 直接验证结构体：
   func ValidateUserData(user UserRegisterRequest) error {
       return util.VerifyUtil.ValidateStruct(user)
   }

4. 自定义验证标签说明：
   - required: 必填字段
   - chinese_name: 中文姓名（2-4个中文字符）
   - phone: 中国手机号
   - email: 邮箱格式
   - strong_password: 强密码（至少8位，包含大小写字母、数字和特殊字符）
   - idcard: 中国身份证号
   - min/max: 最小值/最大值
   - omitempty: 可选字段（为空时不验证）

5. 错误返回格式：
   验证失败时会返回包含字段名和中文错误信息的map，例如：
   {
       "name": "姓名必须是2-4个中文字符",
       "phone": "手机号必须是有效的中国手机号",
       "password": "密码必须至少8位，包含大小写字母、数字和特殊字符"
   }
*/
