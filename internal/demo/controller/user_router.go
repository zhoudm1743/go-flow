package controller

import (
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/validate"
)

// UserRouter 用户路由注册器
type UserRouter struct {
	controller *UserController
}

// NewUserRouter 创建用户路由注册器
func NewUserRouter(controller *UserController) *UserRouter {
	return &UserRouter{
		controller: controller,
	}
}

// RegisterRoutes 注册路由到路由器
func (r *UserRouter) RegisterRoutes(router unified.Router) {
	// 注册API路由
	api := router.Group("/api/users")

	// 获取所有用户
	api.GET("/", r.ListUsers)

	// 获取单个用户
	api.GET("/:id", r.GetUser)

	// 创建用户
	api.POST("/", r.CreateUser)

	// 更新用户
	api.PUT("/:id", r.UpdateUser)

	// 删除用户
	api.DELETE("/:id", r.DeleteUser)
}

// ListUsers 获取所有用户
func (r *UserRouter) ListUsers(ctx unified.Context) error {
	users, err := r.controller.userService.GetAllUsers()
	return response.UnifiedCheckAndRespWithData(ctx, users, err)
}

// GetUser 获取单个用户
func (r *UserRouter) GetUser(ctx unified.Context) error {
	// 获取路径参数
	idStr := ctx.Param("id")
	id, err := ctx.ParamUint("id")
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的用户ID: "+idStr)
	}

	user, err := r.controller.userService.GetUserByID(id)
	return response.UnifiedCheckAndRespWithData(ctx, user, err)
}

// CreateUser 创建用户
func (r *UserRouter) CreateUser(ctx unified.Context) error {
	var req UserCreateReq
	if err := validate.UnifiedVerify.VerifyJSON(ctx, &req); err != nil {
		return err
	}

	createdUser, err := r.controller.userService.CreateUser(req.Name, req.Email)
	return response.UnifiedCheckAndRespWithData(ctx, createdUser, err)
}

// UpdateUser 更新用户
func (r *UserRouter) UpdateUser(ctx unified.Context) error {
	// 获取路径参数
	id, err := ctx.ParamUint("id")
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的用户ID")
	}

	// 验证请求体
	var req UserUpdateReq
	if err := validate.UnifiedVerify.VerifyJSON(ctx, &req); err != nil {
		return err
	}

	updatedUser, err := r.controller.userService.UpdateUser(id, req.Name, req.Email)
	return response.UnifiedCheckAndRespWithData(ctx, updatedUser, err)
}

// DeleteUser 删除用户
func (r *UserRouter) DeleteUser(ctx unified.Context) error {
	// 获取路径参数
	id, err := ctx.ParamUint("id")
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的用户ID")
	}

	err = r.controller.userService.DeleteUser(id)
	return response.UnifiedCheckAndResp(ctx, err)
}
