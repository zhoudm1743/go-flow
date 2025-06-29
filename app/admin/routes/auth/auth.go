package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/req"
	"github.com/zhoudm1743/go-flow/app/admin/service/auth"
	httpCore "github.com/zhoudm1743/go-flow/core/http"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/util"
)

type authRoutes struct {
	srv auth.AuthService
}

func NewAuthGroup(srv auth.AuthService) httpCore.Group {
	return httpCore.NewGroup("/auth",
		func() interface{} {
			return &authRoutes{srv: srv}
		},
		regAuth,
	)
}

func regAuth(rg *httpCore.BaseGroup, instance interface{}) error {
	r := instance.(*authRoutes)
	rg.POST("/login", r.login)
	rg.POST("/logout", r.logout)
	rg.GET("/refresh-token", r.refreshToken)
	rg.GET("/captcha", r.captcha)
	rg.GET("/user-info", r.userInfo)
	return nil
}

// @Summary 登录
// @Description 登录
// @Tags auth
// @Accept json
// @Produce json
// @Param loginReq body req.LoginReq true "登录请求"
// @Success 200 {object} response.Response "登录成功"
// @Failure 400 {object} response.Response "登录失败"
// @Router /api/v1/auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req req.LoginReq
	if response.IsFailWithResp(c, util.VerifyUtil.Verify(c, &req)) {
		return
	}
	res, err := r.srv.Login(c, &req)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 退出登录
// @Description 退出登录
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "退出登录成功"
// @Failure 400 {object} response.Response "退出登录失败"
// @Router /api/v1/auth/logout [post]
func (r *authRoutes) logout(c *gin.Context) {
	token := getToken(c)
	if token == "" {
		response.FailWithMsg(c, response.TokenEmpty, "token不能为空")
		c.Abort()
		return
	}
	err := r.srv.Logout(token)
	response.CheckAndRespWithData(c, nil, err)
}

// @Summary 刷新token
// @Description 刷新token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "刷新token成功"
// @Failure 400 {object} response.Response "刷新token失败"
// @Router /api/v1/auth/refresh-token [get]
func (r *authRoutes) refreshToken(c *gin.Context) {
	token := getToken(c)
	if token == "" {
		response.FailWithMsg(c, response.TokenEmpty, "token不能为空")
		c.Abort()
		return
	}
	res, err := r.srv.RefreshToken(token)
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 获取验证码
// @Description 获取验证码
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "获取验证码成功"
// @Failure 400 {object} response.Response "获取验证码失败"
// @Router /api/v1/auth/captcha [get]
func (r *authRoutes) captcha(c *gin.Context) {
	res, err := r.srv.Captcha()
	response.CheckAndRespWithData(c, res, err)
}

// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "获取用户信息成功"
// @Failure 400 {object} response.Response "获取用户信息失败"
// @Router /api/v1/auth/user-info [get]
func (r *authRoutes) userInfo(c *gin.Context) {
	token := getToken(c)
	if token == "" {
		response.FailWithMsg(c, response.TokenEmpty, "token不能为空")
		c.Abort()
		return
	}
	res, err := r.srv.GetUserInfo(token)
	response.CheckAndRespWithData(c, res, err)
}

func getToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	// 如果token以Bearer开头，则去掉Bearer
	if len(token) > 0 {
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}
	}
	return token
}
