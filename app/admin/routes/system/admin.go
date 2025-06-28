package system

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/req"
	"github.com/zhoudm1743/go-flow/app/admin/service/system"
	httpCore "github.com/zhoudm1743/go-flow/core/http"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/util"
)

type adminRoutes struct {
	srv system.AdminService
}

func NewAdminGroup(srv system.AdminService) httpCore.Group {
	return httpCore.NewGroup("/system",
		func() interface{} {
			return &adminRoutes{srv: srv}
		},
		regAdmin,
		httpCore.AuthMiddleware,
	)
}

func regAdmin(rg *httpCore.BaseGroup, instance interface{}) error {
	r := instance.(*adminRoutes)
	rg.GET("/admin/all", r.all)
	rg.GET("/admin/list", r.list)
	rg.GET("/admin/detail", r.detail)
	rg.POST("/admin/add", r.add)
	rg.POST("/admin/edit", r.edit)
	rg.POST("/admin/delete", r.delete)
	rg.POST("/admin/change-status", r.changeStatus)
	return nil
}

// @Summary 获取所有管理员
// @Description 获取所有管理员
// @Tags 管理员
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param queryReq body req.SystemAdminListReq true "查询参数"
// @Success 200 {object} response.Response{data=[]resp.SystemAdminResp} "成功"
// @Router /system/admin/all [get]
func (r *adminRoutes) all(c *gin.Context) {
	var queryReq req.SystemAdminListReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &queryReq)) {
		return
	}
	res, err := r.srv.All(&queryReq)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 获取管理员列表
// @Description 获取管理员列表
// @Tags 管理员
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param pageReq body req.PageReq true "分页参数"
// @Param queryReq body req.SystemAdminListReq true "查询参数"
// @Success 200 {object} response.Response{data=response.PageResp{lists=[]resp.SystemAdminResp}} "成功"
// @Router /system/admin/list [get]
func (r *adminRoutes) list(c *gin.Context) {
	var pageReq req.PageReq
	var queryReq req.SystemAdminListReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &pageReq, &queryReq)) {
		return
	}
	res, err := r.srv.List(&pageReq, &queryReq)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 获取管理员详情
// @Description 获取管理员详情
// @Tags 管理员
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param idReq body req.SystemAdminDetailReq true "管理员ID"
// @Success 200 {object} response.Response{data=resp.SystemAdminResp} "成功"
// @Router /system/admin/detail [get]
func (r *adminRoutes) detail(c *gin.Context) {
	var idReq req.IdReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &idReq)) {
		return
	}
	res, err := r.srv.Find(idReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 添加管理员
// @Description 添加管理员
// @Tags 管理员
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param addReq body req.SystemAdminAddReq true "添加参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/admin/add [post]
func (r *adminRoutes) add(c *gin.Context) {
	var addReq req.SystemAdminAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &addReq)) {
		return
	}
	err := r.srv.Add(&addReq)
	response.CheckAndResp(c, err)
}

// @Summary 编辑管理员
// @Description 编辑管理员
// @Tags 管理员
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param editReq body req.SystemAdminEditReq true "编辑参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/admin/edit [post]
func (r *adminRoutes) edit(c *gin.Context) {
	var editReq req.SystemAdminEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &editReq)) {
		return
	}
	err := r.srv.Edit(&editReq)
	response.CheckAndResp(c, err)
}

// @Summary 删除管理员
// @Description 删除管理员
// @Tags 管理员
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param idReq body req.IdReq true "管理员ID"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/admin/delete [post]
func (r *adminRoutes) delete(c *gin.Context) {
	var idReq req.IdReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &idReq)) {
		return
	}
	err := r.srv.Delete(idReq.ID)
	response.CheckAndResp(c, err)
}

// @Summary 修改管理员状态
// @Description 修改管理员状态
// @Tags 管理员
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param idReq body req.IdReq true "管理员ID"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/admin/change-status [post]
func (r *adminRoutes) changeStatus(c *gin.Context) {
	var idReq req.IdReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &idReq)) {
		return
	}
	err := r.srv.ChangeStatus(idReq.ID)
	response.CheckAndResp(c, err)
}
