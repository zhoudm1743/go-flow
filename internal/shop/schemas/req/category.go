package req

// CategoryCreateReq 分类创建请求
type CategoryCreateReq struct {
	Name        string `json:"name" binding:"required" msg:"请输入分类名称"`
	Description string `json:"description"`
	ParentID    uint   `json:"parent_id"`
	Level       int    `json:"level"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"`
}

// CategoryUpdateReq 分类更新请求
type CategoryUpdateReq struct {
	Name        string `json:"name" binding:"required" msg:"请输入分类名称"`
	Description string `json:"description"`
	ParentID    uint   `json:"parent_id"`
	Level       int    `json:"level"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"`
}
