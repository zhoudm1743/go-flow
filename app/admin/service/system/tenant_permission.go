package system

import (
	"github.com/zhoudm1743/go-flow/app/admin/schemas/req"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/resp"
	"github.com/zhoudm1743/go-flow/app/models"
	"github.com/zhoudm1743/go-flow/core/database"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

// TenantPermissionService 租户权限管理服务
type TenantPermissionService interface {
	// 权限验证
	CheckPermission(tenantID, userID uint, resource, action string) (bool, error)

	// 用户权限管理
	AssignRolesToUser(req *req.SystemTenantUserRoleReq) error
	RemoveRolesFromUser(tenantID, userID uint, roleIDs []uint) error
	GetUserRoles(tenantID, userID uint) ([]resp.SystemTenantRoleResp, error)
	GetUserPermissions(tenantID, userID uint) ([]resp.SystemTenantPermissionResp, error)

	// 角色权限管理
	AssignPermissionsToRole(req *req.SystemTenantRolePermissionReq) error
	RemovePermissionsFromRole(tenantID, roleID uint, permissionIDs []uint) error
	GetRolePermissions(tenantID, roleID uint) ([]resp.SystemTenantPermissionResp, error)
}

type tenantPermissionService struct {
	db              database.Database
	tenantDBService TenantDatabaseService
}

// CheckPermission 检查用户权限
func (t *tenantPermissionService) CheckPermission(tenantID, userID uint, resource, action string) (bool, error) {
	tenantDB, err := t.tenantDBService.GetTenantDB(tenantID)
	if err != nil {
		return false, err
	}

	// 通过用户角色查询权限
	var count int64
	err = tenantDB.Table("system_tenant_permissions").
		Joins("JOIN system_tenant_role_permissions ON system_tenant_permissions.id = system_tenant_role_permissions.permission_id").
		Joins("JOIN system_tenant_user_roles ON system_tenant_role_permissions.role_id = system_tenant_user_roles.role_id").
		Where("system_tenant_user_roles.tenant_id = ? AND system_tenant_user_roles.user_id = ?", tenantID, userID).
		Where("system_tenant_permissions.resource = ? AND system_tenant_permissions.action = ?", resource, action).
		Where("system_tenant_permissions.status = 1").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// AssignRolesToUser 为用户分配角色
func (t *tenantPermissionService) AssignRolesToUser(req *req.SystemTenantUserRoleReq) error {
	tenantDB, err := t.tenantDBService.GetTenantDB(req.TenantID)
	if err != nil {
		return err
	}

	// 开始事务
	tx := tenantDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户现有角色
	if err := tx.Where("tenant_id = ? AND user_id = ?", req.TenantID, req.UserID).
		Delete(&models.SystemTenantUserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 分配新角色
	for _, roleID := range req.RoleIDs {
		userRole := &models.SystemTenantUserRole{
			TenantID: req.TenantID,
			UserID:   req.UserID,
			RoleID:   roleID,
		}
		if err := tx.Create(userRole).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemoveRolesFromUser 移除用户角色
func (t *tenantPermissionService) RemoveRolesFromUser(tenantID, userID uint, roleIDs []uint) error {
	tenantDB, err := t.tenantDBService.GetTenantDB(tenantID)
	if err != nil {
		return err
	}

	// 删除指定角色
	return tenantDB.Where("tenant_id = ? AND user_id = ? AND role_id IN ?",
		tenantID, userID, roleIDs).Delete(&models.SystemTenantUserRole{}).Error
}

// GetUserRoles 获取用户角色
func (t *tenantPermissionService) GetUserRoles(tenantID, userID uint) ([]resp.SystemTenantRoleResp, error) {
	tenantDB, err := t.tenantDBService.GetTenantDB(tenantID)
	if err != nil {
		return nil, err
	}

	var roles []models.SystemTenantRole
	err = tenantDB.Table("system_tenant_roles").
		Joins("JOIN system_tenant_user_roles ON system_tenant_roles.id = system_tenant_user_roles.role_id").
		Where("system_tenant_user_roles.tenant_id = ? AND system_tenant_user_roles.user_id = ?",
			tenantID, userID).
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	var respRoles []resp.SystemTenantRoleResp
	response.Copy(&respRoles, roles)
	return respRoles, nil
}

// GetUserPermissions 获取用户权限
func (t *tenantPermissionService) GetUserPermissions(tenantID, userID uint) ([]resp.SystemTenantPermissionResp, error) {
	tenantDB, err := t.tenantDBService.GetTenantDB(tenantID)
	if err != nil {
		return nil, err
	}

	var permissions []models.SystemTenantPermission
	err = tenantDB.Table("system_tenant_permissions").
		Joins("JOIN system_tenant_role_permissions ON system_tenant_permissions.id = system_tenant_role_permissions.permission_id").
		Joins("JOIN system_tenant_user_roles ON system_tenant_role_permissions.role_id = system_tenant_user_roles.role_id").
		Where("system_tenant_user_roles.tenant_id = ? AND system_tenant_user_roles.user_id = ?",
			tenantID, userID).
		Group("system_tenant_permissions.id").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	var respPermissions []resp.SystemTenantPermissionResp
	response.Copy(&respPermissions, permissions)
	return respPermissions, nil
}

// AssignPermissionsToRole 为角色分配权限
func (t *tenantPermissionService) AssignPermissionsToRole(req *req.SystemTenantRolePermissionReq) error {
	tenantDB, err := t.tenantDBService.GetTenantDB(req.TenantID)
	if err != nil {
		return err
	}

	// 开始事务
	tx := tenantDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除角色现有权限
	if err := tx.Where("tenant_id = ? AND role_id = ?", req.TenantID, req.RoleID).
		Delete(&models.SystemTenantRolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 分配新权限
	for _, permissionID := range req.PermissionIDs {
		rolePermission := &models.SystemTenantRolePermission{
			TenantID:     req.TenantID,
			RoleID:       req.RoleID,
			PermissionID: permissionID,
		}
		if err := tx.Create(rolePermission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemovePermissionsFromRole 移除角色权限
func (t *tenantPermissionService) RemovePermissionsFromRole(tenantID, roleID uint, permissionIDs []uint) error {
	tenantDB, err := t.tenantDBService.GetTenantDB(tenantID)
	if err != nil {
		return err
	}

	// 删除指定权限
	return tenantDB.Where("tenant_id = ? AND role_id = ? AND permission_id IN ?",
		tenantID, roleID, permissionIDs).Delete(&models.SystemTenantRolePermission{}).Error
}

// GetRolePermissions 获取角色权限
func (t *tenantPermissionService) GetRolePermissions(tenantID, roleID uint) ([]resp.SystemTenantPermissionResp, error) {
	tenantDB, err := t.tenantDBService.GetTenantDB(tenantID)
	if err != nil {
		return nil, err
	}

	var permissions []models.SystemTenantPermission
	err = tenantDB.Table("system_tenant_permissions").
		Joins("JOIN system_tenant_role_permissions ON system_tenant_permissions.id = system_tenant_role_permissions.permission_id").
		Where("system_tenant_role_permissions.tenant_id = ? AND system_tenant_role_permissions.role_id = ?",
			tenantID, roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	var respPermissions []resp.SystemTenantPermissionResp
	response.Copy(&respPermissions, permissions)
	return respPermissions, nil
}

// NewTenantPermissionService 创建租户权限服务
func NewTenantPermissionService(db database.Database, tenantDBService TenantDatabaseService) TenantPermissionService {
	return &tenantPermissionService{
		db:              db,
		tenantDBService: tenantDBService,
	}
}
