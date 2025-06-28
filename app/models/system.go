package models

import "github.com/zhoudm1743/go-flow/pkg/types"

type SystemAdmin struct {
	BaseModel
	Username    string      `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Password    string      `gorm:"size:255;not null" json:"-"`
	Status      int8        `gorm:"default:1" json:"status"`
	Role        []string    `gorm:"type:json;serializer:json" json:"role"`
	Email       string      `gorm:"uniqueIndex;size:200;not null" json:"email"`
	Phone       types.Phone `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Avatar      string      `gorm:"size:255;not null" json:"avatar"`
	LastLoginAt int64       `json:"last_login_at"`
	LastLoginIP string      `gorm:"size:45;not null" json:"last_login_ip"`
}

type SystemMenu struct {
	BaseModel
	Pid  uint           `gorm:"default:0" json:"pid"`
	Name string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Meta SystemMenuMeta `gorm:"type:json;serializer:json" json:"meta"`
	Sort int            `gorm:"default:0" json:"sort"`
}

type SystemMenuMeta struct {
	Title      string   `json:"title"`
	SvgIcon    string   `json:"svgIcon"`
	ElIcon     string   `json:"elIcon"`
	Hidden     bool     `json:"hidden"`
	Roles      []string `json:"roles"`
	Breadcrumb bool     `json:"breadcrumb"`
	Affix      bool     `json:"affix"`
	AlwaysShow bool     `json:"alwaysShow"`
	ActiveMenu string   `json:"activeMenu"`
	KeepAlive  bool     `json:"keepAlive"`
}
