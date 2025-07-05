package req

// AccountCreateReq 创建请求
type AccountCreateReq struct {
	Username  string `json:"username" validate:"required" label:"用户名"`
	Nickname  string `json:"nickname" validate:"required" label:"昵称"`
	Email     string `json:"email" validate:"required,email" label:"邮箱"`
	Password  string `json:"password" validate:"required" label:"密码"`
	Status    int    `json:"status" validate:"omitempty" label:"状态"`
	LastLogin string `json:"last_login" validate:"omitempty" label:"最后登录时间"`
}

// AccountUpdateReq 更新请求
type AccountUpdateReq struct {
	Username  string `json:"username" validate:"omitempty" label:"用户名"`
	Nickname  string `json:"nickname" validate:"omitempty" label:"昵称"`
	Email     string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Password  string `json:"password" validate:"omitempty" label:"密码"`
	Status    int    `json:"status" validate:"omitempty" label:"状态"`
	LastLogin string `json:"last_login" validate:"omitempty" label:"最后登录时间"`
}

// AccountQueryReq 查询请求
type AccountQueryReq struct {
	PageReq
	Keyword string `form:"keyword" validate:"omitempty" json:"keyword"` // 关键字
	Status  int    `form:"status" validate:"omitempty" json:"status"`   // 状态
	Email   string `form:"email" validate:"omitempty" json:"email"`     // 邮箱
}
