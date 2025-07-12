package req

// IdReq ID请求
type IdReq struct {
	ID uint `uri:"id" form:"id" binding:"required,min=1"`
}

// PageReq 分页请求
type PageReq struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=100"`
}

// GetOffset 获取偏移量
func (r *PageReq) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// GetLimit 获取限制
func (r *PageReq) GetLimit() int {
	return r.PageSize
}
