package model

import (
	"github.com/zhoudm1743/go-flow/pkg/types"
)

// Account 模型
type Account struct {
	types.GormModel
	Username  string `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Nickname  string `json:"nickname" gorm:"size:100"`
	Email     string `json:"email" gorm:"size:100;uniqueIndex"`
	Password  string `json:"-" gorm:"size:100;not null"` // 密码不输出到JSON
	Status    int    `json:"status" gorm:"default:1"`    // 状态：1正常，0禁用
	LastLogin string `json:"last_login" gorm:"size:30"`  // 最后登录时间
}
