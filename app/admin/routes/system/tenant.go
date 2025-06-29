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

type tenantRoutes struct {
	srv system.TenantService
	jwt jwt.JwtService
}

func NewTenantGroup(srv system.TenantService, jwt jwt.JwtService) httpCore.Group {
	return httpCore.NewGroup("/system",
		func() interface{} {
			return &tenantRoutes{srv: srv, jwt: jwt}
		},
		regTenant,
		middleware.AuthMiddleware(jwt),
	)
}

func regTenant(rg *httpCore.BaseGroup, instance interface{}) error {
	r := instance.(*tenantRoutes)
	rg.GET("/tenant/all", r.all)
	rg.GET("/tenant/list", r.list)
	rg.GET("/tenant/detail", r.detail)
	rg.POST("/tenant/add", r.add)
	rg.POST("/tenant/edit", r.edit)
	rg.POST("/tenant/delete", r.delete)
	rg.POST("/tenant/change-status", r.changeStatus)
	rg.POST("/tenant/test-connection", r.testConnection)
	return nil
}

// @Summary 获取所有租户
// @Description 获取所有租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param queryReq query req.SystemTenantListReq false "查询参数"
// @Success 200 {object} response.Response{data=[]resp.SystemTenantResp} "成功"
// @Router /system/tenant/all [get]
func (r *tenantRoutes) all(c *gin.Context) {
	var queryReq req.SystemTenantListReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &queryReq)) {
		return
	}
	res, err := r.srv.All(&queryReq)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 获取租户列表
// @Description 获取租户列表
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param pageReq query req.PageReq true "分页参数"
// @Param queryReq query req.SystemTenantListReq false "查询参数"
// @Success 200 {object} response.Response{data=response.PageResp{lists=[]resp.SystemTenantResp}} "成功"
// @Router /system/tenant/list [get]
func (r *tenantRoutes) list(c *gin.Context) {
	var pageReq req.PageReq
	var queryReq req.SystemTenantListReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &pageReq, &queryReq)) {
		return
	}
	res, err := r.srv.List(&pageReq, &queryReq)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 获取租户详情
// @Description 获取租户详情
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param idReq query req.IdReq true "租户ID"
// @Success 200 {object} response.Response{data=resp.SystemTenantResp} "成功"
// @Router /system/tenant/detail [get]
func (r *tenantRoutes) detail(c *gin.Context) {
	var idReq req.IdReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &idReq)) {
		return
	}
	res, err := r.srv.Find(idReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 添加租户
// @Description 添加租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param addReq body req.SystemTenantAddReq true "添加参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/tenant/add [post]
func (r *tenantRoutes) add(c *gin.Context) {
	var addReq req.SystemTenantAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &addReq)) {
		return
	}
	err := r.srv.Add(&addReq)
	response.CheckAndResp(c, err)
}

// @Summary 编辑租户
// @Description 编辑租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param editReq body req.SystemTenantEditReq true "编辑参数"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/tenant/edit [post]
func (r *tenantRoutes) edit(c *gin.Context) {
	var editReq req.SystemTenantEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &editReq)) {
		return
	}
	err := r.srv.Edit(&editReq)
	response.CheckAndResp(c, err)
}

// @Summary 删除租户
// @Description 删除租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param idReq body req.IdReq true "租户ID"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/tenant/delete [post]
func (r *tenantRoutes) delete(c *gin.Context) {
	var idReq req.IdReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &idReq)) {
		return
	}
	err := r.srv.Delete(idReq.ID)
	response.CheckAndResp(c, err)
}

// @Summary 修改租户状态
// @Description 修改租户状态
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param idReq body req.IdReq true "租户ID"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/tenant/change-status [post]
func (r *tenantRoutes) changeStatus(c *gin.Context) {
	var idReq req.IdReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &idReq)) {
		return
	}
	err := r.srv.ChangeStatus(idReq.ID)
	response.CheckAndResp(c, err)
}

// @Summary 测试租户数据库连接
// @Description 测试租户数据库连接
// @Tags 租户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "{token}"
// @Param idReq body req.IdReq true "租户ID"
// @Success 200 {object} response.Response{} "成功"
// @Router /system/tenant/test-connection [post]
func (r *tenantRoutes) testConnection(c *gin.Context) {
	var idReq req.IdReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &idReq)) {
		return
	}
	err := r.srv.TestConnection(idReq.ID)
	response.CheckAndResp(c, err)
}
