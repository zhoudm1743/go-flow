package model

import (
	"github.com/zhoudm1743/go-frame/pkg/types"
)

// Demo 示例模型
type Demo struct {
	types.GormModel
	Name        string `json:"name" gorm:"size:100;not null;comment:示例名称"`
	Description string `json:"description" gorm:"size:500;comment:示例描述"`
	Status      int8    `json:"status" gorm:"default:1;comment:状态 1-启用 0-禁用"`
	// 可以根据需要添加更多字段
}
