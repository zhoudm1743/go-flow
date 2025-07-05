package model

import (
	"github.com/zhoudm1743/go-flow/pkg/types"
)

// Product 产品模型
type Product struct {
	types.GormModel
	Name        string  `json:"name" gorm:"size:100;not null;comment:产品名称"`
	Description string  `json:"description" gorm:"size:500;comment:产品描述"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null;default:0.00;comment:产品价格"`
	Stock       int     `json:"stock" gorm:"not null;default:0;comment:库存"`
	CategoryID  uint    `json:"category_id" gorm:"comment:分类ID"`
	Status      int     `json:"status" gorm:"default:1;comment:状态 1-上架 0-下架"`
}
