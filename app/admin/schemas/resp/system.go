package resp

import "github.com/zhoudm1743/go-flow/pkg/types"

type SystemAdminResp struct {
	ID          uint         `json:"id" structs:"id"`
	Username    string       `json:"username" structs:"username"`
	Nickname    string       `json:"nickname" structs:"nickname"`
	Status      int8         `json:"status" structs:"status"`
	Roles       []string     `json:"roles" structs:"roles"`
	Email       string       `json:"email" structs:"email"`
	Phone       string       `json:"phone" structs:"phone"`
	Avatar      string       `json:"avatar" structs:"avatar"`
	LastLoginAt types.TsTime `json:"last_login_at" structs:"last_login_at"`
	LastLoginIP string       `json:"last_login_ip" structs:"last_login_ip"`
	CreatedAt   types.TsTime `json:"created_at" structs:"created_at"`
	UpdatedAt   types.TsTime `json:"updated_at" structs:"updated_at"`
}

// 租户管理相关响应
type SystemTenantResp struct {
	ID           uint         `json:"id" structs:"id"`
	Name         string       `json:"name" structs:"name"`
	Code         string       `json:"code" structs:"code"`
	Type         string       `json:"type" structs:"type"`
	Status       int8         `json:"status" structs:"status"`
	ExpireAt     int64        `json:"expire_at" structs:"expire_at"`
	DatabaseHost string       `json:"database_host" structs:"database_host"`
	DatabasePort string       `json:"database_port" structs:"database_port"`
	DatabaseName string       `json:"database_name" structs:"database_name"`
	DatabaseUser string       `json:"database_user" structs:"database_user"`
	CreatedAt    types.TsTime `json:"created_at" structs:"created_at"`
	UpdatedAt    types.TsTime `json:"updated_at" structs:"updated_at"`
}

type SystemTenantRoleResp struct {
	ID          uint         `json:"id" structs:"id"`
	TenantID    uint         `json:"tenant_id" structs:"tenant_id"`
	Name        string       `json:"name" structs:"name"`
	Code        string       `json:"code" structs:"code"`
	Description string       `json:"description" structs:"description"`
	Status      int8         `json:"status" structs:"status"`
	CreatedAt   types.TsTime `json:"created_at" structs:"created_at"`
	UpdatedAt   types.TsTime `json:"updated_at" structs:"updated_at"`
}

type SystemTenantPermissionResp struct {
	ID          uint         `json:"id" structs:"id"`
	TenantID    uint         `json:"tenant_id" structs:"tenant_id"`
	Name        string       `json:"name" structs:"name"`
	Code        string       `json:"code" structs:"code"`
	Resource    string       `json:"resource" structs:"resource"`
	Action      string       `json:"action" structs:"action"`
	Description string       `json:"description" structs:"description"`
	Status      int8         `json:"status" structs:"status"`
	CreatedAt   types.TsTime `json:"created_at" structs:"created_at"`
	UpdatedAt   types.TsTime `json:"updated_at" structs:"updated_at"`
}

type SystemTenantRolePermissionResp struct {
	ID             uint         `json:"id" structs:"id"`
	TenantID       uint         `json:"tenant_id" structs:"tenant_id"`
	RoleID         uint         `json:"role_id" structs:"role_id"`
	PermissionID   uint         `json:"permission_id" structs:"permission_id"`
	RoleName       string       `json:"role_name" structs:"role_name"`
	PermissionName string       `json:"permission_name" structs:"permission_name"`
	Resource       string       `json:"resource" structs:"resource"`
	Action         string       `json:"action" structs:"action"`
	CreatedAt      types.TsTime `json:"created_at" structs:"created_at"`
	UpdatedAt      types.TsTime `json:"updated_at" structs:"updated_at"`
}

type SystemTenantUserRoleResp struct {
	ID        uint         `json:"id" structs:"id"`
	TenantID  uint         `json:"tenant_id" structs:"tenant_id"`
	UserID    uint         `json:"user_id" structs:"user_id"`
	RoleID    uint         `json:"role_id" structs:"role_id"`
	Username  string       `json:"username" structs:"username"`
	RoleName  string       `json:"role_name" structs:"role_name"`
	CreatedAt types.TsTime `json:"created_at" structs:"created_at"`
	UpdatedAt types.TsTime `json:"updated_at" structs:"updated_at"`
}

// 权限验证响应
type SystemPermissionCheckResp struct {
	HasPermission bool `json:"has_permission"`
}

// 用户权限列表响应
type SystemUserPermissionResp struct {
	UserID      uint     `json:"user_id"`
	TenantID    uint     `json:"tenant_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
