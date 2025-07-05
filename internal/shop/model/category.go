package model

import (
	"github.com/zhoudm1743/go-flow/pkg/types"
)

// Category 分类模型
type Category struct {
	types.GormModel
	Name        string `json:"name" gorm:"size:50;not null;comment:分类名称"`
	Description string `json:"description" gorm:"size:200;comment:分类描述"`
	ParentID    uint   `json:"parent_id" gorm:"default:0;comment:父分类ID"`
	Level       int    `json:"level" gorm:"default:1;comment:分类层级"`
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`
	Status      int    `json:"status" gorm:"default:1;comment:状态 1-启用 0-禁用"`
}
