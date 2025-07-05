package req

// ProductCreateReq 产品创建请求
type ProductCreateReq struct {
	Name        string  `json:"name" binding:"required" msg:"请输入产品名称"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0" msg:"请输入正确的产品价格"`
	Stock       int     `json:"stock" binding:"gte=0" msg:"库存不能为负数"`
	CategoryID  uint    `json:"category_id"`
	Status      int     `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"`
}

// ProductUpdateReq 产品更新请求
type ProductUpdateReq struct {
	Name        string  `json:"name" binding:"required" msg:"请输入产品名称"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0" msg:"请输入正确的产品价格"`
	Stock       int     `json:"stock" binding:"gte=0" msg:"库存不能为负数"`
	CategoryID  uint    `json:"category_id"`
	Status      int     `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"`
}
