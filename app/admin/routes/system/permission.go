package system

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/app/admin/middleware"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/req"
	"github.com/zhoudm1743/go-flow/app/admin/service/system"
	httpCore "github.com/zhoudm1743/go-flow/core/http"
	"github.com/zhoudm1743/go-flow/pkg/jwt"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/util"
)

type permissionRoutes struct {
	srv system.TenantPermissionService
	jwt jwt.JwtService
}

func NewPermissionGroup(srv system.TenantPermissionService, jwt jwt.JwtService) httpCore.Group {
	return httpCore.NewGroup("/system",
		func() interface{} {
			return &permissionRoutes{srv: srv, jwt: jwt}
		},
		regPermission,
		middleware.AuthMiddleware(jwt),
	)
}

func regPermission(rg *httpCore.BaseGroup, instance interface{}) error {
	r := instance.(*permissionRoutes)

	// 权限验证接口
	rg.POST("/permission/check", r.checkPermission)

	// 用户权限管理
	rg.POST("/user/assign-roles", r.assignRolesToUser)
	rg.POST("/user/remove-roles", r.removeRolesFromUser)
	rg.GET("/user/roles", r.getUserRoles)
	rg.GET("/user/permissions", r.getUserPermissions)

	// 角色权限管理
	rg.POST("/role/assign-permissions", r.assignPermissionsToRole)
	rg.POST("/role/remove-permissions", r.removePermissionsFromRole)
	rg.GET("/role/permissions", r.getRolePermissions)

	return nil
}

// @Summary 检查用户权限
// @Description 检查用户是否具有特定资源的操作权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param checkReq body req.SystemPermissionCheckReq true "权限检查参数"
// @Success 200 {object} response.Response{data=resp.SystemPermissionCheckResp} "成功"
// @Router /system/permission/check [post]
func (r *permissionRoutes) checkPermission(c *gin.Context) {
	var checkReq req.SystemPermissionCheckReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &checkReq)) {
		return
	}

	hasPermission, err := r.srv.CheckPermission(checkReq.TenantID, checkReq.UserID, checkReq.Resource, checkReq.Action)
	if err != nil {
		response.CheckAndResp(c, err)
		return
	}

	result := map[string]bool{"has_permission": hasPermission}
	response.CheckAndRespWithData(c, result, nil)
}

// @Summary 为用户分配角色
// @Description 为用户分配角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param assignReq body req.SystemTenantUserRoleReq true "角色分配参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/user/assign-roles [post]
func (r *permissionRoutes) assignRolesToUser(c *gin.Context) {
	var assignReq req.SystemTenantUserRoleReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &assignReq)) {
		return
	}
	err := r.srv.AssignRolesToUser(&assignReq)
	response.CheckAndResp(c, err)
}

// @Summary 移除用户角色
// @Description 移除用户的指定角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param removeReq body req.SystemTenantUserRoleReq true "角色移除参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/user/remove-roles [post]
func (r *permissionRoutes) removeRolesFromUser(c *gin.Context) {
	var removeReq req.SystemTenantUserRoleReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &removeReq)) {
		return
	}
	err := r.srv.RemoveRolesFromUser(removeReq.TenantID, removeReq.UserID, removeReq.RoleIDs)
	response.CheckAndResp(c, err)
}

// @Summary 获取用户角色
// @Description 获取用户拥有的角色列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param tenant_id query uint true "租户ID"
// @Param user_id query uint true "用户ID"
// @Success 200 {object} response.Response{data=[]resp.SystemTenantRoleResp} "成功"
// @Router /system/user/roles [get]
func (r *permissionRoutes) getUserRoles(c *gin.Context) {
	var queryReq struct {
		TenantID uint `form:"tenant_id" json:"tenant_id" validate:"required"`
		UserID   uint `form:"user_id" json:"user_id" validate:"required"`
	}
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &queryReq)) {
		return
	}
	res, err := r.srv.GetUserRoles(queryReq.TenantID, queryReq.UserID)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 获取用户权限
// @Description 获取用户拥有的权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param tenant_id query uint true "租户ID"
// @Param user_id query uint true "用户ID"
// @Success 200 {object} response.Response{data=[]resp.SystemTenantPermissionResp} "成功"
// @Router /system/user/permissions [get]
func (r *permissionRoutes) getUserPermissions(c *gin.Context) {
	var queryReq struct {
		TenantID uint `form:"tenant_id" json:"tenant_id" validate:"required"`
		UserID   uint `form:"user_id" json:"user_id" validate:"required"`
	}
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &queryReq)) {
		return
	}
	res, err := r.srv.GetUserPermissions(queryReq.TenantID, queryReq.UserID)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 为角色分配权限
// @Description 为角色分配权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param assignReq body req.SystemTenantRolePermissionReq true "权限分配参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/role/assign-permissions [post]
func (r *permissionRoutes) assignPermissionsToRole(c *gin.Context) {
	var assignReq req.SystemTenantRolePermissionReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &assignReq)) {
		return
	}
	err := r.srv.AssignPermissionsToRole(&assignReq)
	response.CheckAndResp(c, err)
}

// @Summary 移除角色权限
// @Description 移除角色的指定权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param removeReq body req.SystemTenantRolePermissionReq true "权限移除参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/role/remove-permissions [post]
func (r *permissionRoutes) removePermissionsFromRole(c *gin.Context) {
	var removeReq req.SystemTenantRolePermissionReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &removeReq)) {
		return
	}
	err := r.srv.RemovePermissionsFromRole(removeReq.TenantID, removeReq.RoleID, removeReq.PermissionIDs)
	response.CheckAndResp(c, err)
}

// @Summary 获取角色权限
// @Description 获取角色拥有的权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param tenant_id query uint true "租户ID"
// @Param role_id query uint true "角色ID"
// @Success 200 {object} response.Response{data=[]resp.SystemTenantPermissionResp} "成功"
// @Router /system/role/permissions [get]
func (r *permissionRoutes) getRolePermissions(c *gin.Context) {
	var queryReq struct {
		TenantID uint `form:"tenant_id" json:"tenant_id" validate:"required"`
		RoleID   uint `form:"role_id" json:"role_id" validate:"required"`
	}
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &queryReq)) {
		return
	}
	res, err := r.srv.GetRolePermissions(queryReq.TenantID, queryReq.RoleID)
	response.CheckAndRespWithData(c, res, err)
}
