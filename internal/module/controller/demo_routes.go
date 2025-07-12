package controller

import (
	"github.com/zhoudm1743/go-frame/internal/module/schemas/req"
	"github.com/zhoudm1743/go-frame/internal/module/service"
	"github.com/zhoudm1743/go-frame/pkg/http"
	"github.com/zhoudm1743/go-frame/pkg/http/unified"
	"github.com/zhoudm1743/go-frame/pkg/response"
	"github.com/zhoudm1743/go-frame/util"
	"go.uber.org/fx"
)

// DemoHandler 示例控制器
type DemoHandler struct {
	service *service.DemoService
}

// NewDemoHandler 创建示例控制器
func NewDemoHandler(service *service.DemoService) *DemoHandler {
	return &DemoHandler{
		service: service,
	}
}

// RegisterDemoRoutes 注册示例路由
func RegisterDemoRoutes(server http.Server, handler *DemoHandler) {
	// 使用完整的路由路径，确保显示正确的URL
	router := server.Router().Group("/api/demoies")

	router.GET("", handler.List)
	router.GET("/:id", handler.Get)
	router.POST("", handler.Create)
	router.PUT("/:id", handler.Update)
	router.DELETE("/:id", handler.Delete)
	router.GET("/page", handler.Page)
}

// DemoModule 示例模块
var DemoModule = fx.Options(
	fx.Provide(NewDemoHandler),
	fx.Invoke(RegisterDemoRoutes),
)

// List 获取列表
func (c *DemoHandler) List(ctx unified.Context) error {
	// 这里不需要验证参数，直接获取所有数据
	items, err := c.service.GetAll()
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, items)
}

// Get 获取单个记录
func (c *DemoHandler) Get(ctx unified.Context) error {
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

// Create 创建记录
func (c *DemoHandler) Create(ctx unified.Context) error {
	var createReq req.DemoCreateReq
	if err := util.VerifyUtil.VerifyJSON(ctx, &createReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	created, err := c.service.Create(&createReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, created)
}

// Update 更新记录
func (c *DemoHandler) Update(ctx unified.Context) error {
	// 验证路径参数
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	// 验证请求体
	var updateReq req.DemoUpdateReq
	if err := util.VerifyUtil.VerifyJSON(ctx, &updateReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, err.Error())
	}

	updated, err := c.service.Update(idReq.ID, &updateReq)
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, updated)
}

// Delete 删除记录
func (c *DemoHandler) Delete(ctx unified.Context) error {
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

// Page 分页查询
func (c *DemoHandler) Page(ctx unified.Context) error {
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
