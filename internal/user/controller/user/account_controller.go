package user

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/internal/user/schemas/req"
	"github.com/zhoudm1743/go-flow/internal/user/service/user"
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/util"
)

// AccountController 控制器
type AccountController struct {
	service user.AccountService
	logger  log.Logger
}

// NewAccountController 创建控制器
func NewAccountController(
	service user.AccountService,
	logger log.Logger,
) *AccountController {
	return &AccountController{
		service: service,
		logger:  logger,
	}
}

// AccountRouter 路由注册器
type AccountRouter struct {
	controller *AccountController
}

// NewAccountRouter 创建路由注册器
func NewAccountRouter(controller *AccountController) *AccountRouter {
	return &AccountRouter{
		controller: controller,
	}
}

// Register 注册路由
func (r *AccountRouter) Register(engine *gin.Engine) {
	// 默认使用无前缀注册
	r.RegisterWithPrefix(engine, "")
}

// RegisterWithPrefix 使用前缀注册路由
func (r *AccountRouter) RegisterWithPrefix(engine *gin.Engine, prefix string) {
	// 构建API路径
	apiPath := "/api/account"
	if prefix != "" {
		apiPath = prefix + apiPath
	}

	api := engine.Group(apiPath)
	{
		api.GET("", r.controller.List)
		api.GET("/:id", r.controller.Get)
		api.POST("", r.controller.Create)
		api.PUT("/:id", r.controller.Update)
		api.DELETE("/:id", r.controller.Delete)
		api.GET("/page", r.controller.Page)
	}
}

// RegisterRoutes 实现http.RouterRegister接口
func (r *AccountRouter) RegisterRoutes(router unified.Router) {
	// 使用统一的路由器API注册路由
	accountGroup := router.Group("/api/account")

	// 使用适配器将Gin处理函数转换为统一的处理函数
	accountGroup.GET("", adaptGinHandler(r.controller.List))
	accountGroup.GET("/:id", adaptGinHandler(r.controller.Get))
	accountGroup.POST("", adaptGinHandler(r.controller.Create))
	accountGroup.PUT("/:id", adaptGinHandler(r.controller.Update))
	accountGroup.DELETE("/:id", adaptGinHandler(r.controller.Delete))
	accountGroup.GET("/page", adaptGinHandler(r.controller.Page))
}

// adaptGinHandler 将Gin处理函数转换为统一的处理函数
func adaptGinHandler(handler func(*gin.Context)) unified.HandlerFunc {
	return func(ctx unified.Context) error {
		// 尝试获取Gin上下文
		if ginCtx, ok := ctx.GinContext().(*gin.Context); ok {
			handler(ginCtx)
			return nil
		}
		// 如果不是Gin上下文，返回错误
		return response.SystemError
	}
}

// ListUnified 获取列表（统一接口）
func (c *AccountController) ListUnified(ctx unified.Context) error {
	// 这里不需要验证参数，直接获取所有数据
	items, err := c.service.GetAll()
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, items)
}

// GetUnified 获取单个记录（统一接口）
func (c *AccountController) GetUnified(ctx unified.Context) error {
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	item, err := c.service.GetByID(idReq.ID)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, item)
}

// CreateUnified 创建记录（统一接口）
func (c *AccountController) CreateUnified(ctx unified.Context) error {
	var createReq req.AccountCreateReq
	if err := util.VerifyUtil.VerifyJSON(ctx, &createReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	created, err := c.service.Create(&createReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, created)
}

// UpdateUnified 更新记录（统一接口）
func (c *AccountController) UpdateUnified(ctx unified.Context) error {
	// 验证路径参数
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	// 验证请求体
	var updateReq req.AccountUpdateReq
	if err := util.VerifyUtil.VerifyJSON(ctx, &updateReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	updated, err := c.service.Update(idReq.ID, &updateReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, updated)
}

// DeleteUnified 删除记录（统一接口）
func (c *AccountController) DeleteUnified(ctx unified.Context) error {
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	err := c.service.Delete(idReq.ID)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOk(ctx)
}

// PageUnified 分页查询（统一接口）
func (c *AccountController) PageUnified(ctx unified.Context) error {
	var pageReq req.PageReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &pageReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	pageResult, err := c.service.GetPage(&pageReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, pageResult)
}

// List 获取列表
func (c *AccountController) List(ctx *gin.Context) {
	// 这里不需要验证参数，直接获取所有数据
	items, err := c.service.GetAll()
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, err.Error())
		return
	}
	response.OkWithData(ctx, items)
}

// Get 获取单个记录
func (c *AccountController) Get(ctx *gin.Context) {
	var idReq req.IdReq
	if err := ctx.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "无效的ID")
		return
	}

	item, err := c.service.GetByID(idReq.ID)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, err.Error())
		return
	}
	response.OkWithData(ctx, item)
}

// Create 创建记录
func (c *AccountController) Create(ctx *gin.Context) {
	var createReq req.AccountCreateReq
	if err := ctx.ShouldBindJSON(&createReq); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, err.Error())
		return
	}

	created, err := c.service.Create(&createReq)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, err.Error())
		return
	}
	response.OkWithData(ctx, created)
}

// Update 更新记录
func (c *AccountController) Update(ctx *gin.Context) {
	// 验证路径参数
	var idReq req.IdReq
	if err := ctx.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "无效的ID")
		return
	}

	// 验证请求体
	var updateReq req.AccountUpdateReq
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, err.Error())
		return
	}

	updated, err := c.service.Update(idReq.ID, &updateReq)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, err.Error())
		return
	}
	response.OkWithData(ctx, updated)
}

// Delete 删除记录
func (c *AccountController) Delete(ctx *gin.Context) {
	var idReq req.IdReq
	if err := ctx.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "无效的ID")
		return
	}

	err := c.service.Delete(idReq.ID)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, err.Error())
		return
	}
	response.Ok(ctx)
}

// Page 分页查询
func (c *AccountController) Page(ctx *gin.Context) {
	var pageReq req.PageReq
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, err.Error())
		return
	}

	pageResult, err := c.service.GetPage(&pageReq)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, err.Error())
		return
	}
	response.OkWithData(ctx, pageResult)
}
