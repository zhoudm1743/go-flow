package controller

import (
	"github.com/zhoudm1743/go-flow/internal/shop/schemas/req"
	"github.com/zhoudm1743/go-flow/internal/shop/service"
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/util"
)

// CategoryController 分类控制器
type CategoryController struct {
	service *service.CategoryService
}

// NewCategoryController 创建分类控制器
func NewCategoryController(service *service.CategoryService) *CategoryController {
	return &CategoryController{
		service: service,
	}
}

// List 获取列表
func (c *CategoryController) List(ctx unified.Context) error {
	// 这里不需要验证参数，直接获取所有数据
	items, err := c.service.GetAll()
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, items)
}

// Get 获取单个记录
func (c *CategoryController) Get(ctx unified.Context) error {
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
func (c *CategoryController) Create(ctx unified.Context) error {
	var createReq req.CategoryCreateReq
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
func (c *CategoryController) Update(ctx unified.Context) error {
	// 验证路径参数
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	// 验证请求体
	var updateReq req.CategoryUpdateReq
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
func (c *CategoryController) Delete(ctx unified.Context) error {
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
func (c *CategoryController) Page(ctx unified.Context) error {
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
