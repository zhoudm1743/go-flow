package controller

import (
	"github.com/zhoudm1743/go-flow/internal/demo/service"
	"github.com/zhoudm1743/go-flow/pkg/log"
)

// 请求结构体
type (
	// UserCreateReq 创建用户请求
	UserCreateReq struct {
		Name  string `json:"name" binding:"required" label:"用户名"`
		Email string `json:"email" binding:"required,email" label:"邮箱"`
	}

	// UserUpdateReq 更新用户请求
	UserUpdateReq struct {
		Name  string `json:"name" binding:"omitempty" label:"用户名"`
		Email string `json:"email" binding:"omitempty,email" label:"邮箱"`
	}
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
	logger      log.Logger
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService, logger log.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

// GetUserService 获取用户服务
func (c *UserController) GetUserService() service.UserService {
	return c.userService
}
