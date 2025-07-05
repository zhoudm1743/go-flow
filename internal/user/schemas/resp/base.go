package resp

// PageResp 分页响应
type PageResp struct {
	Total    int64       `json:"total"`    // 总记录数
	PageNo   int         `json:"pageNo"`   // 当前页码
	PageSize int         `json:"pageSize"` // 每页条数
	Pages    int         `json:"pages"`    // 总页数
	HasNext  bool        `json:"hasNext"`  // 是否有下一页
	Data     interface{} `json:"data"`     // 列表数据
}

// NewPageResp 创建分页响应
func NewPageResp(pageNo, pageSize int, total int64, data interface{}) *PageResp {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	hasNext := pageNo < pages

	return &PageResp{
		Total:    total,
		PageNo:   pageNo,
		PageSize: pageSize,
		Pages:    pages,
		HasNext:  hasNext,
		Data:     data,
	}
}
