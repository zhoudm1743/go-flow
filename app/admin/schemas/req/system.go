package req

type SystemAdminListReq struct {
	Username string `form:"username" json:"username"`
	Nickname string `form:"nickname" json:"nickname"`
	Status   int8   `form:"status,default=-1" json:"status"`
	Role     string `form:"role" json:"role"`
}

type SystemAdminAddReq struct {
	Username string   `form:"username" json:"username" validate:"required,min=3,max=30"`
	Nickname string   `form:"nickname" json:"nickname" validate:"required,min=3,max=30"`
	Password string   `form:"password" json:"password" validate:"required,min=6,max=30"`
	Role     []string `form:"role" json:"role" validate:"required"`
	Email    string   `form:"email" json:"email" validate:"required,email"`
	Phone    string   `form:"phone" json:"phone" validate:"required,phone"`
	Avatar   string   `form:"avatar" json:"avatar"`
}

type SystemAdminEditReq struct {
	ID       uint     `form:"id" json:"id" validate:"required"`
	Username string   `form:"username" json:"username" validate:"required,min=3,max=30"`
	Nickname string   `form:"nickname" json:"nickname" validate:"required,min=3,max=30"`
	Role     []string `form:"role" json:"role" validate:"required"`
	Email    string   `form:"email" json:"email" validate:"required,email"`
	Phone    string   `form:"phone" json:"phone" validate:"required,phone"`
	Avatar   string   `form:"avatar" json:"avatar"`
}
