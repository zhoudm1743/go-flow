package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/pkg/jwt"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		token, err := jwtService.GetTokenFromContext(c)
		if err != nil {
			response.FailWithMsg(c, response.TokenEmpty, err.Error())
			c.Abort()
			return
		}

		// 验证token
		isValid, err := jwtService.VerifyToken(token)
		if err != nil || !isValid {
			response.FailWithMsg(c, response.TokenInvalid, "token验证失败")
			c.Abort()
			return
		}

		// 解析token获取用户信息
		userID, err := jwtService.ParseToken(token)
		if err != nil {
			response.FailWithMsg(c, response.TokenInvalid, "token解析失败")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", userID)
		c.Set("token", token)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求token）
func OptionalAuthMiddleware(jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头中获取token
		token, err := jwtService.GetTokenFromContext(c)
		if err != nil {
			// 没有token，继续处理请求
			c.Next()
			return
		}

		// 验证token
		isValid, err := jwtService.VerifyToken(token)
		if err != nil || !isValid {
			// token无效，继续处理请求但不设置用户信息
			c.Next()
			return
		}

		// 解析token获取用户信息
		userID, err := jwtService.ParseToken(token)
		if err != nil {
			// token解析失败，继续处理请求但不设置用户信息
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", userID)
		c.Set("token", token)

		c.Next()
	}
}

// GetUserIDFromContext 从gin上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		return 0, false
	}

	return userID, true
}

// GetTokenFromContext 从gin上下文中获取token
func GetTokenFromContext(c *gin.Context) (string, bool) {
	tokenInterface, exists := c.Get("token")
	if !exists {
		return "", false
	}

	token, ok := tokenInterface.(string)
	if !ok {
		return "", false
	}

	return token, true
}

// RequireUserMiddleware 要求用户登录中间件（配合OptionalAuthMiddleware使用）
func RequireUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请先登录",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnlyMiddleware 仅管理员访问中间件
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists || userID == 0 {
			response.FailWithMsg(c, response.TokenEmpty, "请先登录")
			c.Abort()
			return
		}

		// TODO: 这里应该查询数据库检查用户是否为管理员
		// 现在先简单判断用户ID是否为1（假设ID为1的是管理员）
		if userID != 1 {
			response.Fail(c, response.NoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RefreshTokenHandler 刷新token处理器
func RefreshTokenHandler(jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		token, err := jwtService.GetTokenFromContext(c)
		if err != nil {
			response.FailWithMsg(c, response.TokenEmpty, err.Error())
			return
		}

		// 刷新token
		newToken, err := jwtService.RefreshToken(token)
		if err != nil {
			response.FailWithMsg(c, response.TokenExpired, "token刷新失败")
			return
		}

		// 返回新token
		response.OkWithData(c, gin.H{
			"token": newToken,
		})
	}
}

/*
使用示例：

1. 在路由中使用认证中间件：

func SetupRoutes(r *gin.Engine, jwtService jwt.JwtService) {
	// 公开接口（无需认证）
	public := r.Group("/api/public")
	{
		public.POST("/login", loginHandler)
		public.POST("/register", registerHandler)
	}

	// 需要认证的接口
	auth := r.Group("/api/auth")
	auth.Use(AuthMiddleware(jwtService)) // 使用认证中间件
	{
		auth.GET("/profile", getUserProfile)
		auth.PUT("/profile", updateUserProfile)
		auth.POST("/refresh", RefreshTokenHandler(jwtService))
	}

	// 需要管理员权限的接口
	admin := r.Group("/api/admin")
	admin.Use(AuthMiddleware(jwtService)) // 先认证
	admin.Use(AdminOnlyMiddleware())      // 再检查管理员权限
	{
		admin.GET("/users", getUserList)
		admin.POST("/users", createUser)
		admin.DELETE("/users/:id", deleteUser)
	}

	// 可选认证的接口（有token时提供用户信息，没有token时也可以访问）
	optional := r.Group("/api/optional")
	optional.Use(OptionalAuthMiddleware(jwtService))
	{
		optional.GET("/posts", getPostList) // 登录用户可以看到更多信息
	}

	// 需要登录但使用两步验证的接口
	twoStep := r.Group("/api/user")
	twoStep.Use(OptionalAuthMiddleware(jwtService)) // 先尝试获取用户信息
	twoStep.Use(RequireUserMiddleware())            // 然后要求必须登录
	{
		twoStep.GET("/dashboard", getUserDashboard)
	}
}

2. 在处理器中获取用户信息：

func getUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.FailWithMsg(c, response.SystemError, "获取用户信息失败")
		return
	}

	// 使用用户ID查询数据库
	// user := userService.GetByID(userID)

	response.OkWithData(c, gin.H{
		"user_id": userID,
		"message": "获取用户信息成功",
	})
}

3. 登录接口示例：

func loginHandler(c *gin.Context, jwtService jwt.JwtService) {
	// 验证用户凭据
	// ... 验证逻辑 ...

	userID := uint(123) // 假设验证成功，用户ID为123

	// 生成JWT token
	token, err := jwtService.GenerateToken(userID)
	if err != nil {
		response.FailWithMsg(c, response.SystemError, "生成token失败")
		return
	}

	response.OkWithData(c, gin.H{
		"token": token,
		"user_id": userID,
	})
}

4. 前端使用token：

在前端请求时，需要在请求头中添加Authorization字段：

Header: Authorization: Bearer <your-jwt-token>

5. Fx依赖注入示例：

在您的应用模块中，可以这样注入JWT服务：

func NewAuthRoutes(jwtService jwt.JwtService) *gin.RouterGroup {
	// 使用注入的jwtService创建路由
}
*/
