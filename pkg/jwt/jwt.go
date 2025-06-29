package jwt

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhoudm1743/go-flow/core/config"
	"go.uber.org/fx"
)

type JwtService interface {
	GenerateToken(userID uint) (string, error)
	ParseToken(token string) (uint, error)
	RefreshToken(token string) (string, error)
	VerifyToken(token string) (bool, error)
	GetTokenFromContext(c *gin.Context) (string, error)
}

type jwtService struct {
	cfg *config.Config
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func (j *jwtService) GenerateToken(userID uint) (string, error) {
	// 解析过期时间
	expiresIn, err := time.ParseDuration(j.cfg.Jwt.ExpiresIn)
	if err != nil {
		return "", errors.New("无效的token过期时间配置")
	}

	// 创建Claims
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.cfg.Jwt.Issuer,
			Subject:   strconv.Itoa(int(userID)),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(j.cfg.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetTokenFromContext 从gin上下文中获取token
func (j *jwtService) GetTokenFromContext(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("token不能为空")
	}
	// 如果token以Bearer开头，则去掉Bearer
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	return token, nil
}

// ParseToken 解析token并返回用户ID
func (j *jwtService) ParseToken(token string) (uint, error) {
	// 解析token
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(j.cfg.Jwt.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	// 验证token是否有效
	if !parsedToken.Valid {
		return 0, errors.New("无效的token")
	}

	return claims.UserID, nil
}

// RefreshToken 刷新token
func (j *jwtService) RefreshToken(token string) (string, error) {
	// 首先解析旧token
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(j.cfg.Jwt.Secret), nil
	})

	// 如果token解析失败且不是过期错误，则返回错误
	if err != nil && !strings.Contains(err.Error(), "token is expired") {
		return "", err
	}

	// 检查刷新token的有效期
	refreshExpires, err := time.ParseDuration(j.cfg.Jwt.RefreshExpires)
	if err != nil {
		return "", errors.New("无效的刷新token过期时间配置")
	}

	// 如果token发布时间超过了刷新期限，则不允许刷新
	if time.Since(claims.IssuedAt.Time) > refreshExpires {
		return "", errors.New("token已超过刷新期限")
	}

	// 生成新的token
	return j.GenerateToken(claims.UserID)
}

// VerifyToken 验证token是否有效
func (j *jwtService) VerifyToken(token string) (bool, error) {
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(j.cfg.Jwt.Secret), nil
	})

	if err != nil {
		return false, err
	}

	return parsedToken.Valid, nil
}

// NewJwtService 创建JwtService
func NewJwtService(cfg *config.Config) JwtService {
	return &jwtService{
		cfg: cfg,
	}
}

// Module fx模块
var Module = fx.Options(
	fx.Provide(NewJwtService),
)
