package req

type SystemAdminListReq struct {
	Username string `form:"username" json:"username"`
	Nickname string `form:"nickname" json:"nickname"`
	Status   int8   `form:"status,default=-1" json:"status"`
	Role     string `form:"role" json:"role"`
}

type SystemAdminAddReq struct {
	Username string   `form:"username" json:"username" validate:"required,min=3,max=30"`
	Nickname string   `form:"nickname" json:"nickname" validate:"required,min=3,max=30"`
	Password string   `form:"password" json:"password" validate:"required,min=6,max=30"`
	Roles    []string `form:"roles" json:"roles" validate:"required"`
	Email    string   `form:"email" json:"email" validate:"required,email"`
	Phone    string   `form:"phone" json:"phone" validate:"required,phone"`
	Avatar   string   `form:"avatar" json:"avatar"`
}

type SystemAdminEditReq struct {
	ID       uint     `form:"id" json:"id" validate:"required"`
	Username string   `form:"username" json:"username" validate:"required,min=3,max=30"`
	Nickname string   `form:"nickname" json:"nickname" validate:"required,min=3,max=30"`
	Roles    []string `form:"roles" json:"roles" validate:"required"`
	Email    string   `form:"email" json:"email" validate:"required,email"`
	Phone    string   `form:"phone" json:"phone" validate:"required,phone"`
	Avatar   string   `form:"avatar" json:"avatar"`
}

// 租户管理相关请求
type SystemTenantListReq struct {
	Name   string `form:"name" json:"name"`
	Code   string `form:"code" json:"code"`
	Type   string `form:"type" json:"type"`
	Status int8   `form:"status,default=-1" json:"status"`
}

type SystemTenantAddReq struct {
	Name             string `form:"name" json:"name" validate:"required,min=2,max=100"`
	Code             string `form:"code" json:"code" validate:"required,min=2,max=50"`
	Type             string `form:"type" json:"type" validate:"required"`
	Status           int8   `form:"status" json:"status" validate:"required,oneof=0 1"`
	ExpireAt         int64  `form:"expire_at" json:"expire_at"`
	DatabaseHost     string `form:"database_host" json:"database_host" validate:"required"`
	DatabasePort     string `form:"database_port" json:"database_port" validate:"required"`
	DatabaseName     string `form:"database_name" json:"database_name" validate:"required"`
	DatabaseUser     string `form:"database_user" json:"database_user" validate:"required"`
	DatabasePassword string `form:"database_password" json:"database_password" validate:"required"`
}

type SystemTenantEditReq struct {
	ID               uint   `form:"id" json:"id" validate:"required"`
	Name             string `form:"name" json:"name" validate:"required,min=2,max=100"`
	Code             string `form:"code" json:"code" validate:"required,min=2,max=50"`
	Type             string `form:"type" json:"type" validate:"required"`
	Status           int8   `form:"status" json:"status" validate:"required,oneof=0 1"`
	ExpireAt         int64  `form:"expire_at" json:"expire_at"`
	DatabaseHost     string `form:"database_host" json:"database_host" validate:"required"`
	DatabasePort     string `form:"database_port" json:"database_port" validate:"required"`
	DatabaseName     string `form:"database_name" json:"database_name" validate:"required"`
	DatabaseUser     string `form:"database_user" json:"database_user" validate:"required"`
	DatabasePassword string `form:"database_password" json:"database_password"`
}

// 租户角色相关请求
type SystemTenantRoleListReq struct {
	TenantID uint   `form:"tenant_id" json:"tenant_id"`
	Name     string `form:"name" json:"name"`
	Code     string `form:"code" json:"code"`
	Status   int8   `form:"status,default=-1" json:"status"`
}

type SystemTenantRoleAddReq struct {
	TenantID    uint   `form:"tenant_id" json:"tenant_id" validate:"required"`
	Name        string `form:"name" json:"name" validate:"required,min=2,max=50"`
	Code        string `form:"code" json:"code" validate:"required,min=2,max=50"`
	Description string `form:"description" json:"description"`
	Status      int8   `form:"status" json:"status" validate:"required,oneof=0 1"`
}

type SystemTenantRoleEditReq struct {
	ID          uint   `form:"id" json:"id" validate:"required"`
	TenantID    uint   `form:"tenant_id" json:"tenant_id" validate:"required"`
	Name        string `form:"name" json:"name" validate:"required,min=2,max=50"`
	Code        string `form:"code" json:"code" validate:"required,min=2,max=50"`
	Description string `form:"description" json:"description"`
	Status      int8   `form:"status" json:"status" validate:"required,oneof=0 1"`
}

// 租户权限相关请求
type SystemTenantPermissionListReq struct {
	TenantID uint   `form:"tenant_id" json:"tenant_id"`
	Name     string `form:"name" json:"name"`
	Code     string `form:"code" json:"code"`
	Resource string `form:"resource" json:"resource"`
	Action   string `form:"action" json:"action"`
	Status   int8   `form:"status,default=-1" json:"status"`
}

type SystemTenantPermissionAddReq struct {
	TenantID    uint   `form:"tenant_id" json:"tenant_id" validate:"required"`
	Name        string `form:"name" json:"name" validate:"required,min=2,max=50"`
	Code        string `form:"code" json:"code" validate:"required,min=2,max=50"`
	Resource    string `form:"resource" json:"resource" validate:"required"`
	Action      string `form:"action" json:"action" validate:"required"`
	Description string `form:"description" json:"description"`
	Status      int8   `form:"status" json:"status" validate:"required,oneof=0 1"`
}

type SystemTenantPermissionEditReq struct {
	ID          uint   `form:"id" json:"id" validate:"required"`
	TenantID    uint   `form:"tenant_id" json:"tenant_id" validate:"required"`
	Name        string `form:"name" json:"name" validate:"required,min=2,max=50"`
	Code        string `form:"code" json:"code" validate:"required,min=2,max=50"`
	Resource    string `form:"resource" json:"resource" validate:"required"`
	Action      string `form:"action" json:"action" validate:"required"`
	Description string `form:"description" json:"description"`
	Status      int8   `form:"status" json:"status" validate:"required,oneof=0 1"`
}

// 角色权限分配请求
type SystemTenantRolePermissionReq struct {
	TenantID      uint   `form:"tenant_id" json:"tenant_id" validate:"required"`
	RoleID        uint   `form:"role_id" json:"role_id" validate:"required"`
	PermissionIDs []uint `form:"permission_ids" json:"permission_ids" validate:"required"`
}

// 用户角色分配请求
type SystemTenantUserRoleReq struct {
	TenantID uint   `form:"tenant_id" json:"tenant_id" validate:"required"`
	UserID   uint   `form:"user_id" json:"user_id" validate:"required"`
	RoleIDs  []uint `form:"role_ids" json:"role_ids" validate:"required"`
}

// 权限验证请求
type SystemPermissionCheckReq struct {
	TenantID uint   `form:"tenant_id" json:"tenant_id" validate:"required"`
	UserID   uint   `form:"user_id" json:"user_id" validate:"required"`
	Resource string `form:"resource" json:"resource" validate:"required"`
	Action   string `form:"action" json:"action" validate:"required"`
}
