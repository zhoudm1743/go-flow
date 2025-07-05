package controller

import (
	"github.com/zhoudm1743/go-flow/internal/test/schemas/req"
	"github.com/zhoudm1743/go-flow/internal/test/service"
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/util"
)

// ProductController 产品控制器
type ProductController struct {
	service *service.ProductService
}

// NewProductController 创建产品控制器
func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

// List 获取列表
func (c *ProductController) List(ctx unified.Context) error {
	// 这里不需要验证参数，直接获取所有数据
	items, err := c.service.GetAll()
	if err != nil {
		return response.UnifiedFailWithMsg(ctx, response.SystemError, err.Error())
	}
	return response.UnifiedOkWithData(ctx, items)
}

// Get 获取单个记录
func (c *ProductController) Get(ctx unified.Context) error {
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
func (c *ProductController) Create(ctx unified.Context) error {
	var createReq req.ProductCreateReq
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
func (c *ProductController) Update(ctx unified.Context) error {
	// 验证路径参数
	var idReq req.IdReq
	if err := util.VerifyUtil.VerifyQuery(ctx, &idReq); err != nil {
		return response.UnifiedFailWithMsg(ctx, response.ParamsValidError, "无效的ID")
	}

	// 验证请求体
	var updateReq req.ProductUpdateReq
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
func (c *ProductController) Delete(ctx unified.Context) error {
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
func (c *ProductController) Page(ctx unified.Context) error {
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
