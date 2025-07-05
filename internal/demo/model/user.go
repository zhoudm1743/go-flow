package model

import (
	"github.com/zhoudm1743/go-flow/pkg/types"
)

// User 用户模型
type User struct {
	types.GormModel
	Name  string `json:"name" gorm:"size:100;not null"`
	Email string `json:"email" gorm:"size:100;uniqueIndex;not null"`
}
