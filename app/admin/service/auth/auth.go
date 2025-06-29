package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/req"
	"github.com/zhoudm1743/go-flow/app/admin/schemas/resp"
	"github.com/zhoudm1743/go-flow/app/admin/service/system"
	"github.com/zhoudm1743/go-flow/app/models"
	"github.com/zhoudm1743/go-flow/core/database"
	"github.com/zhoudm1743/go-flow/pkg/captcha"
	"github.com/zhoudm1743/go-flow/pkg/jwt"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/util"
)

type AuthService interface {
	Login(c *gin.Context, req *req.LoginReq) (*resp.LoginResp, error)
	Logout(token string) error
	GetUserInfo(token string) (*resp.SystemAdminResp, error)
	RefreshToken(token string) (*resp.LoginResp, error)
	Captcha() (*resp.CaptchaResp, error)
}

type authService struct {
	usr     system.AdminService
	db      database.Database
	captcha *captcha.Service
	jwt     jwt.JwtService
}

// Captcha implements AuthService.
func (a *authService) Captcha() (*resp.CaptchaResp, error) {
	captchaID, b64s, _, err := a.captcha.Generate()
	if err != nil {
		return nil, err
	}
	return &resp.CaptchaResp{
		CaptchaID: captchaID,
		Captcha:   b64s,
	}, nil
}

// GetUserInfo implements AuthService.
func (a *authService) GetUserInfo(token string) (*resp.SystemAdminResp, error) {
	userID, err := a.jwt.ParseToken(token)
	if err != nil {
		return nil, err
	}
	admin, err := a.usr.Find(userID)
	if err != nil {
		return nil, err
	}
	var res resp.SystemAdminResp
	response.Copy(&res, admin)
	return &res, nil
}

// Login implements AuthService.
func (a *authService) Login(c *gin.Context, req *req.LoginReq) (*resp.LoginResp, error) {
	var admin models.SystemAdmin
	a.db.GetDB().Model(&models.SystemAdmin{}).Where("username = ?", req.Username).First(&admin)
	if admin.ID == 0 {
		return nil, errors.New("用户不存在")
	}
	pwd := util.ToolsUtil.MakeMd5(req.Password + "zhoudm1743")
	if pwd != admin.Password {
		return nil, errors.New("密码错误")
	}
	token, err := a.jwt.GenerateToken(admin.ID)
	if err != nil {
		return nil, err
	}
	return &resp.LoginResp{
		Token: token,
	}, nil
}

// Logout implements AuthService.
func (a *authService) Logout(token string) error {
	return nil
}

// RefreshToken implements AuthService.
func (a *authService) RefreshToken(token string) (*resp.LoginResp, error) {
	newToken, err := a.jwt.RefreshToken(token)
	if err != nil {
		return nil, err
	}
	return &resp.LoginResp{
		Token: newToken,
	}, nil
}

// NewAuthService 创建AuthService
func NewAuthService(usr system.AdminService, db database.Database, jwt jwt.JwtService, captcha *captcha.Service) AuthService {
	return &authService{
		usr:     usr,
		db:      db,
		jwt:     jwt,
		captcha: captcha,
	}
}
