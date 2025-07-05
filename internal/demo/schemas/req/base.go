package req

// PageReq 分页请求参数
type PageReq struct {
	PageNo   int `form:"pageNo,default=1" validate:"omitempty,gte=1"`          // 页码
	PageSize int `form:"pageSize,default=20" validate:"omitempty,gt=0,lte=60"` // 每页大小
}

type IdReq struct {
	ID uint `form:"id" validate:"required" json:"id"` // 主键ID
}

type IdListReq struct {
	Ids []string `form:"ids" validate:"required,dive" json:"ids"` // 主键ID列表
}

type KeyReq struct {
	Key string `form:"key" validate:"required" json:"key"` // 关键字
}
