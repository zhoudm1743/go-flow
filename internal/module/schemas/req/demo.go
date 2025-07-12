package req

// DemoCreateReq 示例创建请求
type DemoCreateReq struct {
	Name        string `json:"name" binding:"required" msg:"请输入示例名称"`
	Description string `json:"description"`
	Status      int8    `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"`
	// 可以根据需要添加更多字段
}

// DemoUpdateReq 示例更新请求
type DemoUpdateReq struct {
	Name        string `json:"name" binding:"required" msg:"请输入示例名称"`
	Description string `json:"description"`
	Status      int8    `json:"status" binding:"oneof=0 1" msg:"状态只能是0或1"`
	// 可以根据需要添加更多字段
}
