package validate

// 用户注册结构体示例
type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required,chinese_name" label:"姓名"`
	Phone    string `json:"phone" binding:"required,phone" label:"手机号"`
	Email    string `json:"email" binding:"required,email" label:"邮箱"`
	Password string `json:"password" binding:"required,strong_password" label:"密码"`
	IDCard   string `json:"id_card" binding:"omitempty,idcard" label:"身份证号"`
	Age      int    `json:"age" binding:"required,min=18,max=100" label:"年龄"`
}

// 用户信息更新结构体示例
type UserUpdateRequest struct {
	ID      uint   `json:"id" binding:"required,min=1" label:"用户ID"`
	Name    string `json:"name" binding:"omitempty,chinese_name" label:"姓名"`
	Phone   string `json:"phone" binding:"omitempty,phone" label:"手机号"`
	Email   string `json:"email" binding:"omitempty,email" label:"邮箱"`
	Age     int    `json:"age" binding:"omitempty,min=18,max=100" label:"年龄"`
	ZipCode string `json:"zip_code" binding:"omitempty,zipcode" label:"邮政编码"`
}

// 查询参数结构体示例
type UserQueryRequest struct {
	Page     int    `form:"page" binding:"required,min=1" label:"页码"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100" label:"每页数量"`
	Keyword  string `form:"keyword" binding:"omitempty,min=1,max=50" label:"关键词"`
}

// 地址信息结构体示例
type AddressRequest struct {
	Province  string `json:"province" binding:"required,chinese" label:"省份"`
	City      string `json:"city" binding:"required,chinese" label:"城市"`
	District  string `json:"district" binding:"required,chinese" label:"区县"`
	Address   string `json:"address" binding:"required" label:"详细地址"`
	ZipCode   string `json:"zip_code" binding:"required,zipcode" label:"邮政编码"`
	IsDefault bool   `json:"is_default" label:"是否默认地址"`
}

// 文章信息结构体示例
type ArticleRequest struct {
	Title       string `json:"title" binding:"required,min=5,max=100" label:"标题"`
	Content     string `json:"content" binding:"required,min=10" label:"内容"`
	PublishDate string `json:"publish_date" binding:"required,date" label:"发布日期"`
	AuthorID    uint   `json:"author_id" binding:"required,min=1" label:"作者ID"`
	CategoryID  uint   `json:"category_id" binding:"required,min=1" label:"分类ID"`
	Tags        string `json:"tags" binding:"omitempty" label:"标签"`
	CoverImage  string `json:"cover_image" binding:"omitempty,url" label:"封面图片"`
}

/*
使用示例：

1. 在启动时初始化验证器（通常在main.go或bootstrap中）：
   import "github.com/zhoudm1743/go-frame/pkg/validate"

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
       return validate.ValidateStruct(user)
   }

4. 自定义验证标签说明：
   - required: 必填字段
   - chinese_name: 中文姓名（2-4个中文字符）
   - chinese: 中文字符
   - phone: 中国手机号
   - email: 邮箱格式
   - strong_password: 强密码（至少8位，包含大小写字母、数字和特殊字符）
   - idcard: 中国身份证号
   - date: 日期格式(YYYY-MM-DD)
   - url: URL格式
   - zipcode: 中国邮政编码(6位数字)
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
