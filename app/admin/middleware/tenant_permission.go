package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/app/admin/service/system"
	"github.com/zhoudm1743/go-flow/pkg/jwt"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

// TenantPermissionMiddleware 多租户权限验证中间件
func TenantPermissionMiddleware(jwtService jwt.JwtService, permissionService system.TenantPermissionService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			response.FailWithMsg(c, response.TokenEmpty, "未提供认证令牌")
			c.Abort()
			return
		}

		// 解析token获取用户ID
		userID, err := jwtService.ParseToken(token)
		if err != nil {
			response.FailWithMsg(c, response.TokenInvalid, "认证令牌无效")
			c.Abort()
			return
		}

		// 获取租户ID（从header或query参数）
		tenantIDStr := c.GetHeader("X-Tenant-ID")
		if tenantIDStr == "" {
			tenantIDStr = c.Query("tenant_id")
		}
		if tenantIDStr == "" {
			response.FailWithMsg(c, response.ParamsValidError, "未提供租户ID")
			c.Abort()
			return
		}

		tenantID64, err := strconv.ParseUint(tenantIDStr, 10, 32)
		if err != nil {
			response.FailWithMsg(c, response.ParamsTypeError, "租户ID格式错误")
			c.Abort()
			return
		}
		tenantID := uint(tenantID64)

		// 获取请求的资源和操作
		resource := c.Request.URL.Path
		action := c.Request.Method

		// 检查权限
		hasPermission, err := permissionService.CheckPermission(tenantID, userID, resource, action)
		if err != nil {
			response.FailWithMsg(c, response.SystemError, fmt.Sprintf("权限检查失败: %v", err))
			c.Abort()
			return
		}

		if !hasPermission {
			response.FailWithMsg(c, response.NoPermission, "没有访问权限")
			c.Abort()
			return
		}

		// 将租户ID和用户ID存储到上下文中
		c.Set("tenant_id", tenantID)
		c.Set("user_id", userID)

		c.Next()
	})
}

// OptionalTenantPermissionMiddleware 可选的多租户权限验证中间件（用于公开接口）
func OptionalTenantPermissionMiddleware(jwtService jwt.JwtService, permissionService system.TenantPermissionService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 如果提供了token，则进行权限验证
		token := c.GetHeader("Authorization")
		if token != "" {
			// 执行完整的权限验证
			TenantPermissionMiddleware(jwtService, permissionService)(c)
			return
		}

		// 如果没有token，继续执行
		c.Next()
	})
}

// GetTenantID 从上下文获取租户ID
func GetTenantID(c *gin.Context) (uint, bool) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return 0, false
	}

	if id, ok := tenantID.(uint); ok {
		return id, true
	}

	return 0, false
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if id, ok := userID.(uint); ok {
		return id, true
	}

	return 0, false
}

// RequireTenantPermission 装饰器函数，用于快速添加权限验证到特定路由
func RequireTenantPermission(resource, action string, jwtService jwt.JwtService, permissionService system.TenantPermissionService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			response.FailWithMsg(c, response.TokenEmpty, "未提供认证令牌")
			c.Abort()
			return
		}

		// 解析token获取用户ID
		userID, err := jwtService.ParseToken(token)
		if err != nil {
			response.FailWithMsg(c, response.TokenInvalid, "认证令牌无效")
			c.Abort()
			return
		}

		// 获取租户ID
		tenantIDStr := c.GetHeader("X-Tenant-ID")
		if tenantIDStr == "" {
			tenantIDStr = c.Query("tenant_id")
		}
		if tenantIDStr == "" {
			response.FailWithMsg(c, response.ParamsValidError, "未提供租户ID")
			c.Abort()
			return
		}

		tenantID64, err := strconv.ParseUint(tenantIDStr, 10, 32)
		if err != nil {
			response.FailWithMsg(c, response.ParamsTypeError, "租户ID格式错误")
			c.Abort()
			return
		}
		tenantID := uint(tenantID64)

		// 检查特定权限
		hasPermission, err := permissionService.CheckPermission(tenantID, userID, resource, action)
		if err != nil {
			response.FailWithMsg(c, response.SystemError, fmt.Sprintf("权限检查失败: %v", err))
			c.Abort()
			return
		}

		if !hasPermission {
			response.FailWithMsg(c, response.NoPermission, fmt.Sprintf("没有访问 %s 的 %s 权限", resource, action))
			c.Abort()
			return
		}

		// 将租户ID和用户ID存储到上下文中
		c.Set("tenant_id", tenantID)
		c.Set("user_id", userID)

		c.Next()
	})
}
