package models

import "github.com/zhoudm1743/go-flow/pkg/types"

// 租户表
type SystemTenant struct {
	BaseModel
	Name     string `gorm:"not null;default:'';comment:租户名称" json:"name"`
	Code     string `gorm:"not null;default:'';comment:租户编码" json:"code"`
	Type     string `gorm:"not null;default:'';comment:租户类型" json:"type"`
	Status   int8   `gorm:"default:1" json:"status"`
	ExpireAt int64  `gorm:"default:0" json:"expire_at"`
	// 数据库配置
	DatabaseHost     string `gorm:"not null;default:'';comment:数据库主机" json:"database_host"`
	DatabasePort     string `gorm:"not null;default:'3306';comment:数据库端口" json:"database_port"`
	DatabaseName     string `gorm:"not null;default:'';comment:数据库名称" json:"database_name"`
	DatabaseUser     string `gorm:"not null;default:'';comment:数据库用户" json:"database_user"`
	DatabasePassword string `gorm:"not null;default:'';comment:数据库密码" json:"-"`
}

// 租户角色表
type SystemTenantRole struct {
	BaseModel
	TenantID    uint   `gorm:"not null;comment:租户ID" json:"tenant_id"`
	Name        string `gorm:"not null;default:'';comment:角色名称" json:"name"`
	Code        string `gorm:"not null;default:'';comment:角色编码" json:"code"`
	Description string `gorm:"not null;default:'';comment:角色描述" json:"description"`
	Status      int8   `gorm:"default:1;comment:状态:1正常,0禁用" json:"status"`
}

// 租户权限表
type SystemTenantPermission struct {
	BaseModel
	TenantID    uint   `gorm:"not null;comment:租户ID" json:"tenant_id"`
	Name        string `gorm:"not null;default:'';comment:权限名称" json:"name"`
	Code        string `gorm:"not null;default:'';comment:权限编码" json:"code"`
	Resource    string `gorm:"not null;default:'';comment:资源路径" json:"resource"`
	Action      string `gorm:"not null;default:'';comment:操作动作" json:"action"`
	Description string `gorm:"not null;default:'';comment:权限描述" json:"description"`
	Status      int8   `gorm:"default:1;comment:状态:1正常,0禁用" json:"status"`
}

// 租户角色权限关联表
type SystemTenantRolePermission struct {
	BaseModel
	TenantID     uint `gorm:"not null;comment:租户ID" json:"tenant_id"`
	RoleID       uint `gorm:"not null;comment:角色ID" json:"role_id"`
	PermissionID uint `gorm:"not null;comment:权限ID" json:"permission_id"`
}

// 租户用户角色关联表
type SystemTenantUserRole struct {
	BaseModel
	TenantID uint `gorm:"not null;comment:租户ID" json:"tenant_id"`
	UserID   uint `gorm:"not null;comment:用户ID" json:"user_id"`
	RoleID   uint `gorm:"not null;comment:角色ID" json:"role_id"`
}

// Casbin策略规则表
type SystemCasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PType string `gorm:"size:100;not null" json:"ptype"`
	V0    string `gorm:"size:100;not null" json:"v0"`
	V1    string `gorm:"size:100;not null" json:"v1"`
	V2    string `gorm:"size:100;not null" json:"v2"`
	V3    string `gorm:"size:100" json:"v3"`
	V4    string `gorm:"size:100" json:"v4"`
	V5    string `gorm:"size:100" json:"v5"`
}

// TableName 指定表名
func (SystemCasbinRule) TableName() string {
	return "casbin_rule"
}

type SystemAdmin struct {
	BaseModel
	TenantID    uint        `gorm:"not null;default:0" json:"tenant_id"`
	Username    string      `gorm:"size:100;not null" json:"username"`
	Password    string      `gorm:"size:255;not null" json:"-"`
	Status      int8        `gorm:"default:1" json:"status"`
	Roles       []string    `gorm:"type:json;serializer:json" json:"roles"`
	Email       string      `gorm:"uniqueIndex;size:200;not null;default:''" json:"email"`
	Phone       types.Phone `gorm:"uniqueIndex;size:20;not null;default:''" json:"phone"`
	Avatar      string      `gorm:"size:255;not null;default:''" json:"avatar"`
	LastLoginAt int64       `json:"last_login_at" gorm:"default:0"`
	LastLoginIP string      `gorm:"size:45;not null;default:''" json:"last_login_ip"`
}

type SystemMenu struct {
	BaseModel
	CreateUser uint   `gorm:"not null" json:"create_user"`
	UpdateUser uint   `gorm:"not null" json:"update_user"`
	Name       string `gorm:"not null;default:'';comment:菜单名称" json:"name"`
	Type       string `gorm:"not null;default:'';comment:菜单类型" json:"type"`
	ParentID   uint   `gorm:"not null;default:0;comment:父id" json:"parent_id"`
	Path       string `gorm:"not null;default:'';comment:路由地址" json:"path"`
	Component  string `gorm:"not null;default:'';comment:组件地址" json:"component"`
	Params     string `gorm:"not null;default:'';comment:路由参数" json:"params"`
	Icon       string `gorm:"not null;default:'';comment:图标" json:"icon"`
	SortNum    int    `gorm:"not null;default:0;comment:排序" json:"sort_num"`
	Status     string `gorm:"not null;default:'';comment:显示隐藏状态" json:"status"`
}
